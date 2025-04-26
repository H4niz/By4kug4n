import pytest
import grpc
from unittest.mock import MagicMock, patch, mock_open
from concurrent import futures

from byakugan.core.comms.grpc.server import ScannerServer
from byakugan.config.comms import CommsConfig
from byakugan.core.comms.exceptions import ConnectionError, SecurityError
from byakugan.core.comms.grpc import scanner_pb2, scanner_pb2_grpc

class TestScannerServer:
    @pytest.fixture
    def mock_servicer(self):
        """Mock scanner servicer"""
        servicer = MagicMock()
        servicer.ExecuteScan = MagicMock()
        servicer.StreamResults = MagicMock()
        return servicer

    @pytest.fixture
    def config(self):
        """Test configuration"""
        return CommsConfig(
            grpc_host="localhost",
            grpc_port=50051,
            grpc_max_workers=10,
            grpc_max_message_size=10 * 1024 * 1024,
            grpc_timeout=30,
            connection_timeout=5,
            connection_retries=3,
            tls_enabled=False,
            tls_cert_path=None,
            tls_key_path=None,
            worker_pool_size=5,
            queue_size=100,
            http_timeout=30,
            verify_ssl=False,
            max_retries=3
        )

    @pytest.fixture 
    def test_task(self):
        """Sample scan task"""
        return scanner_pb2.ScanTask(
            id="TEST-001",
            target=scanner_pb2.Target(
                url="http://test-api.com/users",
                method="GET",
                protocol="http"
            ),
            auth_context=scanner_pb2.AuthContext(
                type="bearer",
                headers={"Authorization": "Bearer test-token"}
            ),
            payload=scanner_pb2.Payload(
                headers={"Content-Type": "application/json"},
                query_params={"search": "test"},
                insertion_points=[
                    scanner_pb2.InsertionPoint(
                        location="query.search",
                        type="sql_injection",
                        payloads=["' OR '1'='1"],
                        encoding="url"
                    )
                ]
            ),
            validation=scanner_pb2.Validation(
                success_conditions=scanner_pb2.SuccessConditions(
                    status_codes=[200],
                    timing_threshold=5000,
                    error_patterns=["SQL syntax"]
                ),
                evidence_collection=scanner_pb2.EvidenceCollection(
                    save_request=True,
                    save_response=True,
                    timing_analysis=True
                )
            ),
            rule_context=scanner_pb2.RuleContext(
                rule_id="SQL-001",
                severity="HIGH",
                category="Injection"
            )
        )

    def test_initialization(self, config):
        """Test server initialization"""
        server = ScannerServer(config)
        assert server.config == config
        assert server._server is None
        assert server._servicer is None

    @patch('grpc.server')
    def test_start_success(self, mock_grpc_server, config, mock_servicer):
        """Test successful server start"""
        server = ScannerServer(config)
        mock_server = MagicMock()
        mock_grpc_server.return_value = mock_server

        # Mock servicer creation
        with patch('byakugan.core.comms.grpc.scanner_pb2_grpc.add_ScannerServiceServicer_to_server') as mock_add_servicer:
            assert server.start() is True
            
            assert server._server is mock_server
            mock_server.add_insecure_port.assert_called_once_with(
                f"{config.grpc_host}:{config.grpc_port}"
            )
            mock_add_servicer.assert_called_once()
            mock_server.start.assert_called_once()

    @patch('grpc.server')
    def test_start_error(self, mock_grpc_server, config):
        """Test server start error"""
        server = ScannerServer(config)
        mock_grpc_server.side_effect = Exception("Start failed")

        with pytest.raises(ConnectionError) as exc:
            server.start()
        assert "Failed to start server" in str(exc.value)
        assert server._server is None

    def test_configure_security_missing_certs(self, config):
        """Test security configuration with missing certs"""
        server = ScannerServer(config)
        config.tls_enabled = True
        
        with pytest.raises(SecurityError) as exc:
            server._configure_security()
        assert "TLS enabled but cert/key paths not provided" in str(exc.value)

    @patch('grpc.ssl_server_credentials')
    def test_configure_security_success(self, mock_creds, config):
        """Test successful security configuration"""
        config.tls_enabled = True
        config.tls_cert_path = "cert.pem"
        config.tls_key_path = "key.pem"
        
        server = ScannerServer(config)
        server._server = MagicMock()

        mock_creds.return_value = "test_creds"
        cert_data = b"TEST CERTIFICATE"
        key_data = b"TEST PRIVATE KEY"
        
        m = mock_open()
        m.side_effect = [
            mock_open(read_data=cert_data).return_value,
            mock_open(read_data=key_data).return_value
        ]
        
        with patch('builtins.open', m):
            server._configure_security()
            
        mock_creds.assert_called_once_with([(key_data, cert_data)])
        server._server.add_secure_port.assert_called_once_with(
            f"{config.grpc_host}:{config.grpc_port}",
            "test_creds"
        )

    def test_execute_scan(self, config):
        """Test scan execution"""
        server = ScannerServer(config)
        server._servicer = MagicMock()

        # Create test task
        task = scanner_pb2.ScanTask(
            id="TEST-001",
            target=scanner_pb2.Target(
                url="http://test.com",
                method="GET",
                protocol="http"
            )
        )

        # Mock servicer response
        mock_result = scanner_pb2.ScanResult(
            task_id="TEST-001",
            success=True,
            findings=[]
        )
        server._servicer.ExecuteScan.return_value = mock_result

        # Execute scan and verify
        result = server.execute_scan(task)
        assert result.task_id == "TEST-001"
        assert result.success is True
        server._servicer.ExecuteScan.assert_called_once_with(task, None)

    def test_stop(self, config):
        """Test server stop"""
        server = ScannerServer(config)
        mock_server = MagicMock()
        server._server = mock_server
        
        server.stop()
        mock_server.stop.assert_called_once_with(grace=5)
        assert server._server is None
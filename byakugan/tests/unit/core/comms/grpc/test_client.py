import pytest
import grpc
from unittest.mock import MagicMock, patch, mock_open
from concurrent import futures

from byakugan.core.comms.grpc.client import ScannerClient
from byakugan.config.comms import CommsConfig
from byakugan.core.comms.exceptions import ConnectionError, TimeoutError, SecurityError
from byakugan.core.comms.grpc.proto import scanner_pb2


class TestScannerClient:
    @pytest.fixture
    def mock_stub_class(self):
        """Mock ScannerStub class"""
        with patch('byakugan.core.comms.grpc.scanner_pb2_grpc.ScannerServiceStub') as mock:
            yield mock

    @pytest.fixture
    def mock_channel(self):
        """Mock gRPC channel"""
        channel = MagicMock() 
        return channel

    @pytest.fixture
    def config(self):
        """Test configuration"""
        return CommsConfig(
            grpc_host="localhost",
            grpc_port=50051,
            grpc_timeout=30,
            connection_timeout=5,
            connection_retries=3,
            tls_enabled=False
        )

    @pytest.fixture
    def test_task(self):
        """Sample scan task"""
        return {
            "id": "TEST-001",
            "target": {
                "url": "http://test-api.com/users",
                "method": "GET",
                "protocol": "http"
            },
            "auth_context": {
                "type": "bearer",
                "headers": {
                    "Authorization": "Bearer test-token"
                }
            },
            "payload": {
                "headers": {"Content-Type": "application/json"},
                "query_params": {"search": "test"},
                "insertion_points": [
                    {
                        "location": "query.search",
                        "type": "sql_injection",
                        "payloads": ["' OR '1'='1"],
                        "encoding": "url"
                    }
                ]
            },
            "validation": {
                "success_conditions": {
                    "status_codes": [200],
                    "timing_threshold": 5000,
                    "error_patterns": ["SQL syntax"]
                },
                "evidence_collection": {
                    "save_request": True,
                    "save_response": True,
                    "timing_analysis": True
                }
            },
            "rule_context": {
                "rule_id": "SQL-001",
                "severity": "HIGH",
                "category": "Injection"
            }
        }

    def test_initialization(self, config):
        """Test client initialization"""
        client = ScannerClient(config)
        assert client.config == config
        assert client._channel is None
        assert client._stub is None

    @patch('grpc.insecure_channel')
    @patch('grpc.channel_ready_future')
    def test_connect_success(self, mock_ready_future, mock_channel, mock_stub_class, config):
        """Test successful connection"""
        mock_ready_future.return_value.result.return_value = None
        client = ScannerClient(config)
        
        assert client.connect() is True
        assert client._channel is not None
        assert client._stub is not None
        mock_channel.assert_called_once_with(f"{config.grpc_host}:{config.grpc_port}")

    @patch('grpc.secure_channel')
    @patch('grpc.ssl_channel_credentials')
    @patch('grpc.channel_ready_future')
    def test_secure_connection(self, mock_ready_future, mock_creds, mock_secure_channel, config):
        """Test secure connection with TLS"""
        config.tls_enabled = True
        config.tls_cert_path = "cert.pem"
        mock_ready_future.return_value.result.return_value = None
        cert_data = b"TEST CERTIFICATE"
        
        with patch('builtins.open', mock_open(read_data=cert_data)):
            client = ScannerClient(config)
            assert client.connect() is True
            
        mock_creds.assert_called_once_with(cert_data)
        mock_secure_channel.assert_called_once()

    def test_retry_mechanism(self, config, test_task, mock_stub_class):
        """Test retry mechanism"""
        client = ScannerClient(config)
        client._stub = MagicMock()
        
        # Mock responses
        error_response = grpc.RpcError()
        success_response = MagicMock(
            task_id="TEST-001",
            success=True,
            findings=[
                MagicMock(
                    id="FIND-001",
                    rule_id="SQL-001",
                    severity="HIGH",
                    evidence=MagicMock(
                        data={"request": "GET /users", "response": "200 OK"},
                        description="SQL injection found"
                    )
                )
            ]
        )
        
        client._stub.ExecuteScan.side_effect = [error_response, success_response]
        
        result = client.execute_scan(test_task)
        
        assert client._stub.ExecuteScan.call_count == 2
        assert result["task_id"] == "TEST-001"
        assert result["success"] is True
        assert len(result["findings"]) == 1
        assert result["findings"][0]["rule_id"] == "SQL-001"

    def test_close(self, config):
        """Test closing client connection"""
        client = ScannerClient(config)
        # Create mock channel
        mock_channel = MagicMock()
        client._channel = mock_channel
        # Create mock executor
        mock_executor = MagicMock()
        client._executor = mock_executor
        
        client.close()
        
        # Verify channel closed
        mock_channel.close.assert_called_once()
        assert client._channel is None
        # Verify executor shutdown
        mock_executor.shutdown.assert_called_once_with(wait=True)

def test_scanner_client(comms_config, mock_channel):
    client = ScannerClient(comms_config)
    client._channel = mock_channel
    assert client._channel is not None
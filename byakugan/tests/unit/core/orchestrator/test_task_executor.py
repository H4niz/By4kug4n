import pytest
from unittest.mock import MagicMock, patch

from byakugan.core.orchestrator.task_executor import TaskExecutor
from byakugan.config.comms import CommsConfig
from byakugan.core.comms.grpc.proto.scanner_pb2 import ScanTask, Target, Payload, InsertionPoint, ScanResult, Finding, Evidence
from byakugan.core.comms.grpc.proto.scanner_pb2_grpc import ScannerServiceStub

class TestTaskExecutor:
    @pytest.fixture
    def mock_channel(self):
        """Create mock gRPC channel"""
        channel = MagicMock()
        channel.close = MagicMock()
        return channel
    
    @pytest.fixture
    def mock_future(self):
        """Create mock channel ready future"""
        future = MagicMock()
        future.result = MagicMock(return_value=None)
        return future

    @pytest.fixture
    def config(self):
        """Test configuration"""
        return CommsConfig(
            grpc_host="localhost",
            grpc_port=50051,
            grpc_timeout=30
        )

    @pytest.fixture
    def test_task(self):
        """Create test scan task"""
        return ScanTask(
            id="TEST-001",
            target=Target(
                url="http://test.com",
                method="GET",
                protocol="http"
            ),
            payload=Payload(
                headers={"Content-Type": "application/json"},
                query_params={"q": "test"},
                insertion_points=[
                    InsertionPoint(
                        location="query.q",
                        type="sql_injection",
                        payloads=["' OR '1'='1"]
                    )
                ]
            )
        )

    def test_task_executor_init(self, config, mock_channel):
        """Test TaskExecutor initialization"""
        executor = TaskExecutor(config)
        assert executor.config == config

    def test_initialization(self, config, mock_channel, mock_future):
        """Test initialization with proper mocking"""
        with patch('byakugan.core.comms.grpc.client.ScannerClient') as mock_client_class:
            mock_client = mock_client_class.return_value
            mock_client.connect.return_value = True
            mock_client._create_channel.return_value = mock_channel
            
            # Force mock client creation before executor
            mock_client_class.reset_mock()
            
            executor = TaskExecutor(config)
            assert executor.config == config
            assert executor._client is not None
            mock_client_class.assert_called_once()

    def test_execute_task_success(self, config, mock_channel, mock_future, test_task):
        """Test successful task execution"""
        with patch('byakugan.core.comms.grpc.client.ScannerClient') as mock_client_class, \
             patch('grpc.insecure_channel', return_value=mock_channel), \
             patch('grpc.channel_ready_future', return_value=mock_future):

            # Setup mocks
            mock_client = mock_client_class.return_value
            mock_client._create_channel = MagicMock(return_value=mock_channel)
            mock_client.connect.return_value = True

            # Mock successful scan result
            mock_result = ScanResult(
                task_id="TEST-001",
                success=True,
                findings=[
                    Finding(
                        id="FIND-001",
                        rule_id="SQLI-001",
                        severity="HIGH",
                        evidence=Evidence(
                            data={"request": "GET /test"},
                            description="SQL injection found"
                        )
                    )
                ]
            )
            mock_client.execute_scan.return_value = mock_result

            # Execute test
            executor = TaskExecutor(config)
            result = executor.execute_task(test_task)

            # Verify execution
            assert result.task_id == "TEST-001"
            assert result.success is True
            mock_client.execute_scan.assert_called_once_with(test_task)

    def test_execute_task_failure(self, config, test_task):
        """Test failed task execution"""
        with patch('byakugan.core.comms.grpc.client.ScannerClient') as mock_client_class:
            # Setup mock client
            mock_client = mock_client_class.return_value
            mock_client.connect.return_value = True
            mock_client.execute_scan.side_effect = Exception("Scan failed")

            # Create executor and execute task
            executor = TaskExecutor(config)
            result = executor.execute_task(test_task)

            # Verify failure handling 
            assert result.success is False
            assert result.task_id == test_task.id
            mock_client.execute_scan.assert_called_once_with(test_task)

    def test_task_validation(self, config):
        """Test task validation"""
        with patch('byakugan.core.comms.grpc.client.ScannerClient') as mock_client_class:
            # Setup mock client
            mock_client = mock_client_class.return_value
            mock_client.connect.return_value = True

            # Create executor
            executor = TaskExecutor(config)

            # Test invalid task
            invalid_task = ScanTask()
            result = executor.execute_task(invalid_task)

            # Verify validation
            assert result.success is False
            assert result.task_id == ""
            mock_client.execute_scan.assert_not_called()
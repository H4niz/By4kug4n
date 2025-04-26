import pytest
import grpc
from unittest.mock import MagicMock
from byakugan.core.comms.grpc.servicer import ScannerServicer
from byakugan.core.comms.grpc import scanner_pb2

class TestScannerServicer:
    @pytest.fixture
    def servicer(self):
        """Create servicer instance"""
        servicer = ScannerServicer()
        servicer.engine = MagicMock()
        return servicer

    @pytest.fixture
    def test_task(self):
        """Create test scan task"""
        return scanner_pb2.ScanTask(
            id="TEST-001",
            target=scanner_pb2.Target(
                url="http://test.com",
                method="GET",
                protocol="http"
            ),
            auth_context=scanner_pb2.AuthContext(
                type="bearer",
                headers={"Authorization": "Bearer test"}
            ),
            payload=scanner_pb2.Payload(
                headers={"Content-Type": "application/json"},
                query_params={"q": "test"},
                insertion_points=[
                    scanner_pb2.InsertionPoint(
                        location="query.q",
                        type="sql_injection",
                        payloads=["' OR '1'='1"]
                    )
                ]
            ),
            validation=scanner_pb2.Validation(
                success_conditions=scanner_pb2.SuccessConditions(
                    status_codes=[200],
                    timing_threshold=5000,
                    error_patterns=["error"]
                ),
                evidence_collection=scanner_pb2.EvidenceCollection(
                    save_request=True,
                    save_response=True,
                    timing_analysis=True
                )
            )
        )

    def test_execute_scan_success(self, servicer, test_task):
        """Test successful scan execution"""
        # Mock engine response
        mock_result = {
            "task_id": "TEST-001",
            "success": True,
            "findings": [
                {
                    "id": "FIND-001",
                    "rule_id": "SQL-001",
                    "severity": "HIGH",
                    "evidence": {
                        "data": {"request": "GET /test"},
                        "description": "SQL injection found"
                    }
                }
            ]
        }
        servicer.engine.execute_task.return_value = mock_result

        # Execute scan
        result = servicer.ExecuteScan(test_task, None)

        # Verify result
        assert result.task_id == "TEST-001"
        assert result.success is True
        assert len(result.findings) == 1
        assert result.findings[0].id == "FIND-001"

        # Verify engine called correctly
        servicer.engine.execute_task.assert_called_once()

    def test_execute_scan_error(self, servicer, test_task):
        """Test error handling"""
        # Mock engine error
        servicer.engine.execute_task.side_effect = Exception("Test error")

        # Create mock context
        context = MagicMock()

        # Execute scan
        result = servicer.ExecuteScan(test_task, context)

        # Verify error handling
        assert isinstance(result, scanner_pb2.ScanResult)
        context.set_code.assert_called_with(grpc.StatusCode.INTERNAL)
        context.set_details.assert_called_with("Test error")

    def test_result_conversion(self, servicer):
        """Test result conversion"""
        internal_result = {
            "task_id": "TEST-001",
            "success": True,
            "findings": [
                {
                    "id": "FIND-001",
                    "rule_id": "SQL-001",
                    "severity": "HIGH",
                    "evidence": {
                        "data": {"request": "GET /test"},
                        "description": "SQL injection found"
                    }
                }
            ]
        }

        result = servicer._convert_result(internal_result)
        assert isinstance(result, scanner_pb2.ScanResult)
        assert result.task_id == "TEST-001"
        assert result.success is True
        assert len(result.findings) == 1
        assert result.findings[0].id == "FIND-001"
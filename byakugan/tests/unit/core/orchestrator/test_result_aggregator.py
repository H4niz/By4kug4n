import pytest
from datetime import datetime
from byakugan.core.orchestrator.result_aggregator import ResultAggregator, ScanSummary
from byakugan.core.comms.grpc import scanner_pb2

class TestResultAggregator:
    @pytest.fixture
    def aggregator(self):
        """Create result aggregator instance"""
        return ResultAggregator()

    def test_add_successful_result(self, aggregator):
        """Test adding successful scan results"""
        # Create test results
        result1 = scanner_pb2.ScanResult(
            task_id="TASK-001",
            success=True,
            findings=[
                scanner_pb2.Finding(
                    id="FIND-001",
                    rule_id="SQLI-001",
                    severity="HIGH",
                    evidence=scanner_pb2.Evidence(
                        data={"request": "GET /test"},
                        description="SQL injection found"
                    )
                )
            ]
        )
        
        result2 = scanner_pb2.ScanResult(
            task_id="TASK-002",
            success=True,
            findings=[
                scanner_pb2.Finding(
                    id="FIND-002",
                    rule_id="XSS-001",
                    severity="MEDIUM",
                    evidence=scanner_pb2.Evidence(
                        data={"request": "POST /test"},
                        description="XSS found"
                    )
                )
            ]
        )

        # Add results
        aggregator.add_result(result1)
        aggregator.add_result(result2)

        # Get summary
        summary = aggregator.get_summary()

        # Verify summary
        assert summary.total_tasks == 2
        assert summary.completed_tasks == 2
        assert summary.failed_tasks == 0
        assert summary.total_findings == 2
        assert summary.severity_counts == {
            'critical': 0,
            'high': 1,
            'medium': 1, 
            'low': 0,
            'info': 0
        }

    def test_add_failed_result(self, aggregator):
        """Test adding failed scan results"""
        result = scanner_pb2.ScanResult(
            task_id="TASK-003",
            success=False,
            findings=[]
        )

        aggregator.add_result(result)
        summary = aggregator.get_summary()

        assert summary.total_tasks == 1
        assert summary.completed_tasks == 0
        assert summary.failed_tasks == 1
        assert summary.total_findings == 0
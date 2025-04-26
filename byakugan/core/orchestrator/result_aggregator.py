from typing import Dict, List
from dataclasses import dataclass
from datetime import datetime
from byakugan.core.comms.grpc import scanner_pb2
from .models import ScanResult

@dataclass
class ScanSummary:
    """Summary of scan results"""
    total_tasks: int = 0
    completed_tasks: int = 0
    failed_tasks: int = 0
    total_findings: int = 0
    severity_counts: Dict[str, int] = None
    execution_time: float = 0

    def __post_init__(self):
        if self.severity_counts is None:
            self.severity_counts = {
                'critical': 0,
                'high': 0,
                'medium': 0,
                'low': 0,
                'info': 0
            }

@dataclass 
class AggregatedResults:
    """Aggregated scan results"""
    total_tasks: int
    successful_tasks: int 
    failed_tasks: int
    findings: List[Dict]
    errors: List[str]

class ResultAggregator:
    """Aggregates and analyzes scan results"""

    def __init__(self):
        self._summary = ScanSummary()
        self._results = []
        self.start_time = None
        self.end_time = None

    def start_scan(self):
        """Mark scan start time"""
        self.start_time = datetime.utcnow()

    def end_scan(self):
        """Mark scan end time"""
        self.end_time = datetime.utcnow()

    def add_result(self, result: scanner_pb2.ScanResult):
        """Add scan result to aggregator"""
        self._results.append(result)
        self._update_summary(result)

    def _update_summary(self, result: scanner_pb2.ScanResult):
        """Update summary statistics"""
        self._summary.total_tasks += 1
        
        if result.success:
            self._summary.completed_tasks += 1
            self._summary.total_findings += len(result.findings)
            
            # Update severity counts
            for finding in result.findings:
                severity = finding.severity.lower()
                if severity in self._summary.severity_counts:
                    self._summary.severity_counts[severity] += 1
        else:
            self._summary.failed_tasks += 1

    def get_summary(self) -> ScanSummary:
        """Get current summary"""
        self._summary.execution_time = (self.end_time - self.start_time).total_seconds() if self.end_time and self.start_time else 0
        return self._summary

    def generate_report(self) -> Dict:
        """Generate detailed scan report"""
        return {
            "summary": self.get_summary(),
            "findings": [finding for result in self._results for finding in result.findings],
            "errors": [
                {
                    "task_id": result.task_id,
                    "error": result.error_message,
                    "status": "failed"
                }
                for result in self._results if not result.success
            ],
            "metadata": {
                "start_time": self.start_time.isoformat() if self.start_time else None,
                "end_time": self.end_time.isoformat() if self.end_time else None
            }
        }

    def process_results(self, results: List[ScanResult]) -> AggregatedResults:
        """Process and aggregate scan results"""
        successful = 0
        failed = 0
        all_findings = []
        errors = []

        for result in results:
            if result.success:
                successful += 1
                all_findings.extend(result.findings)
            else:
                failed += 1
                if result.error_details:
                    errors.append(result.error_details)

        return AggregatedResults(
            total_tasks=len(results),
            successful_tasks=successful,
            failed_tasks=failed, 
            findings=all_findings,
            errors=errors
        )
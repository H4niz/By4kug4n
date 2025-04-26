import pytest
from byakugan.core.metrics.collector import MetricsCollector

@pytest.fixture
def collector():
    return MetricsCollector()

def test_request_metrics(collector):
    collector.record_request(0.5, True)
    collector.record_request(1.0, True)
    collector.record_request(35.0, False)
    
    metrics = collector.get_metrics()
    assert metrics.total_requests == 3
    assert metrics.successful_requests == 2
    assert metrics.failed_requests == 1
    assert 0.5 < metrics.avg_response_time < 1.5

def test_finding_metrics(collector):
    collector.record_finding({"severity": "high", "title": "Finding 1"})
    collector.record_finding({"severity": "high", "title": "Finding 2"})
    collector.record_finding({"severity": "medium", "title": "Finding 3"})
    
    metrics = collector.get_metrics()
    assert metrics.findings_by_severity["high"] == 2
    assert metrics.findings_by_severity["medium"] == 1

def test_scan_duration(collector):
    collector.start_scan()
    import time
    time.sleep(0.1)  # Simulate scan duration
    collector.end_scan()
    
    metrics = collector.get_metrics()
    assert metrics.scan_duration > 0
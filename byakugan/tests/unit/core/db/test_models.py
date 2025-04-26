import pytest
from datetime import datetime
from byakugan.core.db.models import Scan, Finding

def test_scan_model():
    scan = Scan(
        target_url="https://api.example.com",
        status="running",
        config={"concurrent_tasks": 5}
    )
    
    assert scan.target_url == "https://api.example.com"
    assert scan.status == "running"
    assert isinstance(scan.start_time, datetime)
    assert scan.end_time is None
    
    scan_dict = scan.to_dict()
    assert "id" in scan_dict
    assert "target_url" in scan_dict
    assert "status" in scan_dict

def test_finding_model():
    finding = Finding(
        title="SQL Injection",
        severity="high",
        cvss_score=8.5,
        description="SQL injection vulnerability found",
        evidence={"request": "...", "response": "..."}
    )
    
    assert finding.title == "SQL Injection"
    assert finding.severity == "high"
    assert finding.cvss_score == 8.5
    assert isinstance(finding.evidence, dict)
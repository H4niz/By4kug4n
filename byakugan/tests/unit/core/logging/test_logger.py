import pytest
import logging
from pathlib import Path
from byakugan.core.logging.logger import ScanLogger

@pytest.fixture
def logger():
    config = {"log_dir": "test_logs"}
    logger = ScanLogger(config)
    yield logger
    # Cleanup
    Path("test_logs").rmdir()

def test_log_scan_start(logger, caplog):
    scan_config = {"target": "http://example.com", "rules": ["RULE-001"]}
    with caplog.at_level(logging.INFO):
        logger.log_scan_start(scan_config)
    
    assert "Starting new scan" in caplog.text
    assert "scan_start" in caplog.text

def test_log_finding(logger, caplog):
    finding = {
        "title": "SQL Injection",
        "severity": "high",
        "description": "SQL injection vulnerability"
    }
    
    with caplog.at_level(logging.WARNING):
        logger.log_finding(finding)
    
    assert "SQL Injection" in caplog.text
    assert "high" in caplog.text

def test_log_error(logger, caplog):
    error = ValueError("Invalid configuration")
    context = {"phase": "initialization"}
    
    with caplog.at_level(logging.ERROR):
        logger.log_error(error, context)
    
    assert "Invalid configuration" in caplog.text
    assert "initialization" in caplog.text
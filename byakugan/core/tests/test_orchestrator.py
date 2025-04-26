import pytest
import asyncio
from datetime import datetime
from uuid import UUID

from core.orchestrator import (
    ScanManager,
    ScanConfig,
    ScanStatus,
    TaskCoordinator,
    ResultAggregator
)

@pytest.fixture
def scan_config():
    return ScanConfig(
        api_definition=r"byakugan\scanner\tests\testdata\test_api.yaml",
        rules_dir="test_rules",
        concurrent_limit=2,
        task_timeout=5
    )

@pytest.fixture
def scan_manager(scan_config):
    return ScanManager(scan_config)

@pytest.mark.asyncio
async def test_scan_lifecycle(scan_manager):
    # Start scan
    scan_id = await scan_manager.start_scan()
    assert isinstance(scan_id, UUID)
    
    # Check initial status
    status = scan_manager.get_scan_status(scan_id)
    assert status.status == "running"
    assert status.progress >= 0
    
    # Wait for completion
    while status.status == "running":
        await asyncio.sleep(0.1)
        status = scan_manager.get_scan_status(scan_id)
    
    assert status.status == "completed"
    assert status.progress == 100.0
    assert status.end_time is not None

@pytest.mark.asyncio
async def test_scan_stopping(scan_manager):
    # Start scan
    scan_id = await scan_manager.start_scan()
    
    # Stop scan
    stopped = await scan_manager.stop_scan(scan_id)
    assert stopped is True
    
    # Verify status
    status = scan_manager.get_scan_status(scan_id)
    assert status.status == "stopped"
    assert status.end_time is not None

def test_error_handling(scan_manager):
    # Test with invalid config
    with pytest.raises(Exception):
        asyncio.run(scan_manager.start_scan())

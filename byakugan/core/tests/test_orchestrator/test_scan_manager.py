import pytest
import asyncio
from pathlib import Path
from uuid import UUID
from byakugan.core.orchestrator import ScanManager, ScanConfig, ScanStatus

@pytest.fixture
def test_api_content():
    """Test API spec content"""
    return """
    openapi: 3.0.0
    info:
      title: Test API
      version: 1.0.0
    paths:
      /api/test:
        get:
          parameters:
            - name: id
              in: query 
              required: true
              schema:
                type: string 
    """

@pytest.fixture
def scan_config(tmp_path, test_api_content):
    """Create scan config with test API file"""
    api_file = tmp_path / "test_api.yaml"
    api_file.write_text(test_api_content)
    
    return ScanConfig(
        api_definition=str(api_file),
        rules_dir="testdata/test_rules",
        concurrent_limit=2,
        task_timeout=5
    )

@pytest.mark.asyncio
async def test_scan_lifecycle(scan_config):
    """Test complete scan lifecycle"""
    # Create manager
    manager = ScanManager(scan_config)

    # Start scan
    scan_id = await manager.start_scan()
    assert isinstance(scan_id, UUID)

    # Get initial status 
    status = manager.get_scan_status(scan_id)
    assert status is not None
    assert status.status == "running"

    # Let scan complete
    await asyncio.sleep(0.1)

    # Check final status
    final_status = manager.get_scan_status(scan_id)
    assert final_status.status == "completed"
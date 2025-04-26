import pytest
from uuid import uuid4
from datetime import datetime
from byakugan.core.orchestrator import (
    TaskCoordinator, 
    TaskCoordinatorConfig,
    ScanTask
)

@pytest.fixture
def task_coordinator_config():
    """Create task coordinator config"""
    return TaskCoordinatorConfig(
        concurrent_limit=2,
        task_timeout=5,
        retry_count=2,
        retry_delay=1
    )

@pytest.fixture
def task_coordinator(task_coordinator_config):
    """Create task coordinator"""
    return TaskCoordinator(task_coordinator_config)

@pytest.fixture
def test_scan_task():
    """Create test scan task with complete parameters"""
    return ScanTask(
        id=str(uuid4()),
        scan_id=uuid4(),
        endpoint={
            "path": "/api/test",
            "method": "GET",
            "protocol": "https"
        },
        auth_context={
            "type": "OAuth2",
            "token": "test_token",
            "headers": {
                "Authorization": "Bearer test_token"
            }
        },
        rule={
            "id": "TEST-001",
            "severity": "HIGH",
            "category": "Path Traversal",
            "detection": {
                "pattern": "root:.*:0:0:",
                "location": "response_body"
            }
        },
        payload={
            "path_params": {
                "doc_id": "{{injection_point}}"
            },
            "headers": {
                "Accept": "*/*"
            },
            "insertion_points": [{
                "location": "path.doc_id",
                "type": "path_traversal",
                "payloads": [
                    "../../../etc/passwd",
                    "..%2f..%2f..%2fetc%2fpasswd"
                ]
            }]
        },
        validation={
            "success_conditions": {
                "status_codes": [200],
                "response_patterns": ["success"]
            },
            "evidence_collection": {
                "save_request": True,
                "save_response": True
            }
        }
    )

@pytest.mark.asyncio
async def test_task_execution(task_coordinator, test_scan_task):
    """Test task execution"""
    results = await task_coordinator.execute_tasks([test_scan_task])
    
    assert len(results) == 1
    assert results[0]["task_id"] == test_scan_task.id
    assert results[0]["success"] == True
    assert "findings" in results[0]

@pytest.mark.asyncio
async def test_task_concurrency(task_coordinator):
    """Test concurrent task execution"""
    # Create multiple tasks
    tasks = [
        ScanTask(
            id=str(uuid4()),
            scan_id=uuid4(),
            endpoint={
                "path": f"/api/test/{i}",
                "method": "GET"
            },
            rule={
                "id": f"TEST-{i:03d}",
                "severity": "high"
            }
        )
        for i in range(5)
    ]

    # Execute tasks concurrently  
    results = await task_coordinator.execute_tasks(tasks)
    
    # Verify results
    assert len(results) == len(tasks)
    assert all(r["success"] for r in results)

@pytest.mark.asyncio
async def test_task_retry(task_coordinator, test_scan_task):
    """Test task retry logic"""
    # Force task to fail
    test_scan_task.endpoint["path"] = "/invalid"
    
    results = await task_coordinator.execute_tasks([test_scan_task])
    
    assert len(results) == 1
    assert results[0]["retry_count"] > 0
    assert results[0]["success"] == False
    assert "error_details" in results[0]
    
@pytest.mark.asyncio
async def test_task_success_after_retry(task_coordinator, test_scan_task):
    """Test task succeeds after retry"""
    # Setup task to fail first attempt only
    original_execute = task_coordinator._execute_single_attempt
    attempt = 0
    
    async def mock_execute(task):
        nonlocal attempt
        attempt += 1
        if attempt == 1:
            raise ValueError("First attempt fails")
        return await original_execute(task)
        
    task_coordinator._execute_single_attempt = mock_execute
    
    results = await task_coordinator.execute_tasks([test_scan_task])
    
    assert len(results) == 1
    assert results[0]["retry_count"] == 1
    assert results[0]["success"] == True
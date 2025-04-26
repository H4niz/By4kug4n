import pytest
import asyncio
from uuid import uuid4
from byakugan.core.orchestrator import (
    TaskCoordinator,
    TaskCoordinatorConfig,
    ScanTask
)

@pytest.fixture
def task_coordinator():
    config = TaskCoordinatorConfig(
        concurrent_limit=2,
        task_timeout=5
    )
    return TaskCoordinator(config)

@pytest.fixture
def test_tasks():
    return [
        ScanTask(
            id=str(uuid4()),
            scan_id=uuid4(),
            endpoint={
                "path": "/api/test",
                "method": "GET"
            },
            rule={
                "id": "TEST-001",
                "severity": "HIGH"
            }
        )
    ]

@pytest.mark.asyncio
async def test_task_execution(task_coordinator, test_tasks):
    """Test task execution"""
    results = await task_coordinator.execute_tasks(test_tasks)
    assert len(results) == 1
    assert results[0]["task_id"] == test_tasks[0].id
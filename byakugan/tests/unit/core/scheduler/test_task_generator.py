import pytest
from byakugan.core.scheduler.task_generator import TaskGenerator

@pytest.fixture
def task_generator():
    config = {
        "scanner_version": "1.0.0",
        "validation": {
            "status_codes": [200, 201, 202]
        }
    }
    return TaskGenerator(config)

def test_create_jwt_scan_task(task_generator):
    target_url = "https://api.example.com/auth"
    token = "test.jwt.token"
    
    task = task_generator.create_jwt_scan_task(target_url, token)
    
    assert task["task_id"].startswith("TASK-JWT-")
    assert task["target"]["url"] == target_url
    assert task["auth_context"]["type"] == "JWT"
    assert task["auth_context"]["token"] == token
    assert len(task["payload"]["insertion_points"]) > 0
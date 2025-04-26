import pytest
from uuid import uuid4
from pathlib import Path
from byakugan.core.config import OrchestartorConfig
from byakugan.core.orchestrator import (
    TaskCoordinator,
    TaskCoordinatorConfig, 
    ScanConfig
)

@pytest.fixture
def orchestrator_config():
    return OrchestartorConfig(
        scan_timeout=300,
        max_concurrent_scans=2,
        task_batch_size=10
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
    """Create task coordinator with config"""
    return TaskCoordinator(task_coordinator_config)

@pytest.fixture
def test_api_file(tmp_path):
    """Create test API file"""
    api_file = tmp_path / "test_api.yaml"
    api_file.write_text("""
    openapi: 3.0.0
    info:
      title: Test API
      version: 1.0.0
    """)
    return str(api_file)

@pytest.fixture(scope="session")
def test_api_content():
    """Load test API spec"""
    test_file = Path("byakugan/scanner/tests/testdata/test_api.yaml")
    return test_file.read_text()

@pytest.fixture(scope="session")
def test_rules_dir():
    """Get test rules directory"""
    return Path("byakugan/scanner/tests/testdata/test_rules")
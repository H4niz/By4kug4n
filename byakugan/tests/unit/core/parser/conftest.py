import pytest
import os

@pytest.fixture
def test_data_dir():
    """Return path to test data directory"""
    return os.path.join(os.path.dirname(__file__), "testdata")

@pytest.fixture(autouse=True)
def setup_logger():
    """Configure logging for tests"""
    import logging
    logging.basicConfig(level=logging.INFO)
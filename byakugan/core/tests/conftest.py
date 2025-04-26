import sys
import pytest
import os
import time
from pathlib import Path
from uuid import uuid4

# Add project root to Python path
project_root = Path(__file__).parent.parent.parent
sys.path.insert(0, str(project_root))

from byakugan.core.parser.models import ApiDefinition, Endpoint
from byakugan.core.rule_engine import Rule
from byakugan.core.config import ParserConfig
from byakugan.core.parser import ApiParser
from byakugan.core.orchestrator import TaskCoordinator, ScanTask

@pytest.fixture
def parser_config():
    """Parser configuration fixture"""
    return ParserConfig()

@pytest.fixture  
def api_definition():
    """Test API definition fixture"""
    return ApiDefinition(
        title="Test API",
        version="1.0",
        endpoints=[
            Endpoint(
                path="/users",
                method="GET",
                parameters=[{
                    "name": "id",
                    "location": "query",
                    "type": "string"
                }]
            )
        ]
    )

@pytest.fixture
def test_rules():
    """Test rules fixture"""
    return [
        Rule(
            id="TEST-BOLA-001",
            name="Test BOLA Rule",
            description="Test rule for BOLA detection",
            severity="HIGH",
            category="authorization",
            method="GET",
            required_parameters=[{"name": "id", "location": "path"}],
            detection={
                "locations": {
                    "path_parameters": ["id", "user_id"],
                    "query_parameters": ["id"]
                },
                "strategies": [{"name": "ID Enumeration", "active": True}]
            },
            payloads={
                "id_manipulation": [
                    {"template": "{{base_id + 1}}"},
                    {"template": "{{base_id - 1}}"}
                ]
            },
            patterns=["error", "invalid"]
        )
    ]

@pytest.fixture
def api_parser():
    return ApiParser()

@pytest.fixture
def task_coordinator():
    return TaskCoordinator({
        "concurrent_limit": 2,
        "task_timeout": 5
    })

@pytest.fixture
def test_task():
    return ScanTask(
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

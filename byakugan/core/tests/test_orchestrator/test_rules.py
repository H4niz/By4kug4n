import pytest
from byakugan.core.parser import ApiParser

@pytest.fixture
def api_spec():
    return {
        "openapi": "3.0.0",
        "info": {
            "title": "Test API",
            "version": "1.0.0"
        },
        "paths": {
            "/api/test": {
                "get": {
                    "parameters": [{
                        "name": "id",
                        "in": "query",
                        "required": True, 
                        "schema": {"type": "string"}
                    }]
                }
            }
        }
    }

@pytest.fixture
def api_parser():
    return ApiParser()

@pytest.mark.asyncio 
async def test_rule_execution(api_parser, api_spec):
    """Test rule execution with API spec"""
    api_def = api_parser.parse_dict(api_spec)
    assert api_def is not None
    assert len(api_def.endpoints) == 1
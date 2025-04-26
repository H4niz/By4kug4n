import pytest
from core.parser import GraphQLParser, OpenAPIParser
from core.parser.models import ApiDefinition, Endpoint

def test_graphql_parser():
    parser = GraphQLParser()
    schema = """
    type Query {
        user(id: ID!): User
    }
    type User {
        id: ID!
        name: String! 
    }
    """
    
    # Test parse
    api_def = parser.parse(schema)
    assert isinstance(api_def, ApiDefinition)
    assert len(api_def.endpoints) == 1
    
    # Test endpoint properties
    endpoint = api_def.endpoints[0]
    assert endpoint.path == "/graphql"
    assert endpoint.method == "POST"
    assert len(endpoint.parameters) == 1

def test_openapi_parser():
    parser = OpenAPIParser() 
    spec = {
        "openapi": "3.0.0",
        "paths": {
            "/users": {
                "get": {
                    "parameters": []
                }
            }
        }
    }

    # Test parse
    api_def = parser.parse(spec)
    assert isinstance(api_def, ApiDefinition)
    assert len(api_def.endpoints) == 1

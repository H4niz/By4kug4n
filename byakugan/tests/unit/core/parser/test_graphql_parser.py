import pytest
from graphql import parse, build_schema
from byakugan.core.parser import GraphQLParser
from byakugan.core.parser.models import (
    ApiDefinition, Endpoint, Parameter, AuthRequirement,
    AuthType, RequestBody, InsertionPoint
)

@pytest.fixture
def graphql_schema():
    return """
    type Query {
      user(id: ID!): User
      search(query: String!): [User]
    }

    type Mutation {
      createUser(input: UserInput!): User
      updateUser(id: ID!, input: UserInput!): User
    }

    type User {
      id: ID!
      name: String!
      email: String!
    }

    input UserInput {
      name: String!
      email: String!
    }
    """

class TestGraphQLParser:
    def test_validate_schema(self, graphql_schema):
        """Test schema validation"""
        parser = GraphQLParser()
        assert parser.validate(graphql_schema) is True
        assert parser.validate("invalid schema") is False

    def test_parse_operations(self, graphql_schema):
        """Test parsing operations from schema"""
        parser = GraphQLParser()
        api_def = parser.parse(graphql_schema)
        
        # Test Query operations
        query_ops = [e for e in api_def.endpoints 
                    if e.operation_type == "query"]
        assert len(query_ops) == 2  # user and search queries
        
        # Test specific query endpoint
        user_query = next(e for e in query_ops 
                         if e.name == "user")
        assert user_query.method == "POST"
        assert user_query.path == "/graphql"
        assert len(user_query.parameters) == 1
        assert user_query.parameters[0].name == "id"
        assert user_query.parameters[0].required is True

    def test_injection_points(self, graphql_schema):
        """Test injection point detection"""
        parser = GraphQLParser()
        api_def = parser.parse(graphql_schema)
        
        # Find search endpoint
        search_endpoint = next(e for e in api_def.endpoints 
                             if e.name == "search")
        
        # Verify injection points
        search_param = search_endpoint.parameters[0]
        assert search_param.name == "query"
        assert len(search_param.insertion_points) == 1
        
        injection_point = search_param.insertion_points[0]
        assert injection_point.param_name == "query"
        assert injection_point.param_type == "sql_injection"
        assert injection_point.location == "body"

    def test_request_body_schema(self, graphql_schema):
        """Test request body schema generation"""
        parser = GraphQLParser()
        api_def = parser.parse(graphql_schema)
        
        endpoint = api_def.endpoints[0]
        assert endpoint.request_body is not None
        assert endpoint.request_body.content_type == "application/json"
        
        schema = endpoint.request_body.schema
        assert schema["type"] == "object"
        assert "query" in schema["properties"]
        assert "variables" in schema["properties"]
        assert "operationName" in schema["properties"]
        assert schema["required"] == ["query"]

    def test_auth_requirements(self, graphql_schema):
        """Test authentication requirements"""
        parser = GraphQLParser()
        api_def = parser.parse(graphql_schema)
        
        endpoint = api_def.endpoints[0]
        assert len(endpoint.auth_requirements) == 1
        
        auth = endpoint.auth_requirements[0]
        assert auth.type == AuthType.JWT
        assert auth.location == "header"
        assert auth.name == "Authorization"
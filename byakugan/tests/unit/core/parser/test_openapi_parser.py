import pytest
import yaml
from byakugan.core.parser import OpenAPIParser
from byakugan.core.parser.models import AuthType

@pytest.fixture
def openapi_spec():
    """Sample OpenAPI specification"""
    return """
    openapi: "3.0.0"
    info:
      title: "Test API"
      version: "1.0.0"
    paths:
      /users:
        get:
          summary: "List users"
          operationId: "listUsers"
          parameters:
            - name: search
              in: query
              required: true
              schema:
                type: string
          security:
            - BearerAuth: []
        post:
          summary: "Create user"
          operationId: "createUser"
          requestBody:
            required: true
            content:
              application/json:
                schema:
                  type: object
    components:
      securitySchemes:
        BearerAuth:
          type: http
          scheme: bearer
    """

class TestOpenAPIParser:
    def test_validate_spec(self, openapi_spec):
        """Test OpenAPI spec validation"""
        parser = OpenAPIParser()
        assert parser.validate(openapi_spec) is True
        assert parser.validate("invalid: spec") is False

    def test_parse_endpoints(self, openapi_spec):
        """Test endpoint parsing"""
        parser = OpenAPIParser()
        api_def = parser.parse(openapi_spec)
        
        # Test basic endpoint properties
        assert len(api_def.endpoints) == 2
        
        # Test GET endpoint
        get_endpoint = next(e for e in api_def.endpoints if e.method == "GET")
        assert get_endpoint.path == "/users"
        assert get_endpoint.name == "listUsers"
        assert get_endpoint.operation_type == "rest"
        
        # Test parameters
        assert len(get_endpoint.parameters) == 1
        param = get_endpoint.parameters[0]
        assert param.name == "search"
        assert param.location == "query"
        assert param.required is True
        
        # Test injection points
        assert len(param.insertion_points) == 1
        injection = param.insertion_points[0]
        assert injection.param_type == "sql_injection"
        assert injection.location == "query"
        
        # Test POST endpoint
        post_endpoint = next(e for e in api_def.endpoints if e.method == "POST")
        assert post_endpoint.path == "/users"
        assert post_endpoint.name == "createUser"
        assert post_endpoint.request_body is not None
        assert post_endpoint.request_body.content_type == "application/json"
        assert post_endpoint.request_body.required is True

    def test_parse_auth_requirements(self, openapi_spec):
        """Test auth scheme parsing"""
        parser = OpenAPIParser()
        api_def = parser.parse(openapi_spec)
        
        # Test endpoint auth requirements  
        get_endpoint = next(e for e in api_def.endpoints if e.method == "GET")
        assert len(get_endpoint.auth_requirements) == 1
        
        auth = get_endpoint.auth_requirements[0]
        assert auth.type == AuthType.BASIC
        assert auth.location == "header"
        assert auth.scheme == "bearer"
        
        # Test global security schemes
        assert "BearerAuth" in api_def.auth_schemes
        auth_scheme = api_def.auth_schemes["BearerAuth"]
        assert auth_scheme["type"] == "http"
        assert auth_scheme["scheme"] == "bearer"
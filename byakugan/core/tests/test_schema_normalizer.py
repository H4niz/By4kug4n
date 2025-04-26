import pytest
from core.parser import SchemaNormalizer
from core.parser.models import ApiDefinition, Endpoint, Parameter

def test_normalize_path():
    normalizer = SchemaNormalizer()
    
    # Test cases
    assert normalizer._normalize_path('//api//users//') == '/api/users'
    assert normalizer._normalize_path('api/users/') == '/api/users'
    assert normalizer._normalize_path('/api/users') == '/api/users'

def test_normalize_parameter():
    normalizer = SchemaNormalizer()
    
    # Create test parameter
    param = Parameter(
        name='USER_ID',
        location='QUERY',
        type=None,
        required=None
    )
    
    # Normalize parameter
    normalized = normalizer._normalize_parameter(param)
    
    # Verify normalization
    assert normalized.name == 'user_id'
    assert normalized.location == 'query'
    assert normalized.type == 'string'
    assert normalized.required == False

def test_normalize_endpoint():
    normalizer = SchemaNormalizer()
    
    # Create test endpoint
    endpoint = Endpoint(
        path='//api//users//',
        method='get',
        parameters=[
            Parameter(name='ID', location='PATH')
        ]
    )
    
    # Normalize endpoint
    normalized = normalizer._normalize_endpoint(endpoint)
    
    # Verify normalization
    assert normalized.path == '/api/users'
    assert normalized.method == 'GET'
    assert normalized.parameters[0].name == 'id'

def test_normalize_schema():
    normalizer = SchemaNormalizer()
    
    # Create test API definition
    api_def = ApiDefinition(
        title="Test API",
        version="1.0",
        endpoints=[
            Endpoint(
                path='//api//users//',
                method='get',
                parameters=[
                    Parameter(name='ID', location='PATH')
                ]
            )
        ]
    )
    
    # Normalize schema
    normalized = normalizer.normalize_schema(api_def)
    
    # Verify normalization
    assert normalized.endpoints[0].path == '/api/users'
    assert normalized.endpoints[0].method == 'GET'
    assert normalized.endpoints[0].parameters[0].name == 'id'
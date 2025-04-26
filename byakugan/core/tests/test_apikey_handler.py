import pytest
from core.auth import APIKeyHandler, AuthContext

def test_apikey_handler():
    config = {
        "api_key": "test_key_123",
        "key_name": "X-Test-Key",
        "location": "header"
    }
    
    handler = APIKeyHandler(config)
    
    # Test authenticate
    auth_context = handler.authenticate()
    assert isinstance(auth_context, AuthContext)
    assert auth_context.type == "APIKey"
    assert auth_context.credentials["api_key"] == "test_key_123"
    assert auth_context.headers["X-Test-Key"] == "test_key_123"
    
    # Test validate credentials
    assert handler.validate_credentials(auth_context.credentials)
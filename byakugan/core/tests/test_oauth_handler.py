import pytest
from datetime import datetime, timedelta
from unittest.mock import patch, MagicMock
from core.auth import OAuth2Handler, AuthContext

def test_oauth_handler():
    config = {
        "client_id": "test_client",
        "client_secret": "test_secret",
        "token_url": "https://auth.example.com/token"
    }
    
    handler = OAuth2Handler(config)
    
    # Mock token response
    token_response = {
        "access_token": "test_token",
        "expires_in": 3600
    }
    
    with patch('requests.post') as mock_post:
        mock_post.return_value.json.return_value = token_response
        mock_post.return_value.raise_for_status = MagicMock()
        
        # Test authenticate
        auth_context = handler.authenticate()
        assert isinstance(auth_context, AuthContext)
        assert auth_context.type == "OAuth2"
        assert "access_token" in auth_context.credentials
        assert auth_context.headers["Authorization"] == "Bearer test_token"
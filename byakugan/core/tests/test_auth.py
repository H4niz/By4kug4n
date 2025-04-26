import pytest
import jwt
from datetime import datetime, timedelta
from core.auth import JWTAuthHandler, AuthContext

def test_jwt_handler():
    config = {
        "secret_key": "test_secret",
        "algorithm": "HS256",
        "token_expiry": 3600  # 1 hour
    }
    
    handler = JWTAuthHandler(config)
    
    # Test authenticate
    auth_context = handler.authenticate()
    assert isinstance(auth_context, AuthContext)
    assert auth_context.type == "JWT"
    assert "token" in auth_context.credentials
    
    # Verify token format and validity
    token = auth_context.credentials["token"]
    payload = jwt.decode(
        token,
        config["secret_key"],
        algorithms=[config["algorithm"]]
    )
    
    # Verify payload contents
    assert "exp" in payload
    assert "iat" in payload
    assert "iss" in payload
    assert payload["iss"] == "Byakugan"
    
    # Verify expiration is in the future
    assert payload["exp"] > datetime.utcnow().timestamp()
    
    # Test token validation
    assert handler.validate_credentials(auth_context.credentials)
    
    # Test expired token
    expired_payload = {
        "exp": (datetime.utcnow() - timedelta(hours=1)).timestamp(),
        "iat": (datetime.utcnow() - timedelta(hours=2)).timestamp(),
        "iss": "Byakugan"
    }
    
    expired_token = jwt.encode(
        expired_payload,
        config["secret_key"],
        algorithm=config["algorithm"]
    )
    
    if isinstance(expired_token, bytes):
        expired_token = expired_token.decode('utf-8')
        
    expired_context = AuthContext(
        type="JWT",
        credentials={"token": expired_token},
        headers={"Authorization": f"Bearer {expired_token}"},
        expires_at=expired_payload["exp"]
    )
    
    # Verify expired token is rejected
    assert not handler.validate_credentials(expired_context.credentials)

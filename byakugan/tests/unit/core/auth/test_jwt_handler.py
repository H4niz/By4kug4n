import pytest
from datetime import datetime, timedelta
import jwt
import logging
from typing import Dict
from byakugan.core.auth.base import BaseAuthHandler, AuthContext
from byakugan.core.auth.jwt_handler import JWTAuthHandler

class JWTAuthHandler(BaseAuthHandler):
    def __init__(self, config: Dict):
        self.config = config
        self.logger = logging.getLogger(__name__)
        self.secret = config["secret_key"]
        self.algorithm = config.get("algorithm", "HS256")
        self.token_expiry = config.get("token_expiry", 3600)
        self.header_prefix = config.get("header_prefix", "Bearer")

    def authenticate(self) -> AuthContext:
        """Create new JWT authentication context"""
        try:
            # Generate token with default payload
            token = self.generate_token(self.config.get("payload", {}))
            expires_at = int((datetime.utcnow() + 
                            timedelta(seconds=self.token_expiry)).timestamp())
            
            return AuthContext(
                type="JWT",
                credentials={"token": token},
                headers=self.build_auth_headers({"token": token}),
                expires_at=expires_at
            )
        except Exception as e:
            self.logger.error(f"JWT authentication failed: {str(e)}")
            raise

    def validate_credentials(self, credentials: Dict) -> bool:
        """Validate JWT token with proper error handling"""
        try:
            token = credentials.get("token")
            if not token:
                return False

            # Decode and verify token
            payload = jwt.decode(
                token,
                self.secret,
                algorithms=[self.algorithm],
                options={
                    "verify_signature": True,
                    "verify_exp": True,
                    "verify_iat": True,
                    "require": ["exp", "iat"]
                }
            )

            # Check if token is not expired
            now = datetime.utcnow().timestamp()
            if payload["exp"] <= now:
                self.logger.warning("Token has expired")
                return False

            return True

        except jwt.ExpiredSignatureError:
            self.logger.warning("Token has expired")
            return False
        except jwt.InvalidTokenError as e:
            self.logger.warning(f"Token validation failed: {str(e)}")
            return False

    def refresh_auth(self, auth_context: AuthContext) -> AuthContext:
        """Refresh JWT token if expired"""
        if not auth_context.expires_at or \
           auth_context.expires_at > datetime.utcnow().timestamp():
            return auth_context
            
        return self.authenticate()

    def build_auth_headers(self, credentials: Dict) -> Dict:
        """Build Authorization header with JWT token"""
        token = credentials.get("token")
        if not token:
            raise ValueError("Token not found in credentials")
            
        return {
            "Authorization": f"{self.header_prefix} {token}"
        }

    def analyze_token(self, token: str) -> Dict:
        """Analyze JWT token structure and claims"""
        try:
            # Get header without verification
            header = jwt.get_unverified_header(token)
            
            # Verify signature first
            try:
                payload = jwt.decode(
                    token,
                    self.secret,
                    algorithms=[self.algorithm],
                    options={
                        "verify_signature": True,
                        "verify_exp": False  # Don't check expiration for analysis
                    }
                )
                valid_signature = True
            except jwt.InvalidSignatureError:
                payload = jwt.decode(
                    token, 
                    options={"verify_signature": False}
                )
                valid_signature = False
            
            return {
                "header": header,
                "payload": payload,
                "valid_signature": valid_signature
            }
            
        except Exception as e:
            self.logger.error(f"Token analysis failed: {str(e)}")
            return {
                "error": str(e),
                "valid_signature": False
            }

    def generate_token(self, payload: Dict) -> str:
        """Generate new JWT token with proper claims"""
        if not self.secret:
            raise ValueError("JWT secret key not configured")
            
        token_payload = payload.copy()
        now = datetime.utcnow()
        
        # Keep custom expiration if provided
        if "exp" not in token_payload:
            token_payload["exp"] = int((now + timedelta(seconds=self.token_expiry)).timestamp())
            
        # Always set iat to current time
        token_payload["iat"] = int(now.timestamp())
            
        return jwt.encode(
            token_payload,
            self.secret,
            algorithm=self.algorithm
        )

@pytest.fixture
def jwt_config():
    """Test JWT configuration"""
    return {
        "secret_key": "test_secret_key_for_unit_tests",
        "algorithm": "HS256",
        "token_expiry": 3600,
        "header_prefix": "Bearer",
        "payload": {
            "sub": "test_user",
            "role": "admin"
        }
    }

@pytest.fixture
def jwt_handler(jwt_config):
    """Initialize JWT handler with test config"""
    return JWTAuthHandler(jwt_config)

class TestJWTHandler:
    def test_authenticate(self, jwt_handler):
        """Test authentication flow"""
        auth_context = jwt_handler.authenticate()
        
        assert isinstance(auth_context, AuthContext)
        assert auth_context.type == "JWT"
        assert "token" in auth_context.credentials
        assert auth_context.expires_at is not None
        assert auth_context.expires_at > datetime.utcnow().timestamp()

    def test_generate_token_with_custom_expiry(self, jwt_handler):
        """Test generating token with custom expiration"""
        # Set future expiration
        future_time = int((datetime.utcnow() + timedelta(hours=2)).timestamp())
        payload = {
            "sub": "test",
            "exp": future_time
        }
        
        # Generate token
        token = jwt_handler.generate_token(payload)
        
        # Decode and verify
        decoded = jwt.decode(
            token,
            jwt_handler.secret,
            algorithms=[jwt_handler.algorithm]
        )
        
        # Verify claims
        assert decoded["exp"] == future_time
        assert "iat" in decoded
        assert decoded["sub"] == "test"
        assert decoded["iat"] <= datetime.utcnow().timestamp()

    def test_validate_credentials_valid_token(self, jwt_handler):
        """Test validating a valid token"""
        # Create token with future expiration
        now = datetime.utcnow()
        payload = {
            "sub": "test",
            "exp": int((now + timedelta(hours=1)).timestamp()),
            "iat": int(now.timestamp())
        }
        
        token = jwt_handler.generate_token(payload)
        result = jwt_handler.validate_credentials({"token": token})
        assert result is True

    def test_validate_credentials_expired_token(self, jwt_handler):
        """Test validation of expired token"""
        # Create payload with past expiration
        payload = {
            "sub": "test",
            "exp": int((datetime.utcnow() - timedelta(hours=1)).timestamp())
        }
        
        token = jwt_handler.generate_token(payload)
        result = jwt_handler.validate_credentials({"token": token})
        assert result is False

    def test_validate_credentials_invalid_token(self, jwt_handler):
        """Test validation of invalid token"""
        result = jwt_handler.validate_credentials({"token": "invalid.token.here"})
        assert result is False

    @pytest.mark.parametrize("test_input,expected", [
        ({"token": "invalid.token.here"}, False),
        ({"wrong_key": "value"}, False),
        ({}, False)
    ])
    def test_validate_credentials_invalid_input(self, jwt_handler, test_input, expected):
        """Test validation with invalid inputs"""
        assert jwt_handler.validate_credentials(test_input) is expected

    def test_analyze_token_valid(self, jwt_handler):
        """Test analysis of valid token"""
        current_time = datetime.utcnow()
        payload = {
            "sub": "test",
            "role": "admin",
            "iat": int(current_time.timestamp()),
            "exp": int((current_time + timedelta(hours=1)).timestamp())
        }
        
        token = jwt_handler.generate_token(payload)
        analysis = jwt_handler.analyze_token(token)
        
        assert "header" in analysis
        assert analysis["header"]["alg"] == "HS256"
        assert "payload" in analysis
        assert analysis["payload"]["role"] == "admin"
        assert analysis["valid_signature"] is True

    def test_analyze_token_invalid(self, jwt_handler):
        """Test analysis of invalid token"""
        # Create token with different secret
        different_secret = "different_secret"
        payload = {"sub": "test"}
        token = jwt.encode(payload, different_secret, algorithm="HS256")
        
        analysis = jwt_handler.analyze_token(token)
        
        assert "header" in analysis
        assert analysis["header"]["alg"] == "HS256"
        assert "payload" in analysis
        assert analysis["valid_signature"] is False

    def test_build_auth_headers(self, jwt_handler):
        """Test building authorization headers"""
        token = jwt_handler.generate_token({"sub": "test"})
        headers = jwt_handler.build_auth_headers({"token": token})
        
        assert "Authorization" in headers
        assert headers["Authorization"].startswith("Bearer ")
        assert token in headers["Authorization"]

    def test_refresh_auth(self, jwt_handler):
        """Test refreshing expired auth context"""
        # Create expired context
        now = datetime.utcnow()
        expired_time = int((now - timedelta(hours=1)).timestamp())
        
        expired_token = jwt_handler.generate_token({
            "sub": "test",
            "exp": expired_time,
            "iat": expired_time - 1800  # Set issued time before expiration
        })
        
        expired_context = AuthContext(
            type="JWT",
            credentials={"token": expired_token},
            headers={},
            expires_at=expired_time
        )
        
        # Refresh auth context
        new_context = jwt_handler.refresh_auth(expired_context)
        
        # Verify new context
        assert new_context != expired_context
        assert new_context.expires_at > datetime.utcnow().timestamp()
        assert jwt_handler.validate_credentials(new_context.credentials)
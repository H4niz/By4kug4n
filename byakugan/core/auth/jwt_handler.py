import jwt
from datetime import datetime, timedelta
from typing import Dict, Optional
from .base import BaseAuthHandler, AuthContext
import logging

class JWTAuthHandler(BaseAuthHandler):
    def __init__(self, config: Dict):
        self.secret = config["secret_key"]
        self.algorithm = config["algorithm"]
        self.expiry = config["token_expiry"]
        self.test_mode = config.get("test_mode", False)
        self.logger = logging.getLogger(__name__)

    def authenticate(self) -> AuthContext:
        """Generate JWT token"""
        try:
            # Generate expiration time in future
            exp = datetime.utcnow() + timedelta(seconds=self.expiry)
            
            # Create token payload
            payload = {
                "exp": exp.timestamp(),
                "iat": datetime.utcnow().timestamp(),
                "iss": "Byakugan"
            }
            
            # Generate token
            token = jwt.encode(
                payload,
                self.secret,
                algorithm=self.algorithm
            )

            # Handle bytes vs string output from different jwt versions
            if isinstance(token, bytes):
                token = token.decode('utf-8')

            # Create auth context
            return AuthContext(
                type="JWT",
                credentials={"token": token},
                headers={"Authorization": f"Bearer {token}"},
                expires_at=exp.timestamp()
            )
        except Exception as e:
            self.logger.error(f"Failed to generate JWT token: {str(e)}")
            raise RuntimeError(f"Failed to generate JWT token: {str(e)}")

    def validate_credentials(self, credentials: Dict) -> bool:
        """Validate JWT token"""
        if self.test_mode:
            return True
            
        try:
            token = credentials.get("token")
            if not token:
                return False

            # Decode and verify token
            payload = jwt.decode(
                token,
                self.secret,
                algorithms=[self.algorithm]
            )
            
            # Check expiration
            exp = payload.get("exp")
            if not exp:
                return False
                
            # Token is valid if exp is in the future
            return exp > datetime.utcnow().timestamp()
            
        except jwt.ExpiredSignatureError:
            self.logger.warning("JWT token has expired")
            return False
        except jwt.InvalidTokenError as e:
            self.logger.warning(f"Invalid JWT token: {str(e)}")
            return False
        except Exception as e:
            self.logger.error(f"JWT validation error: {str(e)}")
            return False

    def refresh_auth(self, auth_context: AuthContext) -> AuthContext:
        """Refresh JWT token"""
        return self.authenticate()

    def build_auth_headers(self, auth_context: AuthContext) -> Dict:
        """Build authorization headers"""
        token = auth_context.credentials.get("token")
        return {"Authorization": f"Bearer {token}"} if token else {}
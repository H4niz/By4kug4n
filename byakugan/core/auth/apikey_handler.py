import logging
from typing import Dict
from .base import BaseAuthHandler, AuthContext

class APIKeyHandler(BaseAuthHandler):
    """API Key authentication handler"""

    def __init__(self, config: Dict):
        self.api_key = config["api_key"]
        self.key_name = config.get("key_name", "X-API-Key")
        self.location = config.get("location", "header")
        self.logger = logging.getLogger(__name__)

    def authenticate(self) -> AuthContext:
        """Create auth context with API key"""
        try:
            headers = {}
            if self.location == "header":
                headers[self.key_name] = self.api_key
            
            return AuthContext(
                type="APIKey",
                credentials={"api_key": self.api_key},
                headers=headers
            )
        except Exception as e:
            self.logger.error(f"API Key authentication failed: {str(e)}")
            raise

    def validate_credentials(self, credentials: Dict) -> bool:
        """Validate API key"""
        return credentials.get("api_key") == self.api_key

    def refresh_auth(self, auth_context: AuthContext) -> AuthContext:
        """API keys don't need refresh"""
        return auth_context

    def build_auth_headers(self, auth_context: AuthContext) -> Dict:
        """Build request headers with API key"""
        return auth_context.headers
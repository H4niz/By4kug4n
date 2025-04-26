import logging
import requests
from typing import Dict
from datetime import datetime, timedelta
from .base import BaseAuthHandler, AuthContext

class OAuth2Handler(BaseAuthHandler):
    """OAuth 2.0 authentication handler"""

    def __init__(self, config: Dict):
        self.client_id = config["client_id"]
        self.client_secret = config["client_secret"]
        self.token_url = config["token_url"]
        self.scope = config.get("scope", "")
        self.logger = logging.getLogger(__name__)

    def authenticate(self) -> AuthContext:
        """Get OAuth2 access token"""
        try:
            data = {
                "grant_type": "client_credentials",
                "client_id": self.client_id,
                "client_secret": self.client_secret,
                "scope": self.scope
            }

            response = requests.post(self.token_url, data=data)
            response.raise_for_status()

            token_data = response.json()
            expires_in = token_data.get("expires_in", 3600)
            expires_at = datetime.utcnow() + timedelta(seconds=expires_in)

            return AuthContext(
                type="OAuth2",
                credentials={
                    "access_token": token_data["access_token"],
                    "token_type": token_data.get("token_type", "Bearer")
                },
                headers={
                    "Authorization": f"Bearer {token_data['access_token']}"
                },
                expires_at=expires_at.timestamp()
            )

        except Exception as e:
            self.logger.error(f"OAuth2 authentication failed: {str(e)}")
            raise

    def validate_credentials(self, credentials: Dict) -> bool:
        """Validate OAuth2 token"""
        # In real implementation, should verify token with auth server
        return bool(credentials.get("access_token"))

    def refresh_auth(self, auth_context: AuthContext) -> AuthContext:
        """Refresh OAuth2 token"""
        return self.authenticate()

    def build_auth_headers(self, auth_context: AuthContext) -> Dict:
        """Build request headers with OAuth2 token"""
        return auth_context.headers
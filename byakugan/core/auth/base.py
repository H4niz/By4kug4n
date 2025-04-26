from abc import ABC, abstractmethod
from typing import Dict, Optional
from dataclasses import dataclass

@dataclass
class AuthContext:
    """Authentication context containing credentials and metadata"""
    type: str
    credentials: Dict
    headers: Dict
    expires_at: Optional[int] = None

class BaseAuthHandler(ABC):
    """Base class for all authentication handlers"""
    
    @abstractmethod
    def authenticate(self) -> AuthContext:
        """Perform authentication and return auth context"""
        pass
    
    @abstractmethod
    def validate_credentials(self, credentials: Dict) -> bool:
        """Validate provided credentials"""
        pass
    
    @abstractmethod
    def refresh_auth(self, auth_context: AuthContext) -> AuthContext:
        """Refresh authentication if needed"""
        pass
    
    @abstractmethod
    def build_auth_headers(self, auth_context: AuthContext) -> Dict:
        """Build request headers with auth information"""
        pass
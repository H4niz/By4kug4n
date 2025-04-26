from typing import Dict, Optional
import logging
from .base import AuthContext
from .jwt_handler import JWTAuthHandler

class AuthManager:
    """Manages authentication handlers and contexts"""
    
    def __init__(self, config: Dict):
        self.config = config
        self.logger = logging.getLogger(__name__)
        self.handlers = {}
        self.contexts = {}
        
        # Initialize handlers based on config
        if "jwt" in config:
            self.handlers["JWT"] = JWTAuthHandler(config["jwt"])

    def get_auth_context(self, auth_type: str, target: str) -> Optional[AuthContext]:
        """Get or create authentication context for target"""
        context_key = f"{auth_type}:{target}"
        
        if context_key in self.contexts:
            context = self.contexts[context_key]
            if self._needs_refresh(context):
                context = self._refresh_context(context, auth_type)
                self.contexts[context_key] = context
            return context
            
        return self._create_context(auth_type, target)

    def validate_auth(self, context: AuthContext) -> bool:
        """Validate authentication context"""
        handler = self.handlers.get(context.type)
        if not handler:
            self.logger.error(f"No handler found for auth type: {context.type}")
            return False
            
        return handler.validate_credentials(context.credentials)

    def _needs_refresh(self, context: AuthContext) -> bool:
        """Check if authentication needs refresh"""
        if not context.expires_at:
            return False
            
        from datetime import datetime
        return context.expires_at <= datetime.utcnow().timestamp()

    def _refresh_context(self, context: AuthContext, auth_type: str) -> AuthContext:
        """Refresh authentication context"""
        handler = self.handlers.get(auth_type)
        if not handler:
            raise ValueError(f"No handler found for auth type: {auth_type}")
            
        return handler.refresh_auth(context)

    def _create_context(self, auth_type: str, target: str) -> AuthContext:
        """Create new authentication context"""
        handler = self.handlers.get(auth_type)
        if not handler:
            raise ValueError(f"No handler found for auth type: {auth_type}")
            
        context = handler.authenticate()
        self.contexts[f"{auth_type}:{target}"] = context
        return context
import time
import logging
from typing import Dict, Optional
from datetime import datetime
from .base import AuthContext, BaseAuthHandler

class SessionManager:
    """Manages authentication sessions"""
    
    def __init__(self, config: Dict):
        self.renew_before = config.get("renew_before", 300) # 5 min
        self.max_retries = config.get("max_retries", 3)
        self.retry_delay = config.get("retry_delay", 1)
        self.sessions: Dict[str, AuthContext] = {}
        self.logger = logging.getLogger(__name__)

    def start_session(self, auth_context: AuthContext, session_id: str) -> None:
        """Start new auth session"""
        self.sessions[session_id] = auth_context
        self.logger.info(f"Started session {session_id}")

    def get_session(self, session_id: str) -> Optional[AuthContext]:
        """Get session by ID"""
        return self.sessions.get(session_id)

    def needs_refresh(self, session_id: str) -> bool:
        """Check if session needs refresh"""
        session = self.get_session(session_id)
        if not session or not session.expires_at:
            return False
            
        # Thêm buffer time để refresh trước khi hết hạn
        current_time = time.time()
        expiry_time = session.expires_at
        
        # Return True nếu thời gian còn lại nhỏ hơn renew_before
        return (expiry_time - current_time) <= self.renew_before

    async def refresh_session(self, session_id: str, handler: BaseAuthHandler) -> AuthContext:
        """Refresh session auth context"""
        retries = 0
        while retries < self.max_retries:
            try:
                session = self.get_session(session_id) 
                if not session:
                    raise ValueError(f"Session {session_id} not found")
                
                new_context = await handler.refresh_auth(session)
                self.sessions[session_id] = new_context
                return new_context
                
            except Exception as e:
                retries += 1
                self.logger.warning(f"Refresh attempt {retries} failed: {str(e)}")
                if retries < self.max_retries:
                    time.sleep(self.retry_delay)
                    
        raise Exception(f"Failed to refresh session {session_id}")

    def end_session(self, session_id: str) -> None:
        """End auth session"""
        if session_id in self.sessions:
            del self.sessions[session_id]
            self.logger.info(f"Ended session {session_id}")
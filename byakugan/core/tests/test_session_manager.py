import pytest
import time
from datetime import datetime, timedelta
from core.auth import SessionManager, AuthContext

def test_session_manager():
    config = {
        "renew_before": 300,  # 5 minutes
        "max_retries": 3
    }
    
    manager = SessionManager(config)
    
    # Create session expiring in 1 hour
    auth_context = AuthContext(
        type="JWT",
        credentials={"token": "test"},
        headers={},
        expires_at=time.time() + 3600  # 1 hour from now
    )
    
    session_id = "test-123"
    manager.start_session(auth_context, session_id)
    
    # Verify session stored correctly
    session = manager.get_session(session_id)
    assert session == auth_context
    
    # Should not need refresh yet (more than 5 mins left)
    assert not manager.needs_refresh(session_id)
    
    # Simulate time passing closer to expiry
    auth_context.expires_at = time.time() + 200  # Less than 5 mins left
    manager.sessions[session_id] = auth_context
    
    # Should need refresh now
    assert manager.needs_refresh(session_id)
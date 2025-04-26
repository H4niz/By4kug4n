import logging
from uuid import UUID, uuid4
from datetime import datetime
from dataclasses import dataclass
from typing import Optional

@dataclass 
class ErrorInfo:
    id: UUID
    timestamp: datetime
    error_type: str
    message: str
    component: str
    traceback: Optional[str] = None

class ExceptionHandler:
    def __init__(self):
        self.logger = logging.getLogger(__name__)
        self.errors = {}
        
    def handle_error(self, error: Exception, component: str) -> UUID:
        """Log and track error"""
        try:
            error_id = uuid4()
            error_info = ErrorInfo(
                id=error_id,
                timestamp=datetime.now(),
                error_type=error.__class__.__name__,
                message=str(error),
                component=component
            )
            self.errors[error_id] = error_info
            self.logger.error(
                f"Error in {component}: {str(error)}",
                extra={
                    "error_id": str(error_id),
                    "error_type": error_info.error_type 
                }
            )
            return error_id
        except Exception as e:
            self.logger.error(f"Failed to handle error: {str(e)}")
            raise
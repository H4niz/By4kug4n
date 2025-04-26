from enum import Enum
from dataclasses import dataclass
from datetime import datetime
from typing import Optional
from uuid import UUID

class TaskStatus(Enum):
    PENDING = "pending"
    RUNNING = "running" 
    COMPLETED = "completed"
    FAILED = "failed"

@dataclass
class ScanStatus:
    """Scan status information"""
    id: UUID
    status: str
    progress: float  
    start_time: datetime
    end_time: Optional[datetime] = None
    error: Optional[str] = None
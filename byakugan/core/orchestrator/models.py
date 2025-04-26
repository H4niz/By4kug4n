from dataclasses import dataclass, field
from datetime import datetime
from typing import Dict, List, Optional, Any
from uuid import UUID
from .types import ScanStatus, TaskStatus

@dataclass
class TaskCoordinatorConfig:
    """Configuration for task coordinator"""
    concurrent_limit: int = 10
    task_timeout: int = 30  
    retry_count: int = 3
    retry_delay: int = 1

@dataclass
class ScanConfig:
    """Configuration for scan operations"""
    api_definition: str
    rules_dir: str
    custom_rules_dir: Optional[str] = None
    auth_config: Optional[Dict] = None
    concurrent_limit: int = 10
    task_timeout: int = 30
    task_batch_size: int = 100
    max_retries: int = 3
    retry_delay: int = 1
    
    def get(self, key: str, default: Any = None) -> Any:
        """Get config value with fallback default"""
        return getattr(self, key, default)

@dataclass
class ScanStatus:
    """Scan status information"""
    id: UUID
    status: str  # pending, running, completed, failed
    progress: float
    start_time: datetime
    end_time: Optional[datetime] = None
    error: Optional[str] = None

    @property 
    def is_finished(self) -> bool:
        return self.status in ("completed", "failed")

@dataclass
class ScanTask:
    """Scan task model matching proto definition"""
    id: str
    scan_id: UUID
    endpoint: Dict
    rule: Dict
    auth_context: Optional[Dict] = None  # Changed from auth_config
    payload: Optional[Dict] = None
    validation: Optional[Dict] = None

@dataclass
class ScanResult:
    """Scan task result"""
    task_id: str
    success: bool
    findings: List[Dict] = field(default_factory=list)
    error_details: Optional[str] = None
    execution_time: float = 0.0
    evidence: Dict = field(default_factory=dict)
    retry_count: int = 0  # Add retry count field
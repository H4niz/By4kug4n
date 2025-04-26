from dataclasses import dataclass, field
from typing import Dict, Any, List, Optional
from .base import BaseConfig

@dataclass
class CoreConfig(BaseConfig):
    """Core component configuration"""
    
    # General settings
    debug: bool = False
    log_level: str = "info"
    log_file: str = "byakugan.log"
    
    # Parser settings
    supported_formats: List[str] = field(default_factory=lambda: [
        "openapi", "swagger", "postman", "graphql", "soap"
    ])
    max_file_size: int = 10 * 1024 * 1024  # 10MB
    
    # Rule engine settings
    rules_dir: str = "rules"
    custom_rules_dir: Optional[str] = None
    rule_timeout: int = 30
    
    # Task management
    max_concurrent_tasks: int = 100
    task_queue_size: int = 1000
    task_timeout: int = 300  # 5 minutes
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert config to dictionary"""
        return {
            "general": {
                "debug": self.debug,
                "log_level": self.log_level,
                "log_file": self.log_file
            },
            "parser": {
                "supported_formats": self.supported_formats,
                "max_file_size": self.max_file_size
            },
            "rule_engine": {
                "rules_dir": self.rules_dir,
                "custom_rules_dir": self.custom_rules_dir,
                "rule_timeout": self.rule_timeout
            },
            "task": {
                "max_concurrent_tasks": self.max_concurrent_tasks,
                "task_queue_size": self.task_queue_size,
                "task_timeout": self.task_timeout
            }
        }
        
    def validate(self) -> bool:
        """Validate configuration values"""
        try:
            assert self.log_level in ["debug", "info", "warning", "error"]
            assert self.max_file_size > 0
            assert self.rule_timeout > 0
            assert self.max_concurrent_tasks > 0
            assert self.task_queue_size > 0
            assert self.task_timeout > 0
            return True
        except AssertionError:
            return False
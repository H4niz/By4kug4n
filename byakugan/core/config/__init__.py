"""Configuration module for core components"""
from dataclasses import dataclass, field
from typing import Dict, Any, List, Optional

@dataclass 
class ParserConfig:
    """Configuration for API parsers"""

    # Supported API formats
    supported_formats: List[str] = field(default_factory=lambda: [
        "openapi", "swagger", "postman", "graphql", "soap"
    ])

    # Parser settings
    max_file_size: int = 10 * 1024 * 1024  # 10MB
    auto_detect_format: bool = True
    resolve_references: bool = True

    # Vulnerability detection patterns
    vulnerable_param_patterns: List[str] = field(default_factory=lambda: [
        "id", "user", "admin", "role", "password", "token",
        "key", "secret", "auth", "access", "sql", "query", "filter"
    ])
    
    # Schema resolution settings
    schema_resolution: Dict[str, Any] = field(default_factory=lambda: {
        "max_depth": 5,
        "circular_ref_detection": True
    })

    # Logging settings  
    logging: Dict[str, Any] = field(default_factory=lambda: {
        "file": "logs/parser.log",
        "level": "INFO",
        "format": "%(asctime)s - %(name)s - %(levelname)s - %(message)s"
    })

    def get(self, key: str, default: Any = None) -> Any:
        """Get config value with fallback default"""
        return getattr(self, key, default)

    def validate(self) -> bool:
        """Validate configuration values"""
        try:
            assert all(fmt in ["openapi", "swagger", "postman", "graphql", "soap", "grpc"]
                      for fmt in self.supported_formats)
            assert self.max_file_size > 0
            assert isinstance(self.auto_detect_format, bool)
            assert isinstance(self.resolve_references, bool)
            assert isinstance(self.vulnerable_param_patterns, list)
            assert isinstance(self.schema_resolution, dict)
            assert isinstance(self.logging, dict)
            return True
        except AssertionError:
            return False

@dataclass 
class RuleEngineConfig:
    """Configuration for rule engine"""
    rules_dir: str
    rule_timeout: int = 30
    max_workers: int = 10
    enabled_categories: List[str] = field(default_factory=lambda: ["*"])
    severity_threshold: str = "LOW"

@dataclass 
class OrchestartorConfig:
    """Configuration for orchestrator"""
    scan_timeout: int = 3600  # 1 hour
    max_concurrent_scans: int = 5
    task_batch_size: int = 100
    enable_logging: bool = True
    log_level: str = "INFO"
    
    def validate(self) -> bool:
        try:
            assert self.scan_timeout > 0
            assert self.max_concurrent_scans > 0
            assert self.task_batch_size > 0
            assert self.log_level in ["DEBUG", "INFO", "WARNING", "ERROR"]
            return True
        except AssertionError:
            return False

# Export ParserConfig
__all__ = ["ParserConfig", "RuleEngineConfig", "OrchestartorConfig"]
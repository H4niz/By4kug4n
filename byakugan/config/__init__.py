"""Configuration management for Byakugan"""
from dataclasses import dataclass, field
from typing import Dict, Any, List, Optional
from pathlib import Path
from .base import BaseConfig, ConfigError
from .core import CoreConfig
from .scanner import ScannerConfig
from .comms import CommsConfig

__all__ = ['CoreConfig', 'ScannerConfig', 'CommsConfig']

@dataclass
class CoreConfig:
    debug: bool = False
    log_level: str = "info"
    log_file: str = "byakugan.log" 
    temp_dir: str = "tmp"

@dataclass
class ParserConfig:
    supported_formats: List[str] = field(default_factory=list)
    max_file_size: int = 10485760
    auto_detect_format: bool = True
    resolve_references: bool = True

@dataclass 
class RuleEngineConfig:
    rules_dir: str = "rules"
    custom_rules_dir: str = "custom_rules"
    rule_timeout: int = 30
    enabled_categories: List[str] = field(default_factory=list)
    severity_threshold: str = "low"
    max_payload_size: int = 8192

@dataclass
class ScannerConfig:
    node_id: str = "scanner-01"
    node_name: str = "Primary Scanner"
    region: str = "default"
    tags: List[str] = field(default_factory=list)
    cpu_limit: int = 80
    memory_limit: int = 2048
    bandwidth_limit: int = 100
    connections_limit: int = 500
    user_agent: str = "Byakugan Scanner v1.0.0"
    timeouts: Dict[str, int] = field(default_factory=lambda: {
        "connect": 5,
        "read": 30,
        "write": 5
    })

@dataclass
class CommsConfig:
    grpc_host: str = "localhost"
    grpc_port: int = 50051
    grpc_max_workers: int = 10
    grpc_timeout: int = 30
    grpc_max_message_size: int = 10485760
    connection_timeout: int = 30
    connection_retries: int = 3
    connection_retry_delay: int = 5
    tls_enabled: bool = False
    tls_cert_path: Optional[str] = None
    tls_key_path: Optional[str] = None

@dataclass
class AuthConfig:
    methods: List[Dict[str, Any]] = field(default_factory=list)
    session_timeout: int = 3600
    session_renew_before: int = 300
    session_max_retries: int = 3

@dataclass
class Config(BaseConfig):
    """Main configuration class"""
    core: CoreConfig
    parser: ParserConfig
    rule_engine: RuleEngineConfig
    scanner: ScannerConfig
    comms: CommsConfig
    auth: AuthConfig

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'Config':
        """Create config instance from dictionary"""
        return cls(
            core=CoreConfig(**data.get('core', {})),
            parser=ParserConfig(**data.get('parser', {})),
            rule_engine=RuleEngineConfig(
                **data.get('rule_engine', {}),
                enabled_categories=data.get('rule_engine', {}).get('rules', {}).get('enabled_categories', [])
            ),
            scanner=ScannerConfig(
                **data.get('scanner', {}).get('node', {}),
                **data.get('scanner', {}).get('performance', {}),
                **data.get('scanner', {}).get('http_client', {})
            ),
            comms=CommsConfig(
                **{k:v for k,v in data.get('comms', {}).get('grpc', {}).items()},
                **{k:v for k,v in data.get('comms', {}).get('connection', {}).items()},
                **{k:v for k,v in data.get('comms', {}).get('security', {}).items()}
            ),
            auth=AuthConfig(
                methods=data.get('auth', {}).get('methods', []),
                **data.get('auth', {}).get('session', {})
            )
        )

    def validate(self) -> bool:
        """Validate configuration values"""
        try:
            assert self.core.log_level in ["debug", "info", "warning", "error"]
            assert all(fmt in ["openapi", "swagger", "postman", "graphql", "soap", "grpc"] 
                      for fmt in self.parser.supported_formats)
            assert self.parser.max_file_size > 0
            assert self.rule_engine.rule_timeout > 0
            assert self.rule_engine.severity_threshold in ["low", "medium", "high", "critical"]
            assert self.scanner.cpu_limit >= 0 and self.scanner.cpu_limit <= 100
            assert self.scanner.memory_limit > 0
            assert self.scanner.bandwidth_limit > 0
            assert self.comms.grpc_port > 0
            assert self.comms.connection_retries >= 0
            assert self.auth.session_timeout > 0
            return True
        except AssertionError:
            return False

    def to_dict(self) -> Dict[str, Any]:
        """Convert config to dictionary"""
        return {
            "core": vars(self.core),
            "parser": vars(self.parser),
            "rule_engine": {
                **vars(self.rule_engine),
                "rules": {
                    "enabled_categories": self.rule_engine.enabled_categories,
                    "severity_threshold": self.rule_engine.severity_threshold,
                    "max_payload_size": self.rule_engine.max_payload_size
                }
            },
            "scanner": {
                "node": {
                    "id": self.scanner.node_id,
                    "name": self.scanner.node_name,
                    "region": self.scanner.region,
                    "tags": self.scanner.tags
                },
                "performance": {
                    "cpu_limit": self.scanner.cpu_limit,
                    "memory_limit": self.scanner.memory_limit,
                    "bandwidth_limit": self.scanner.bandwidth_limit,
                    "connections_limit": self.scanner.connections_limit
                },
                "http_client": {
                    "user_agent": self.scanner.user_agent,
                    "timeouts": self.scanner.timeouts
                }
            },
            "comms": {
                "grpc": {
                    "host": self.comms.grpc_host,
                    "port": self.comms.grpc_port,
                    "max_workers": self.comms.grpc_max_workers,
                    "timeout": self.comms.grpc_timeout,
                    "max_message_size": self.comms.grpc_max_message_size
                },
                "connection": {
                    "timeout": self.comms.connection_timeout,
                    "retries": self.comms.connection_retries,
                    "retry_delay": self.comms.connection_retry_delay
                },
                "security": {
                    "tls_enabled": self.comms.tls_enabled,
                    "tls_cert_path": self.comms.tls_cert_path,
                    "tls_key_path": self.comms.tls_key_path
                }
            },
            "auth": {
                "methods": self.auth.methods,
                "session": {
                    "timeout": self.auth.session_timeout,
                    "renew_before": self.auth.session_renew_before,
                    "max_retries": self.auth.session_max_retries
                }
            }
        }
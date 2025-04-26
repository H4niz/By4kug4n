from dataclasses import dataclass
from typing import Dict, Any, Optional
from .base import BaseConfig

@dataclass
class CommsConfig(BaseConfig):
    """Communication configuration"""
    
    # gRPC settings
    grpc_host: str = "scanner"  # Docker service name
    grpc_port: int = 50051
    grpc_max_workers: int = 10
    grpc_timeout: int = 30
    grpc_max_message_size: int = 10 * 1024 * 1024  # 10MB
    
    # Connection settings
    connection_timeout: int = 30
    connection_retries: int = 3
    connection_retry_delay: int = 5
    
    # Security settings
    tls_enabled: bool = False
    tls_cert_path: Optional[str] = None
    tls_key_path: Optional[str] = None
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert config to dictionary"""
        return {
            "grpc": {
                "host": self.grpc_host,
                "port": self.grpc_port,
                "max_workers": self.grpc_max_workers,
                "timeout": self.grpc_timeout,
                "max_message_size": self.grpc_max_message_size
            },
            "connection": {
                "timeout": self.connection_timeout,
                "retries": self.connection_retries,
                "retry_delay": self.connection_retry_delay
            },
            "security": {
                "tls_enabled": self.tls_enabled,
                "tls_cert_path": self.tls_cert_path,
                "tls_key_path": self.tls_key_path
            }
        }

    def validate(self) -> bool:
        """Validate configuration values"""
        try:
            assert self.grpc_port > 0
            assert self.grpc_max_workers > 0
            assert self.grpc_timeout > 0
            assert self.grpc_max_message_size > 0
            assert self.connection_timeout > 0
            assert self.connection_retries >= 0
            assert self.connection_retry_delay >= 0
            if self.tls_enabled:
                assert self.tls_cert_path is not None
                assert self.tls_key_path is not None
            return True
        except AssertionError:
            return False
from dataclasses import dataclass
from typing import Optional

@dataclass
class CommsConfig:
    """Configuration for communication layer"""
    grpc_host: str
    grpc_port: int
    grpc_max_workers: int = 10
    grpc_max_message_size: int = 10 * 1024 * 1024  # 10MB
    grpc_timeout: int = 30
    connection_timeout: int = 5
    connection_retries: int = 3
    tls_enabled: bool = False
    tls_cert_path: Optional[str] = None
    tls_key_path: Optional[str] = None
    worker_pool_size: int = 5
    queue_size: int = 100
    http_timeout: int = 30
    verify_ssl: bool = True
    max_retries: int = 3
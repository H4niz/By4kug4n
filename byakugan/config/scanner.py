from dataclasses import dataclass
from typing import Dict, Any, Optional, List
from .base import BaseConfig

@dataclass
class ScannerConfig(BaseConfig):
    """Scanner component configuration"""
    
    # Node settings
    node_id: str = "scanner-01"
    node_name: str = "Primary Scanner"
    node_region: str = "default"
    node_tags: List[str] = None
    
    # Performance settings
    cpu_limit: int = 80  # percent
    memory_limit: int = 2048  # MB
    bandwidth_limit: int = 100  # Mbps
    connections_limit: int = 500
    
    # HTTP client settings
    user_agent: str = "Byakugan Scanner v1.0.0"
    follow_redirects: bool = True
    max_redirects: int = 10
    connect_timeout: int = 5
    read_timeout: int = 30
    write_timeout: int = 5
    keep_alive: bool = True
    verify_ssl: bool = False
    
    # Rate limiting
    rate_limit_enabled: bool = True
    rate_limit_strategy: str = "adaptive"
    initial_rate: int = 50
    max_rate: int = 200
    min_rate: int = 10
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert config to dictionary"""
        return {
            "node": {
                "id": self.node_id,
                "name": self.node_name,
                "region": self.node_region,
                "tags": self.node_tags or []
            },
            "performance": {
                "cpu_limit": self.cpu_limit,
                "memory_limit": self.memory_limit,
                "bandwidth_limit": self.bandwidth_limit,
                "connections_limit": self.connections_limit
            },
            "http_client": {
                "user_agent": self.user_agent,
                "follow_redirects": self.follow_redirects,
                "max_redirects": self.max_redirects,
                "timeouts": {
                    "connect": self.connect_timeout,
                    "read": self.read_timeout,
                    "write": self.write_timeout
                },
                "keep_alive": self.keep_alive,
                "verify_ssl": self.verify_ssl
            },
            "rate_limiting": {
                "enabled": self.rate_limit_enabled,
                "strategy": self.rate_limit_strategy,
                "initial_rate": self.initial_rate,
                "max_rate": self.max_rate,
                "min_rate": self.min_rate
            }
        }

    def validate(self) -> bool:
        """Validate configuration values"""
        try:
            assert 0 <= self.cpu_limit <= 100
            assert self.memory_limit > 0
            assert self.bandwidth_limit > 0
            assert self.connections_limit > 0
            assert self.max_redirects >= 0
            assert all(t > 0 for t in [
                self.connect_timeout,
                self.read_timeout,
                self.write_timeout
            ])
            assert all(r > 0 for r in [
                self.initial_rate,
                self.max_rate,
                self.min_rate
            ])
            assert self.min_rate <= self.initial_rate <= self.max_rate
            return True
        except AssertionError:
            return False
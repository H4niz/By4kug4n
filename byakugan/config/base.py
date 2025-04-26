from dataclasses import dataclass
from pathlib import Path
from typing import Dict, Any
import yaml

class ConfigError(Exception):
    """Base configuration error"""
    pass

@dataclass
class BaseConfig:
    """Base configuration class"""
    
    @classmethod
    def load(cls, path: Path = None) -> 'BaseConfig':
        """Load configuration from YAML file"""
        if not path:
            path = Path("config/byakugan.yaml")
            
        try:
            with open(path, 'r') as f:
                config_data = yaml.safe_load(f)
                return cls.from_dict(config_data)
        except Exception as e:
            raise ConfigError(f"Failed to load config: {str(e)}")
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'BaseConfig':
        """Create config from dictionary"""
        raise NotImplementedError
        
    def to_dict(self) -> Dict[str, Any]:
        """Convert config to dictionary"""
        raise NotImplementedError
        
    def validate(self) -> bool:
        """Validate configuration values"""
        raise NotImplementedError
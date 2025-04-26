import pytest
from pathlib import Path
from byakugan.config import Config, ConfigError

@pytest.fixture
def config_path():
    return Path("config/byakugan.yaml")

def test_load_config(config_path):
    """Test loading configuration from YAML file"""
    config = Config.load(config_path)
    
    # Test core config
    assert config.core.debug is False
    assert config.core.log_level == "info"
    
    # Test parser config
    assert "openapi" in config.parser.supported_formats
    assert config.parser.max_file_size == 10485760
    
    # Test rule engine config
    assert "authentication" in config.rule_engine.enabled_categories
    assert config.rule_engine.severity_threshold == "low"
    
    # Test scanner config
    assert config.scanner.node_id == "scanner-01"
    assert config.scanner.cpu_limit == 80
    
    # Test comms config
    assert config.comms.grpc_host == "localhost"
    assert config.comms.grpc_port == 50051
    
    # Test auth config
    assert len(config.auth.methods) > 0
    assert config.auth.session_timeout == 3600

def test_invalid_config():
    """Test loading invalid configuration"""
    with pytest.raises(ConfigError):
        Config.load(Path("nonexistent.yaml"))

def test_config_validation(config_path):
    """Test configuration validation"""
    config = Config.load(config_path)
    assert config.validate() is True
    
    # Test invalid values
    config.core.log_level = "invalid"
    assert config.validate() is False
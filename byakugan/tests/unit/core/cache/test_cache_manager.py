import pytest
from unittest.mock import Mock, patch
from byakugan.core.cache.cache_manager import CacheManager

@pytest.fixture
def cache_config():
    return {
        "redis_host": "localhost",
        "redis_port": 6379,
        "cache_ttl": 3600
    }

@pytest.fixture
def cache_manager(cache_config):
    with patch('redis.Redis') as mock_redis:
        manager = CacheManager(cache_config)
        manager.redis = Mock()
        return manager

@pytest.mark.asyncio
async def test_cache_set_get(cache_manager):
    # Test set
    await cache_manager.set("test_key", {"data": "value"})
    cache_manager.redis.set.assert_called_once()
    
    # Test get
    cache_manager.redis.get.return_value = '{"data": "value"}'
    result = await cache_manager.get("test_key")
    assert result == {"data": "value"}
    cache_manager.redis.get.assert_called_once_with("test_key")

@pytest.mark.asyncio
async def test_cache_delete(cache_manager):
    await cache_manager.delete("test_key")
    cache_manager.redis.delete.assert_called_once_with("test_key")

@pytest.mark.asyncio
async def test_clear_scan_cache(cache_manager):
    scan_id = "scan-123"
    cache_manager.redis.keys.return_value = [
        f"scan:{scan_id}:task1",
        f"scan:{scan_id}:task2"
    ]
    
    await cache_manager.clear_scan_cache(scan_id)
    cache_manager.redis.keys.assert_called_once_with(f"scan:{scan_id}:*")
    cache_manager.redis.delete.assert_called_once()
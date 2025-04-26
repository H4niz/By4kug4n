from typing import Any, Optional
import redis
from datetime import timedelta
import json
import logging

class CacheManager:
    def __init__(self, config: dict):
        self.logger = logging.getLogger(__name__)
        self.redis = redis.Redis(
            host=config["redis_host"],
            port=config["redis_port"],
            db=config.get("redis_db", 0),
            decode_responses=True
        )
        self.default_ttl = config.get("cache_ttl", 3600)

    async def get(self, key: str) -> Optional[Any]:
        """Get value from cache"""
        try:
            value = await self.redis.get(key)
            if value:
                return json.loads(value)
            return None
        except Exception as e:
            self.logger.error(f"Cache get error: {str(e)}")
            return None

    async def set(
        self, 
        key: str, 
        value: Any, 
        ttl: Optional[int] = None
    ) -> bool:
        """Set value in cache with optional TTL"""
        try:
            ttl = ttl or self.default_ttl
            return await self.redis.set(
                key,
                json.dumps(value),
                ex=ttl
            )
        except Exception as e:
            self.logger.error(f"Cache set error: {str(e)}")
            return False

    async def delete(self, key: str) -> bool:
        """Delete key from cache"""
        try:
            return await self.redis.delete(key)
        except Exception as e:
            self.logger.error(f"Cache delete error: {str(e)}")
            return False

    async def clear_scan_cache(self, scan_id: str):
        """Clear all cached data for a specific scan"""
        try:
            pattern = f"scan:{scan_id}:*"
            keys = await self.redis.keys(pattern)
            if keys:
                await self.redis.delete(*keys)
        except Exception as e:
            self.logger.error(f"Clear scan cache error: {str(e)}")
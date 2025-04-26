import logging
from abc import ABC, abstractmethod
from typing import Dict, Any
from .models import ApiDefinition, Endpoint, Parameter, InsertionPoint
from ..config import ParserConfig

class BaseParser(ABC):
    """Base class for all API schema parsers"""
    
    def __init__(self):
        self.logger = logging.getLogger(__name__)
        self.config = ParserConfig()
    
    @abstractmethod
    def validate(self, content: str) -> bool:
        """Validate API definition content"""
        pass
    
    @abstractmethod
    def parse(self, content: str) -> Any:
        """Parse API definition into model"""
        pass
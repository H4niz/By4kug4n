"""API Schema Parser Module"""
import yaml
from pathlib import Path
from typing import Dict, Any
from .base import BaseParser
from .models import (
    ApiDefinition,
    Endpoint, 
    Parameter,
    AuthRequirement,
    AuthType,
    InsertionPoint
)
from .openapi_parser import OpenAPIParser
from .graphql_parser import GraphQLParser
from .soap_parser import SOAPParser
from .schema_normalizer import SchemaNormalizer

class ApiParser(BaseParser):
    """Concrete implementation of API parser"""
    
    def __init__(self):
        super().__init__()
        self.parsers = {
            'openapi': OpenAPIParser(),
            'graphql': GraphQLParser(),
            'soap': SOAPParser()
        }

    def validate(self, content: str) -> bool:
        """Validate API definition content"""
        try:
            data = self._load_content(content)
            return any(parser.validate(data) for parser in self.parsers.values())
        except Exception:
            return False

    def parse(self, content: str) -> ApiDefinition:
        """Parse API definition from file path or string content"""
        try:
            data = self._load_content(content)
            return self.parse_dict(data)
        except Exception as e:
            raise ValueError(f"Failed to parse API definition: {str(e)}")
            
    def parse_dict(self, data: Dict[str, Any]) -> ApiDefinition:
        """Parse API definition from dictionary"""
        if 'openapi' in data:
            return self.parsers['openapi'].parse(data)
        elif 'graphql' in data:
            return self.parsers['graphql'].parse(data)
        elif 'definitions' in data:
            return self.parsers['soap'].parse(data)
        raise ValueError("Unknown API format")

    def _load_content(self, content: str) -> Dict:
        """Load content from file path or string"""
        if Path(content).exists():
            content = Path(content).read_text()
        return yaml.safe_load(content)

# Add ApiParser to exports
__all__ = [
    "BaseParser",
    "ApiParser",
    "ApiDefinition", 
    "Endpoint",
    "Parameter",
    "AuthRequirement",
    "AuthType",
    "InsertionPoint",
    "OpenAPIParser",
    "GraphQLParser",
    "SOAPParser",
    "SchemaNormalizer"
]
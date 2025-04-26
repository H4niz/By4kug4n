import logging
from typing import Dict, List, Optional
from .models import ApiDefinition, Endpoint, Parameter
from ..config import ParserConfig

class SchemaNormalizer:
    """Schema normalizer for API definitions"""
    
    def __init__(self):
        self.logger = logging.getLogger(__name__)
        self.config = ParserConfig()

    def normalize_schema(self, api_def: ApiDefinition) -> ApiDefinition:
        """Normalize entire API definition"""
        try:
            # Normalize endpoints
            normalized_endpoints = []
            for endpoint in api_def.endpoints:
                normalized_endpoint = self._normalize_endpoint(endpoint)
                normalized_endpoints.append(normalized_endpoint)
                
            api_def.endpoints = normalized_endpoints
            return api_def
            
        except Exception as e:
            self.logger.error(f"Schema normalization failed: {str(e)}")
            raise

    def _normalize_endpoint(self, endpoint: Endpoint) -> Endpoint:
        """Normalize endpoint attributes"""
        # Normalize path
        endpoint.path = self._normalize_path(endpoint.path)
        
        # Normalize parameters
        normalized_params = []
        for param in endpoint.parameters:
            normalized_param = self._normalize_parameter(param)
            normalized_params.append(normalized_param)
        endpoint.parameters = normalized_params
        
        # Normalize method
        endpoint.method = endpoint.method.upper()
        
        return endpoint
        
    def _normalize_path(self, path: str) -> str:
        """Normalize API path"""
        # Remove duplicate slashes
        while '//' in path:
            path = path.replace('//', '/')
            
        # Remove trailing slash
        path = path.rstrip('/')
        
        # Ensure leading slash
        if not path.startswith('/'):
            path = '/' + path
            
        return path
        
    def _normalize_parameter(self, param: Parameter) -> Parameter:
        """Normalize parameter attributes"""
        # Convert name to lowercase
        param.name = param.name.lower()
        
        # Normalize location
        valid_locations = ['query', 'path', 'header', 'body', 'cookie']
        if param.location not in valid_locations:
            param.location = 'query'
            
        # Set default values
        if param.type is None:
            param.type = 'string'
            
        if param.required is None:
            param.required = False
            
        return param
import yaml
from typing import Dict, List, Optional
from .base import BaseParser
from .models import (
    ApiDefinition, Endpoint, Parameter, AuthRequirement,
    AuthType, RequestBody, InsertionPoint
)

class OpenAPIParser(BaseParser):
    def __init__(self):
        self.spec = None

    def validate(self, spec: Dict) -> bool:
        """Validate OpenAPI specification"""
        required_fields = ['openapi', 'paths']
        if not all(field in spec for field in required_fields):
            return False
            
        # Add info section validation
        if 'info' not in spec:
            spec['info'] = {
                'title': 'Unknown API',
                'version': '1.0.0'
            }
            
        return True

    def parse(self, content: Dict) -> ApiDefinition:
        """Parse OpenAPI spec into normalized format"""
        try:
            self.spec = content
            if isinstance(content, str):
                self.spec = yaml.safe_load(content)
                
            if not self.validate(self.spec):
                raise ValueError("Invalid OpenAPI specification")
                
            return ApiDefinition(
                title=self.spec.get('info', {}).get('title', 'Unknown API'),
                version=self.spec.get('info', {}).get('version', '1.0.0'),
                description=self.spec.get('info', {}).get('description'),
                endpoints=self._parse_paths(),
                auth_schemes=self._parse_security_schemes()
            )
        except Exception as e:
            raise ValueError(f"Failed to parse OpenAPI spec: {str(e)}")

    def _parse_paths(self) -> List[Endpoint]:
        """Extract API endpoints from paths section"""
        endpoints = []
        paths = self.spec.get('paths', {})
        
        for path, methods in paths.items():
            for method, operation in methods.items():
                if method.lower() not in ['get', 'post', 'put', 'delete', 'patch']:
                    continue

                parameters = self._parse_parameters(operation)
                request_body = self._parse_request_body(operation)
                auth_reqs = self._parse_operation_security(operation)

                endpoint = Endpoint(
                    path=path,
                    method=method.upper(),
                    name=operation.get('operationId', ''),
                    operation_type="rest",
                    parameters=parameters,
                    request_body=request_body,
                    auth_requirements=auth_reqs,
                    description=operation.get('description', '')
                )
                endpoints.append(endpoint)
                
        return endpoints

    def _parse_parameters(self, operation: Dict) -> List[Parameter]:
        """Parse operation parameters"""
        parameters = []
        for param in operation.get('parameters', []):
            parameter = Parameter(
                name=param['name'],
                location=param['in'],
                required=param.get('required', False),
                type=param.get('schema', {}).get('type', 'string'),
                description=param.get('description')
            )
            
            # Add injection points for vulnerable parameters
            if self._should_add_injection_point(parameter):
                parameter.insertion_points.append(
                    InsertionPoint(
                        param_name=parameter.name,
                        param_type="sql_injection",
                        location=parameter.location
                    )
                )
            parameters.append(parameter)
        return parameters

    def _parse_request_body(self, operation: Dict) -> Optional[RequestBody]:
        """Parse request body schema"""
        request_body = operation.get('requestBody', {})
        if not request_body:
            return None
            
        content = request_body.get('content', {})
        for content_type, content_schema in content.items():
            return RequestBody(
                content_type=content_type,
                schema=content_schema.get('schema', {}),
                required=request_body.get('required', False)
            )
        return None

    def _parse_operation_security(self, operation: Dict) -> List[AuthRequirement]:
        """Parse operation security requirements"""
        requirements = []
        security = operation.get('security', [])
        security_schemes = self.spec.get('components', {}).get('securitySchemes', {})
        
        for sec_req in security:
            for scheme_name, _ in sec_req.items():
                scheme = security_schemes.get(scheme_name)
                if scheme:
                    requirements.append(
                        AuthRequirement(
                            type=self._map_security_type(scheme['type']),
                            location='header',
                            name=scheme.get('name', 'Authorization'),
                            scheme=scheme.get('scheme')
                        )
                    )
        return requirements

    def _parse_security_schemes(self) -> Dict:
        """Extract security schemes from components section"""
        return self.spec.get('components', {}).get('securitySchemes', {})

    def _should_add_injection_point(self, param: Parameter) -> bool:
        """Check if parameter should have injection points"""
        vulnerable_params = ['query', 'search', 'filter', 'where']
        return any(v in param.name.lower() for v in vulnerable_params)

    def _map_security_type(self, scheme_type: str) -> AuthType:
        """Map OpenAPI security scheme to internal auth type"""
        mapping = {
            'http': AuthType.BASIC,
            'apiKey': AuthType.APIKEY,
            'oauth2': AuthType.OAUTH2,
            'openIdConnect': AuthType.OAUTH2
        }
        return mapping.get(scheme_type, AuthType.NONE)
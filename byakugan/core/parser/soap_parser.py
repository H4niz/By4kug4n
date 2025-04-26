import xmltodict
from typing import Dict, List, Optional
from .base import BaseParser
from .models import (
    ApiDefinition, Endpoint, Parameter, AuthRequirement,
    AuthType, RequestBody, InsertionPoint
)

class SOAPParser(BaseParser):
    """Parser for SOAP/WSDL definitions"""
    
    def __init__(self):
        self.wsdl = None
        self.namespaces = {}

    def validate(self, content: str) -> bool:
        """Validate WSDL document"""
        try:
            wsdl = xmltodict.parse(content)
            return isinstance(wsdl, dict) and 'definitions' in wsdl
        except Exception:
            return False

    def parse(self, content: str) -> ApiDefinition:
        """Parse WSDL into normalized format"""
        try:
            self.wsdl = xmltodict.parse(content)
            self._extract_namespaces()
            
            return ApiDefinition(
                title=self._get_service_name(),
                version="1.0",
                endpoints=self._parse_operations(),
                auth_schemes=self._get_auth_schemes()
            )
        except Exception as e:
            raise ValueError(f"Failed to parse WSDL: {str(e)}")

    def _extract_namespaces(self):
        """Extract XML namespaces"""
        definitions = self.wsdl.get("definitions", {})
        xmlns = definitions.get("@xmlns")
        
        # Handle different xmlns formats
        if isinstance(xmlns, dict):
            self.namespaces = {
                k.replace("xmlns:", ""): v
                for k, v in xmlns.items()
            }
        elif isinstance(xmlns, str):
            # Handle single namespace string
            self.namespaces = {"": xmlns}
        else:
            # Default empty namespaces
            self.namespaces = {}

    def _parse_operations(self) -> List[Endpoint]:
        """Extract operations from WSDL"""
        endpoints = []
        port_types = self.wsdl["definitions"].get("portType", [])
        
        if not isinstance(port_types, list):
            port_types = [port_types]
            
        for port_type in port_types:
            operations = port_type.get("operation", [])
            if not isinstance(operations, list):
                operations = [operations]
                
            for operation in operations:
                endpoint = self._create_endpoint(operation)
                if endpoint:
                    endpoints.append(endpoint)
                    
        return endpoints

    def _create_endpoint(self, operation: Dict) -> Optional[Endpoint]:
        """Create endpoint from WSDL operation"""
        operation_name = operation.get("@name")
        if not operation_name:
            return None
            
        # Parse input message
        input_msg = operation.get("input", {})
        input_params = self._get_message_parameters(
            input_msg.get("@message")
        )
        
        # Create request body
        request_body = RequestBody(
            content_type="application/soap+xml",
            schema=self._create_soap_schema(operation_name, input_params),
            required=True
        )
        
        return Endpoint(
            path=f"/soap/{operation_name}",
            method="POST",
            name=operation_name,
            operation_type="soap",
            parameters=input_params,
            auth_requirements=[
                AuthRequirement(
                    type=AuthType.BASIC,
                    location="header",
                    name="Authorization"
                )
            ],
            request_body=request_body
        )

    def _get_message_parameters(self, message_ref: str) -> List[Parameter]:
        """Extract parameters from message definition"""
        parameters = []
        if not message_ref:
            return parameters
            
        # Remove namespace prefix
        message_name = message_ref.split(":")[-1]
        messages = self.wsdl["definitions"].get("message", [])
        
        if not isinstance(messages, list):
            messages = [messages]
            
        for message in messages:
            if message.get("@name") == message_name:
                parts = message.get("part", [])
                if not isinstance(parts, list):
                    parts = [parts]
                    
                for part in parts:
                    # Create parameter first
                    param = Parameter(
                        name=part.get("@name", ""),
                        location="body",
                        required=True,
                        type=part.get("@type", "string").split(":")[-1],
                        description=part.get("documentation")
                    )
                    
                    # Then add injection points if needed
                    if self._should_add_injection_points(param):
                        injection_points = self._create_injection_points(param)
                        param.insertion_points.extend(injection_points)
                        
                    parameters.append(param)
                    
        return parameters

    def _should_add_injection_points(self, param: Parameter) -> bool:
        """Determine if injection points should be added"""
        vuln_names = ['query', 'search', 'filter']
        return any(v in param.name.lower() for v in vuln_names)

    def _create_injection_points(self, param: Parameter) -> List[InsertionPoint]:
        """Create injection points for a parameter"""
        return [
            InsertionPoint(
                param_name=param.name,
                param_type="sql_injection",
                location="body"
            )
        ]

    def _create_soap_schema(self, operation: str, 
                           parameters: List[Parameter]) -> Dict:
        """Create SOAP envelope schema"""
        return {
            "type": "object",
            "properties": {
                "Envelope": {
                    "type": "object",
                    "properties": {
                        "Header": {
                            "type": "object"
                        },
                        "Body": {
                            "type": "object",
                            "properties": {
                                operation: {
                                    "type": "object",
                                    "properties": {
                                        param.name: {
                                            "type": param.type
                                        }
                                        for param in parameters
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }

    def _get_service_name(self) -> str:
        """Extract service name from WSDL"""
        service = self.wsdl["definitions"].get("service", {})
        return service.get("@name", "SOAP Service")

    def _get_auth_schemes(self) -> Dict:
        """Get default authentication schemes for SOAP"""
        return {
            "BasicAuth": {
                "type": AuthType.BASIC,
                "location": "header",
                "name": "Authorization"
            }
        }
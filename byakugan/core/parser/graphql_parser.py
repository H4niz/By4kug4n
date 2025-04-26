from typing import Dict, List, Optional
from graphql import build_schema, parse, GraphQLNonNull
from .base import BaseParser
from .models import (
    ApiDefinition, Endpoint, Parameter, AuthRequirement,
    AuthType, RequestBody, InsertionPoint
)

class GraphQLParser(BaseParser):
    def __init__(self):
        self.schema = None

    def validate(self, content: str) -> bool:
        try:
            parse(content)
            return True
        except Exception:
            return False

    def parse(self, content: str) -> ApiDefinition:
        try:
            self.schema = build_schema(content)
            
            return ApiDefinition(
                title="GraphQL API",
                version="1.0",
                endpoints=self._parse_operations(),
                auth_schemes=self._get_default_auth_schemes()
            )
        except Exception as e:
            raise ValueError(f"Failed to parse GraphQL schema: {str(e)}")

    def _parse_operations(self) -> List[Endpoint]:
        endpoints = []
        
        # Parse Query operations
        query_type = self.schema.get_type('Query')
        if query_type and hasattr(query_type, 'fields'):
            for field_name, field in query_type.fields.items():
                endpoints.append(self._create_endpoint(
                    path="/graphql",
                    method="POST",
                    name=field_name,
                    operation_type="query", 
                    field=field
                ))

        # Parse Mutation operations
        mutation_type = self.schema.get_type('Mutation')
        if mutation_type and hasattr(mutation_type, 'fields'):
            for field_name, field in mutation_type.fields.items():
                endpoints.append(self._create_endpoint(
                    path="/graphql",
                    method="POST",
                    name=field_name,
                    operation_type="mutation",
                    field=field
                ))

        return endpoints

    def _create_endpoint(
        self, path: str, method: str,
        name: str, operation_type: str,
        field: any
    ) -> Endpoint:
        parameters = []
        
        # Extract field arguments as parameters
        if hasattr(field, 'args'):
            for arg_name, arg in field.args.items():
                # Check if type is non-null
                required = isinstance(arg.type, GraphQLNonNull)
                
                param = Parameter(
                    name=arg_name,
                    location="body",
                    required=required,  # Set based on GraphQLNonNull
                    type=str(arg.type),
                    description=getattr(arg, 'description', None)
                )
                
                # Add injection points for potentially vulnerable parameters
                if arg_name.lower() in ['query', 'search', 'filter']:
                    param.insertion_points.append(
                        InsertionPoint(
                            param_name=arg_name,
                            param_type="sql_injection",
                            location="body"
                        )
                    )
                
                parameters.append(param)

        # Create request body schema
        request_body = RequestBody(
            content_type="application/json",
            required=True,
            schema={
                "type": "object",
                "properties": {
                    "query": {"type": "string"},
                    "variables": {"type": "object"},
                    "operationName": {"type": "string"}
                },
                "required": ["query"]
            }
        )

        # Add JWT auth requirement
        auth_requirements = [
            AuthRequirement(
                type=AuthType.JWT,
                location="header",
                name="Authorization"
            )
        ]

        return Endpoint(
            path=path,
            method=method,
            name=name,
            operation_type=operation_type,
            parameters=parameters,
            request_body=request_body,
            auth_requirements=auth_requirements,
            description=getattr(field, 'description', None)
        )

    def _get_default_auth_schemes(self) -> Dict:
        return {
            "jwt": {
                "type": "http",
                "scheme": "bearer",
                "bearerFormat": "JWT"
            }
        }
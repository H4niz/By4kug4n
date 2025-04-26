from dataclasses import dataclass, asdict, field
from typing import List, Dict, Optional
from enum import Enum

class AuthType(Enum):
    NONE = "none"
    BASIC = "basic"
    APIKEY = "apikey" 
    OAUTH2 = "oauth2"
    JWT = "jwt"

@dataclass
class InsertionPoint:
    param_name: str
    param_type: str  # sql_injection, xss, etc
    location: str    # query, body, header
    payloads: List[str] = field(default_factory=list)

@dataclass
class Parameter:
    name: str
    location: str  # query, path, header
    required: bool = False
    type: str = "string"
    description: Optional[str] = None
    insertion_points: List[InsertionPoint] = field(default_factory=list)

@dataclass
class RequestBody:
    content_type: str
    schema: Dict
    required: bool = False

@dataclass 
class AuthRequirement:
    type: AuthType
    location: str = "header"
    name: str = "Authorization"
    scheme: Optional[str] = None

@dataclass
class Endpoint:
    path: str
    method: str
    name: str = ""
    operation_type: str = "rest"
    parameters: List[Dict] = field(default_factory=list)
    request_body: Optional[Dict] = None
    auth_requirements: List[Dict] = field(default_factory=list)
    description: str = ""

    def to_dict(self) -> Dict:
        """Convert endpoint to dictionary"""
        return asdict(self)

@dataclass
class ApiDefinition:
    title: str
    version: str
    description: Optional[str] = None
    endpoints: List[Endpoint] = field(default_factory=list)
    auth_schemes: Dict = field(default_factory=dict)
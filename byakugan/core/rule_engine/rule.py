from typing import Dict, List, Optional
from dataclasses import dataclass, field

@dataclass
class Rule:
    """Security rule definition"""
    id: str
    name: str
    description: str
    severity: str
    category: str
    method: str = "ANY"  # Default to ANY
    required_parameters: List[Dict] = field(default_factory=list)
    detection: Dict = field(default_factory=dict)
    payloads: Dict = field(default_factory=dict)
    patterns: List[str] = field(default_factory=list)

    @classmethod
    def from_dict(cls, data: Dict) -> 'Rule':
        """Create Rule from dictionary"""
        return cls(
            id=data["id"],
            name=data["name"],
            description=data["description"],
            severity=data["severity"],
            category=data["category"],
            method=data.get("method", "ANY"),
            required_parameters=data.get("required_parameters", []),
            detection=data.get("detection", {}),
            payloads=data.get("payloads", {}),
            patterns=data.get("patterns", [])
        )
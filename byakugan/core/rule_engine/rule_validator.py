import logging
from typing import List, Dict
from .rule import Rule

class RuleValidator:
    """Validates security rules"""
    
    def __init__(self):
        self.logger = logging.getLogger(__name__)

    def validate_rule(self, rule: Rule) -> bool:
        """Validate rule definition"""
        try:
            # Check required fields
            if not all([rule.id, rule.name, rule.severity, rule.category]):
                return False

            # Validate severity
            if rule.severity not in ["LOW", "MEDIUM", "HIGH", "CRITICAL"]:
                return False

            # Validate detection config
            if not self._validate_detection(rule.detection):
                return False

            # Validate payloads
            if not rule.payloads:
                return False

            return True

        except Exception as e:
            self.logger.error(f"Rule validation error: {str(e)}")
            return False

    def _validate_detection(self, detection: Dict) -> bool:
        """Validate detection configuration"""
        try:
            # Check locations
            if "locations" not in detection:
                return False
                
            # Check strategies
            if "strategies" not in detection:
                return False

            # Must have at least one active strategy
            strategies = detection["strategies"]
            if not any(s.get("active", False) for s in strategies):
                return False

            return True

        except Exception:
            return False
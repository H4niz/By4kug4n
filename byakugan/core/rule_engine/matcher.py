from typing import List, Dict
from .rule import Rule
import logging

class RuleMatcher:
    """Matches rules with endpoints"""
    
    def __init__(self, rules: List[Rule] = None):
        self.rules = rules or []
        self.logger = logging.getLogger(__name__)

    def match_endpoint(self, endpoint: Dict) -> List[Rule]:
        """Find rules matching an endpoint"""
        matches = []
        
        for rule in self.rules:
            try:
                if self._matches_method(rule, endpoint) and \
                   self._matches_parameters(rule, endpoint) and \
                   self._matches_location(rule, endpoint):
                    matches.append(rule)
            except Exception as e:
                self.logger.error(f"Error matching rule {rule.id}: {str(e)}")
                continue
                
        return matches

    def _matches_method(self, rule: Rule, endpoint: Dict) -> bool:
        """Check if rule matches HTTP method"""
        if not hasattr(rule, 'method') or rule.method == "ANY":
            return True
        return rule.method.upper() == endpoint.get("method", "").upper()

    def _matches_parameters(self, rule: Rule, endpoint: Dict) -> bool:
        """Check if endpoint has required parameters"""
        required_params = rule.required_parameters
        endpoint_params = {p["name"]: p for p in endpoint.get("parameters", [])}
        
        for required in required_params:
            param_name = required["name"]
            param_location = required["location"]
            
            if param_name not in endpoint_params:
                return False
                
            if endpoint_params[param_name]["location"] != param_location:
                return False
                
        return True

    def _matches_location(self, rule: Rule, endpoint: Dict) -> bool:
        """Check if rule matches location"""
        detection = rule.detection
        if not detection:
            return True
            
        locations = detection.get("locations", {})
        
        # Check path parameters
        path_params = locations.get("path_parameters", [])
        if path_params:
            for param in endpoint.get("parameters", []):
                if param.get("location") == "path" and param.get("name") in path_params:
                    return True

        # Check query parameters            
        query_params = locations.get("query_parameters", [])
        if query_params:
            for param in endpoint.get("parameters", []):
                if param.get("location") == "query" and param.get("name") in query_params:
                    return True
                    
        # Check headers
        headers = locations.get("headers", [])
        if headers and endpoint.get("headers", {}).keys() & set(headers):
            return True
                    
        return False
import logging
from typing import List, Dict, Optional
from .rule import Rule
from .matcher import RuleMatcher
from .rule_loader import RuleLoader
from .rule_validator import RuleValidator
from ..config import RuleEngineConfig

class RuleEngine:
    """Core rule engine for loading and managing security rules"""
    
    def __init__(self, rules_dir: str):
        self.rules = {}
        self.matcher = RuleMatcher()
        self.loader = RuleLoader()
        self.validator = RuleValidator()
        self.load_rules(rules_dir)

    def load_rules(self, rules_dir: str):
        """Load rules from YAML files"""
        for rule_data in self.loader.load_rules_from_directory(rules_dir):
            rule = Rule.from_dict(rule_data)
            if self.validator.validate_rule(rule):
                self.rules[rule.id] = rule

    def get_rule_by_id(self, rule_id: str) -> Optional[Rule]:
        """Get rule by ID"""
        return self.rules.get(rule_id)

    def get_matching_rules(self, endpoint: Dict) -> List[Rule]:
        """Find rules matching endpoint"""
        return self.matcher.match_endpoint(endpoint)
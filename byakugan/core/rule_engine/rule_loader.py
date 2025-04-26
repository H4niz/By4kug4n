import os
import yaml
import logging
from typing import Dict, List, Optional
from pathlib import Path
from .rule import Rule

class RuleLoader:
    """Handles loading and validating security rules"""
    
    def __init__(self):
        self.logger = logging.getLogger(__name__)

    def load_rules_from_file(self, file_path: str) -> List[Dict]:
        """Load rules from a YAML file"""
        try:
            with open(file_path) as f:
                data = yaml.safe_load(f)
                
            if not isinstance(data, dict) or "rules" not in data:
                raise ValueError(f"Invalid rule format in {file_path}")
                
            return data["rules"]
            
        except Exception as e:
            self.logger.error(f"Failed to load rules from {file_path}: {str(e)}")
            return []

    def load_rules_from_directory(self, rules_dir: str) -> List[Dict]:
        """Load all rules from directory"""
        rules = []
        try:
            for yaml_file in Path(rules_dir).glob("*.yaml"):
                file_rules = self.load_rules_from_file(str(yaml_file))
                rules.extend(file_rules)
                
        except Exception as e:
            self.logger.error(f"Failed to load rules from directory: {str(e)}")
            
        return rules

    def validate_rule_dependencies(self, rules: List[Rule]) -> bool:
        """Validate rule dependencies"""
        rule_ids = {rule.id for rule in rules}
        
        for rule in rules:
            # Check if all required rules exist
            for dep in rule.dependencies:
                if dep not in rule_ids:
                    self.logger.error(f"Missing dependency {dep} for rule {rule.id}")
                    return False
                    
        return True
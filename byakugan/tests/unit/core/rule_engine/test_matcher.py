import pytest
from byakugan.core.rule_engine.rule import Rule
from byakugan.core.rule_engine.matcher import RuleMatcher

@pytest.fixture
def sample_rules():
    return [
        Rule.from_dict({
            "id": "JWT-001",
            "name": "JWT None Algorithm",
            "severity": "HIGH",
            "category": "Authentication",
            "payloads": ["eyJ0eXAiOiJKV1QiLCJhbGciOiJub25lIn0"],
            "validation": {
                "response_codes": [200],
                "patterns": ["authenticated"]
            }
        })
    ]

@pytest.fixture
def sample_endpoint():
    return {
        "path": "/auth/verify",
        "method": "POST",
        "security": [{"BearerAuth": []}],
        "requestBody": {
            "content": {
                "application/json": {
                    "schema": {
                        "type": "object",
                        "properties": {
                            "token": {"type": "string"}
                        }
                    }
                }
            }
        }
    }

@pytest.fixture
def sql_injection_rule():
    return Rule(
        id="SQL-001",
        name="SQL Injection",
        description="Detects SQL injection vulnerabilities",
        severity="HIGH",
        category="Injection",
        method="POST",
        required_parameters=[
            {
                "name": "query",
                "location": "body",
                "type": "string"
            }
        ],
        payloads=["' OR '1'='1"],
        validation={
            "response_patterns": ["error in your SQL syntax"]
        }
    )

@pytest.fixture
def auth_bypass_rule():
    return Rule(
        id="AUTH-001", 
        name="Authentication Bypass",
        description="Detects authentication bypass vulnerabilities",
        severity="CRITICAL",
        category="Authentication",
        method=None,  # Applies to any method
        required_parameters=[],
        payloads=["none"],
        validation={},
        auth_requirements=["Bearer"]
    )

def test_rule_matching(sample_rules, sample_endpoint):
    matcher = RuleMatcher(sample_rules)
    matched = matcher.match_endpoint(sample_endpoint)
    
    assert len(matched) == 1
    assert matched[0].id == "JWT-001"

class TestRuleMatcher:
    def test_match_endpoint_with_method(self, sql_injection_rule):
        """Test matching rule with specific HTTP method"""
        matcher = RuleMatcher([sql_injection_rule])
        
        # Should match
        endpoint = {
            "method": "POST",
            "parameters": [
                {
                    "name": "query",
                    "location": "body",
                    "type": "string"
                }
            ]
        }
        matches = matcher.match_endpoint(endpoint)
        assert len(matches) == 1
        assert matches[0].id == "SQL-001"
        
        # Should not match - wrong method
        endpoint["method"] = "GET"
        matches = matcher.match_endpoint(endpoint)
        assert len(matches) == 0

    def test_match_endpoint_parameters(self, sql_injection_rule):
        """Test matching based on required parameters"""
        matcher = RuleMatcher([sql_injection_rule])
        
        # Should match - has required parameter
        endpoint = {
            "method": "POST",
            "parameters": [
                {
                    "name": "query",
                    "location": "body",
                    "type": "string"
                }
            ]
        }
        matches = matcher.match_endpoint(endpoint)
        assert len(matches) == 1
        
        # Should not match - missing required parameter
        endpoint["parameters"] = []
        matches = matcher.match_endpoint(endpoint)
        assert len(matches) == 0

    def test_match_auth_requirements(self, auth_bypass_rule):
        """Test matching based on auth requirements"""
        matcher = RuleMatcher([auth_bypass_rule])
        
        # Should match - has required auth
        endpoint = {
            "method": "GET",
            "security": [{"type": "Bearer"}]
        }
        matches = matcher.match_endpoint(endpoint)
        assert len(matches) == 1
        
        # Should not match - missing required auth
        endpoint["security"] = []
        matches = matcher.match_endpoint(endpoint)
        assert len(matches) == 0

    def test_body_parameter_extraction(self, sql_injection_rule):
        """Test extraction of parameters from request body"""
        matcher = RuleMatcher([sql_injection_rule])
        
        endpoint = {
            "method": "POST",
            "requestBody": {
                "content": {
                    "application/json": {
                        "schema": {
                            "properties": {
                                "query": {
                                    "type": "string"
                                }
                            },
                            "required": ["query"]
                        }
                    }
                }
            }
        }
        
        matches = matcher.match_endpoint(endpoint)
        assert len(matches) == 1
        assert matches[0].id == "SQL-001"
import os
import pytest
from pathlib import Path
from core.rule_engine import RuleEngine, Rule, RuleValidator
from core.config import RuleEngineConfig

@pytest.fixture
def test_rules_dir():
    """Create test rules directory with sample rules"""
    current_dir = Path(__file__).parent
    rules_dir = current_dir / "testdata" / "test_rules"
    rules_dir.mkdir(parents=True, exist_ok=True)
    
    # Create BOLA rule
    bola_rule = rules_dir / "bola.yaml"
    bola_rule.write_text("""
rules:
  - id: TEST-BOLA-001
    name: "Test BOLA Rule"
    description: "Test rule for BOLA detection" 
    severity: "HIGH"
    category: "authorization"
    method: "GET"  # Add method
    required_parameters:
      - name: "id"
        location: "path"
    detection:
      locations:
        path_parameters:
          - id
          - user_id
        query_parameters:
          - id
      strategies:
        - name: "ID Enumeration"
          active: true
    payloads:
      id_manipulation:
        - template: "{{base_id + 1}}"
        - template: "{{base_id - 1}}"
    patterns:
      - "error"
      - "invalid"
""")

    # Create JWT rule
    jwt_rule = rules_dir / "jwt.yaml"
    jwt_rule.write_text("""
rules:
  - id: TEST-JWT-001
    name: "Test JWT Rule" 
    description: "Test rule for JWT vulnerabilities"
    severity: "CRITICAL"
    category: "authentication"
    required_parameters:
      - name: "token"
        location: "header"
    detection:
      locations:
        headers:
          - Authorization
        query_parameters:
          - token
      strategies:
        - name: "None Algorithm"
          active: true
    payloads:
      none_algorithm:
        - template: "eyJ0eXAiOiJKV1QiLCJhbGciOiJub25lIn0.{{payload}}."
    patterns:
      - "authenticated"
      - "authorized"
""")

    return str(rules_dir)

@pytest.fixture
def rule_engine(test_rules_dir):
    """Create RuleEngine with test configuration"""
    config = RuleEngineConfig(
        rules_dir=test_rules_dir,
        custom_rules_dir=None,
        rule_timeout=30,
        enabled_categories=["authentication", "authorization"],
        severity_threshold="LOW",
        max_workers=2
    )
    return RuleEngine(config)

def test_rule_loading(rule_engine, test_rules_dir):
    """Test rule loading"""
    # Load rules
    rule_engine.load_rules(test_rules_dir)
    
    # Verify rules were loaded
    assert len(rule_engine.rules) == 2
    
    # Verify BOLA rule
    bola_rule = rule_engine.rules.get("TEST-BOLA-001")
    assert bola_rule is not None
    assert bola_rule.category == "authorization"
    assert bola_rule.severity == "HIGH"
    assert "id_manipulation" in bola_rule.payloads
    assert "error" in bola_rule.patterns
    assert "invalid" in bola_rule.patterns
    
    # Verify JWT rule
    jwt_rule = rule_engine.rules.get("TEST-JWT-001")
    assert jwt_rule is not None
    assert jwt_rule.category == "authentication" 
    assert jwt_rule.severity == "CRITICAL"
    assert "none_algorithm" in jwt_rule.payloads
    assert "authenticated" in jwt_rule.patterns
    assert "authorized" in jwt_rule.patterns

def test_rule_validation():
    """Test rule validation"""
    validator = RuleValidator()

    valid_rule = Rule(
        id="TEST-001",
        name="Test Rule",
        description="Test rule",
        severity="HIGH",
        category="authentication",
        detection={
            "locations": {
                "headers": ["Authorization"]
            },
            "strategies": [
                {"name": "test", "active": True}
            ]
        },
        payloads={
            "test": ["payload"]
        },
        patterns=["error"]
    )
    assert validator.validate_rule(valid_rule)

def test_rule_matching(rule_engine, test_rules_dir):
    """Test rule matching"""
    rule_engine.load_rules(test_rules_dir)

    # Test BOLA rule matching
    endpoint = {
        "path": "/users/{id}",
        "method": "GET",
        "parameters": [
            {
                "name": "id",
                "location": "path",
                "value": "123"
            }
        ]
    }
    matches = rule_engine.get_matching_rules(endpoint)
    assert len(matches) == 1
    assert matches[0].id == "TEST-BOLA-001"

    # Test JWT rule matching
    endpoint = {
        "path": "/auth",
        "method": "POST",
        "headers": {
            "Authorization": "Bearer token" 
        },
        "parameters": [
            {
                "name": "token",
                "location": "header",
                "value": "xyz"
            }
        ]
    }
    matches = rule_engine.get_matching_rules(endpoint)
    assert len(matches) == 1
    assert matches[0].id == "TEST-JWT-001"

@pytest.mark.asyncio 
async def test_rule_execution(rule_engine, test_rules_dir):
    """Test rule execution"""
    rule_engine.load_rules(test_rules_dir)
    
    # Get BOLA rule
    rule = rule_engine.rules["TEST-BOLA-001"]
    
    # Test endpoint
    endpoint = {
        "path": "/users/123",
        "method": "GET",
        "parameters": [
            {
                "name": "id",
                "location": "path",
                "value": "123"
            }
        ]
    }
    
    # Execute rule
    results = await rule_engine.execute_rule(rule, endpoint)
    
    # Verify results
    assert len(results) > 0
    assert results[0]["rule_id"] == "TEST-BOLA-001"
    assert "payloads" in results[0]

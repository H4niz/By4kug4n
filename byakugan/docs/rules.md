# Rule Development | Phát Triển Rule

## Rule Structure | Cấu Trúc Rule

```yaml
id: BYKG-001
name: SQL Injection Detection
severity: HIGH
category: injection

patterns:
  - type: regex
    value: "(?i)(UNION|SELECT|INSERT|UPDATE|DELETE).*"
  
validation:
  response_patterns:
    - "SQL syntax.*MySQL"
    - "Warning.*SQL"
  
evidence_requirements:
  - response_code: 500
  - response_body_pattern: "database error"
```

## Creating Rules | Tạo Rule

### 1. Basic Structure | Cấu Trúc Cơ Bản
- Rule metadata | Metadata của rule
- Detection patterns | Mẫu phát hiện
- Validation logic | Logic xác thực
- Evidence collection | Thu thập bằng chứng

### 2. Pattern Types | Loại Mẫu
- Regex patterns | Mẫu regex
- Keywords | Từ khóa
- Semantic rules | Rule ngữ nghĩa
- Custom functions | Hàm tùy chỉnh

### 3. Validation | Xác Thực
- Response codes | Mã phản hồi
- Response patterns | Mẫu phản hồi
- Headers | Headers
- Timing | Thời gian

## Rule Testing | Kiểm Thử Rule

### 1. Unit Tests | Unit Test
```python
def test_sql_injection_rule():
    rule = Rule.from_yaml("rules/sql-injection.yaml")
    result = rule.test_payload("' OR '1'='1")
    assert result.matched
    assert result.severity == "HIGH"
```

### 2. Integration Tests | Test Tích Hợp
```python
async def test_rule_execution():
    scanner = Scanner()
    result = await scanner.execute_rule("BYKG-001", target_url)
    assert result.findings
    assert result.evidence
```
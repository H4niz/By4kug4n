# Byakugan Configuration Guide | Hướng Dẫn Cấu Hình Byakugan

## Core Configuration | Cấu Hình Core

```yaml
core:
  debug: false                # Enable debug mode | Bật chế độ debug
  log_level: "info"          # Logging level (debug, info, warning, error) | Mức độ log
  log_file: "byakugan.log"   # Log file path | Đường dẫn file log
  temp_dir: "tmp"            # Temporary directory | Thư mục tạm
```

### Parameters | Tham Số
| Parameter | Type | Default | Description | Mô tả |
|-----------|------|---------|-------------|--------|
| debug | boolean | false | Enable debug mode | Bật chế độ debug |
| log_level | string | "info" | Log level | Mức độ log |
| log_file | string | "byakugan.log" | Log file path | Đường dẫn file log |
| temp_dir | string | "tmp" | Temporary directory | Thư mục tạm |

## Scanner Configuration | Cấu Hình Scanner

```yaml
scanner:
  node:
    id: "scanner-01"         # Scanner node ID | ID node scanner
    name: "Primary Scanner"  # Scanner name | Tên scanner
    region: "default"        # Deployment region | Vùng triển khai
    tags:                    # Node tags | Thẻ node
      - "production"

  performance:
    cpu_limit: 80           # CPU usage limit (%) | Giới hạn CPU
    memory_limit: 2048      # Memory limit (MB) | Giới hạn bộ nhớ
    bandwidth_limit: 100    # Bandwidth limit (Mbps) | Giới hạn băng thông
    connections_limit: 500  # Max connections | Giới hạn kết nối

  http_client:
    user_agent: "Byakugan Scanner v1.0.0"  # User agent string | Chuỗi user agent
    follow_redirects: true   # Follow HTTP redirects | Theo dõi chuyển hướng HTTP
    max_redirects: 10        # Maximum redirects | Số lần chuyển hướng tối đa
    timeouts:
      connect: 5            # Connection timeout (s) | Thời gian timeout kết nối
      read: 30             # Read timeout (s) | Thời gian timeout đọc
      write: 5             # Write timeout (s) | Thời gian timeout ghi
```

## Communication Settings | Cấu Hình Giao Tiếp

```yaml
comms:
  grpc:
    host: "localhost"        # gRPC host | Host gRPC
    port: 50051             # gRPC port | Cổng gRPC
    max_workers: 10         # Max worker threads | Số luồng worker tối đa
    timeout: 30             # Request timeout (s) | Thời gian timeout yêu cầu
    max_message_size: 10485760  # Max message size (bytes) | Kích thước message tối đa

  connection:
    timeout: 30             # Connection timeout (s) | Thời gian timeout kết nối
    retries: 3              # Max retries | Số lần thử lại tối đa
    retry_delay: 5          # Retry delay (s) | Thời gian giữa các lần thử lại

  security:
    tls_enabled: false      # Enable TLS | Bật TLS
    tls_cert_path: null     # TLS certificate path | Đường dẫn chứng chỉ TLS
    tls_key_path: null      # TLS key path | Đường dẫn khóa TLS
```

## Rule Engine Configuration | Cấu Hình Rule Engine

```yaml
rule_engine:
  rules_dir: "rules"        # Rules directory | Thư mục rules
  custom_rules_dir: "custom_rules"  # Custom rules directory | Thư mục rules tùy chỉnh
  rule_timeout: 30          # Rule execution timeout (s) | Thời gian timeout thực thi rule
  rules:
    enabled_categories:     # Enabled rule categories | Danh mục rules được bật
      - "authentication"
      - "authorization"
      - "injection"
      - "business_logic"
    severity_threshold: "low"  # Minimum severity | Ngưỡng mức độ nghiêm trọng
    max_payload_size: 8192    # Max payload size (bytes) | Kích thước payload tối đa
```

## Authentication Configuration | Cấu Hình Xác Thực

```yaml
auth:
  methods:                  # Authentication methods | Phương thức xác thực
    - type: "oauth2"
      enabled: true
      config:
        client_id: ""
        client_secret: ""
        token_url: ""
        scopes: []
    - type: "bearer"
      enabled: true
    - type: "basic"
      enabled: true
    - type: "apikey"
      enabled: true
      location: "header"    # header, query
      name: "X-API-Key"

  session:
    timeout: 3600          # Session timeout (s) | Thời gian timeout phiên
    renew_before: 300      # Renew threshold (s) | Ngưỡng làm mới phiên
    max_retries: 3         # Max retry attempts | Số lần thử lại tối đa
```

## Data Types | Kiểu Dữ Liệu

| Type | Description | Mô tả |
|------|-------------|--------|
| boolean | true/false value | Giá trị true/false |
| string | Text value | Giá trị văn bản |
| integer | Whole number | Số nguyên |
| float | Decimal number | Số thập phân |
| array | List of values | Danh sách giá trị |
| object | Key-value pairs | Cặp khóa-giá trị |

## Environment Variables | Biến Môi Trường

Configuration values can be overridden using environment variables:
Các giá trị cấu hình có thể được ghi đè bằng biến môi trường:

```bash
BYAKUGAN_CORE_DEBUG=true
BYAKUGAN_SCANNER_CPU_LIMIT=50
BYAKUGAN_COMMS_GRPC_HOST=scanner.local
```
# Byakugan - API Security Scanner

<p align="center">
  <img src="media/byakugan-logo.svg" alt="Byakugan Logo" width="200"/>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Python-3.9+-blue.svg" alt="Python Version"/>
  <img src="https://img.shields.io/badge/Go-1.18+-00ADD8.svg" alt="Go Version"/>
  <img src="https://img.shields.io/badge/Version-1.0.0-green.svg" alt="Version"/>
  <img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License"/>
</p>

Byakugan là một công cụ quét lỗ hổng API hiện đại được phát triển với Python và Go, hỗ trợ đa dạng các định dạng API và chiến lược phát hiện lỗ hổng.

## Tính năng chính

- **Hỗ trợ đa định dạng API**
  - OpenAPI/Swagger (2.0, 3.0, 3.1)
  - Postman Collections
  - GraphQL
  - SOAP/WSDL
  - gRPC

- **Engine quét mạnh mẽ**
  - Quét đồng thời nhiều endpoints
  - Tự động phát hiện injection points
  - Rate limiting và retry thông minh
  - Worker pool tối ưu
  - Proxy rotation tích hợp

- **Rule Engine linh hoạt**
  - Rule-based scanning
  - Custom rule development
  - Rule profiles (OWASP API Top 10, PCI-DSS)
  - Plugin architecture

- **Xác thực đa dạng**
  - OAuth 2.0 (tất cả flows)
  - JWT
  - API Keys
  - Basic Auth
  - Custom auth methods

- **Báo cáo chi tiết**
  - HTML/PDF reports
  - JSON/CSV exports
  - SARIF format
  - Evidence collection
  - Remediation guidance

## Cài đặt

### Yêu cầu

- Python 3.9+
- Go 1.18+
- Docker (tùy chọn)

### Cài đặt từ source

...

## Contributing

Xem [CONTRIBUTING.md](CONTRIBUTING.md) để biết thêm chi tiết về cách đóng góp.

## License

Project được phát hành dưới [MIT License](LICENSE).
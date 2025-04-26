# Development Guide | Hướng Dẫn Phát Triển

## Setup Development Environment | Cài Đặt Môi Trường Phát Triển

### Requirements | Yêu Cầu
- Python 3.9+
- Go 1.18+
- Git
- Docker (optional)

### Installation | Cài Đặt
```bash
# Clone repository | Sao chép repo
git clone https://github.com/h4niz/byakugan.git
cd byakugan

# Create virtual environment | Tạo môi trường ảo
python -m venv .venv
.\.venv\Scripts\activate

# Install dependencies | Cài đặt dependencies
pip install -r requirements.txt
go mod download
```

## Project Structure | Cấu Trúc Dự Án

```plaintext
byakugan/
├── cmd/cli/          # CLI interface | Giao diện CLI
├── core/            # Core service | Dịch vụ Core
├── scanner/         # Scanner service | Dịch vụ Scanner
├── proto/           # Protocol definitions | Định nghĩa giao thức
└── docs/           # Documentation | Tài liệu
```

## Development Workflow | Quy Trình Phát Triển

### 1. Code Style | Phong Cách Code
- Python: PEP 8
- Go: gofmt
- Use type hints | Sử dụng type hints
- Write docstrings | Viết docstrings

### 2. Testing | Kiểm Thử
```bash
# Run Python tests | Chạy test Python
pytest

# Run Go tests | Chạy test Go
go test ./...
```

### 3. Building | Xây Dựng
```bash
# Build CLI | Xây dựng CLI
pip install -e .

# Build Scanner | Xây dựng Scanner
cd scanner
go build -o scanner.exe
```
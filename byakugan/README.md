# Byakugan Scanner Service

## Overview
Security vulnerability scanner with distributed architecture using gRPC for service communication.

## Project Structure
```
byakugan/
├── scanner/          # Scanner service (Go)
│   ├── engine/      # Core scanning engine
│   ├── comms/       # Communication layer
│   ├── proto/       # Protocol definitions  
│   ├── config/      # Configuration
│   └── cmd/         # Command line tools
└── core/            # Core service (Python)
    ├── api/         # API handlers
    ├── models/      # Data models
    └── services/    # Business logic
```

## Features

### Core Scanning ✅
- [x] Task-based scanning architecture
- [x] Concurrent scan execution with worker pools
- [x] Evidence collection and analysis
- [x] Multiple vulnerability detection patterns
- [x] HTTP request/response handling
- [x] Configurable scan parameters
- [x] Basic authentication support

### Communication ✅ 
- [x] gRPC service definitions
- [x] Result streaming
- [x] Task status monitoring
- [x] Health checks
- [x] Server-client protocols

### Security Features ✅
- [x] SQL Injection detection
- [x] NoSQL Injection detection
- [x] Path Traversal scanning
- [x] Input validation
- [x] Basic auth

## In Development 🔄

### Engine Enhancements
- [ ] Advanced pattern matching
- [ ] Custom rule engine
- [ ] Result persistence
- [ ] Metrics collection
- [ ] Performance optimization

### Infrastructure
- [ ] Load balancing
- [ ] Service discovery
- [ ] Distributed coordination
- [ ] High availability
- [ ] Container support

## Planned Features 📋

### Core Features
- [ ] Real-time progress updates
- [ ] Report generation
- [ ] Dashboard interface
- [ ] Plugin system
- [ ] API documentation

### Security
- [ ] TLS/mTLS support
- [ ] Role-based access control
- [ ] Rate limiting
- [ ] Advanced authentication

### Integration
- [ ] CI/CD pipeline hooks
- [ ] Webhook notifications
- [ ] External logging
- [ ] Monitoring systems
- [ ] Third-party APIs

## Quick Start

### Prerequisites
- Go 1.22+
- Python 3.9+
- Protocol Buffers
- gRPC tools

### Installation
```bash
# Clone repository
git clone https://github.com/haniz/byakugan.git
cd byakugan

# Install Go dependencies
cd scanner
go mod download

# Install Python dependencies
cd ../core
pip install -r requirements.txt
```

### Running Services
```bash
# Start Scanner Service
cd scanner
go run main.go -config config/scanner.yaml

# Start Core Service
cd ../core
python -m core.main
```

### Running Tests
```bash
# Go tests
cd scanner
go test ./...

# Python tests 
cd core
pytest
```

## Development Status
- Core scanning engine: 80% complete
- gRPC communication: 70% complete
- Security features: 60% complete
- Infrastructure: 40% complete

## Next Steps
1. Complete core engine implementation
2. Add persistence layer
3. Implement metrics collection
4. Add security enhancements
5. Build deployment infrastructure

## Contributing
Please see [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines.

## License
MIT License - see [LICENSE](LICENSE) for details
# API Reference | Tài Liệu API

## gRPC Services | Dịch Vụ gRPC

### ScannerService | Dịch Vụ Scanner

```protobuf
service ScannerService {
    rpc ExecuteScan(ScanTask) returns (ScanResult) {}
    rpc StreamResults(stream ScanResult) returns (stream StreamAck) {}
    rpc GetTaskStatus(TaskStatusRequest) returns (TaskStatus) {}
    rpc Heartbeat(HeartbeatRequest) returns (HeartbeatResponse) {}
}
```

### Message Types | Kiểu Message

```protobuf
message ScanTask {
    string id = 1;
    Target target = 2;
    RuleContext rule_context = 3;
    TaskConfig config = 4;
}

message ScanResult {
    string task_id = 1;
    bool success = 2;
    repeated Finding findings = 3;
}
```

## REST APIs | API REST

### Scan Management | Quản Lý Quét

#### Start Scan | Bắt Đầu Quét
```http
POST /api/v1/scans
Content-Type: application/json

{
    "target": "http://example.com",
    "rules": ["sql-injection", "xss"],
    "config": {
        "concurrent_limit": 10,
        "timeout": 30
    }
}
```

#### Get Status | Lấy Trạng Thái
```http
GET /api/v1/scans/{scan_id}/status
```

## Client SDKs | SDK Client

### Python Client | Client Python
```python
from byakugan.client import ByakuganClient

client = ByakuganClient()
result = await client.execute_scan(
    target="http://example.com",
    rules=["sql-injection"]
)
```

### Go Client | Client Go
```go
client := byakugan.NewClient(config)
result, err := client.ExecuteScan(ctx, &ScanTask{
    Target: "http://example.com",
    Rules: []string{"sql-injection"},
})
```
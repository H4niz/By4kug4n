package detector

import (
	"context"
	"fmt"
	"time"

	"github.com/haniz/byakugan/scanner/engine/logger"
	pb "github.com/haniz/byakugan/scanner/proto"
)

type MockDetector struct{}

func NewMockDetector() *MockDetector {
	return &MockDetector{}
}

func (m *MockDetector) DetectVulnerability(ctx context.Context, dc *DetectionContext) (*pb.Evidence, error) {
	return &pb.Evidence{
		Data: map[string]string{
			"request":         "POST /users/login HTTP/1.1\nContent-Type: application/json\n\n{\"username\":\"admin' --\"}",
			"response":        "HTTP/1.1 500 Internal Server Error\n\n{\"error\":\"mysql_error\"}",
			"injection_point": dc.InsertionPoint.Location,
			"payload":         dc.InsertionPoint.Payloads[0],
		},
	}, nil
}

type mockLogger struct {
	infoMsgs  []string
	errorMsgs []string
	debugMsgs []string
}

// Implement missing WithField method
func (m *mockLogger) WithField(key string, value interface{}) logger.Logger {
	return m.WithFields(map[string]interface{}{key: value})
}

// Implement missing WithFields method
func (m *mockLogger) WithFields(fields map[string]interface{}) logger.Logger {
	// Return self since it's a mock
	return m
}

func (m *mockLogger) Info(msg string, fields map[string]interface{}) {
	m.infoMsgs = append(m.infoMsgs, msg)
}

func (m *mockLogger) Error(msg string, err error, fields map[string]interface{}) {
	m.errorMsgs = append(m.errorMsgs, msg)
}

func (m *mockLogger) Debug(msg string, fields map[string]interface{}) {
	m.debugMsgs = append(m.debugMsgs, msg)
}

func (m *mockLogger) LogPayload(id string, payload interface{}) {
	m.infoMsgs = append(m.infoMsgs, fmt.Sprintf("Payload: %v", payload))
}

func (m *mockLogger) LogRequest(method string, url string, headers, params map[string]interface{}) {
	m.infoMsgs = append(m.infoMsgs, fmt.Sprintf("Request: %s %s", method, url))
}

func (m *mockLogger) LogResponse(statusCode int, body map[string]interface{}, duration time.Duration) {
	m.infoMsgs = append(m.infoMsgs, fmt.Sprintf("Response: %d", statusCode))
}

func (m *mockLogger) LogMetrics(component string, metrics map[string]interface{}) {
	m.infoMsgs = append(m.infoMsgs, fmt.Sprintf("Metrics: %s", component))
}

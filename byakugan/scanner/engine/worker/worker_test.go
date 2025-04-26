package worker

import (
	"context"
	"testing"
	"time"

	"github.com/haniz/byakugan/scanner/engine/detector"
	"github.com/haniz/byakugan/scanner/engine/logger"
	pb "github.com/haniz/byakugan/scanner/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockDetector struct {
	evidence *pb.Evidence
	err      error
}

// Implement Detector interface
func (m *mockDetector) DetectVulnerability(ctx context.Context, dc *detector.DetectionContext) (*pb.Evidence, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.evidence != nil {
		return m.evidence, nil
	}
	return &pb.Evidence{
		Data: map[string]string{
			"test": "value",
		},
		Validated: true,
		Timestamp: time.Now().Unix(),
	}, nil
}

type mockLogger struct {
	infoMsgs  []string
	errorMsgs []string
	debugMsgs []string
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

func (m *mockLogger) LogPayload(id string, payload interface{})                                    {}
func (m *mockLogger) LogRequest(method string, url string, headers, params map[string]interface{}) {}
func (m *mockLogger) LogResponse(statusCode int, body map[string]interface{}, duration time.Duration) {
}
func (m *mockLogger) LogMetrics(component string, metrics map[string]interface{}) {}
func (m *mockLogger) WithField(key string, value interface{}) logger.Logger       { return m }
func (m *mockLogger) WithFields(fields map[string]interface{}) logger.Logger      { return m }

// CreateTestDetectionContext tạo context mẫu cho test
func CreateTestDetectionContext() *detector.DetectionContext {
	return &detector.DetectionContext{
		Target: &pb.Target{
			Url:    "http://test.com",
			Method: "GET",
		},
		AuthContext: &pb.AuthContext{},
		InsertionPoint: &pb.InsertionPoint{
			Location: "query",
			Type:     "parameter",
		},
		Validation: &pb.Validation{},
	}
}

// Helper function to create test tasks
func CreateTestTask() *pb.ScanTask {
	return &pb.ScanTask{
		Id: "TEST-001",
		Target: &pb.Target{
			Url:    "http://test.com",
			Method: "GET",
		},
		AuthContext: &pb.AuthContext{
			Type: "none",
		},
		Payload: &pb.Payload{
			InsertionPoints: []*pb.InsertionPoint{
				{
					Location: "query",
					Type:     "parameter",
					Payloads: []string{"test_payload"},
				},
			},
		},
		RuleContext: &pb.RuleContext{
			Id:       "RULE-001",
			Category: "test",
			Severity: "low",
		},
		Validation: &pb.Validation{
			SuccessConditions: &pb.SuccessConditions{
				StatusCodes:      []int32{200},
				ResponsePatterns: []string{"success"},
			},
		},
	}
}

func TestWorker_ProcessTask(t *testing.T) {
	tests := []struct {
		name          string
		task          *pb.ScanTask
		mockEvidence  *pb.Evidence
		mockError     error
		expectError   bool
		expectSuccess bool
	}{
		{
			name: "successful jwt none detection",
			task: &pb.ScanTask{
				Id: "TEST-JWT-001",
				Target: &pb.Target{
					Url:    "http://test.com/auth",
					Method: "POST",
				},
				AuthContext: &pb.AuthContext{
					Type:  "jwt",
					Token: "eyJ0eXAiOiJKV1QiLCJhbGciOiJub25lIn0.eyJzdWIiOiIxMjM0In0.",
				},
				RuleContext: &pb.RuleContext{
					Id:       "BYAKUGAN-JWT-001",
					Category: "authentication",
					Severity: "critical",
				},
			},
			mockEvidence: &pb.Evidence{
				Data: map[string]string{
					"request":  "POST /auth HTTP/1.1\nAuthorization: Bearer eyJ0...",
					"response": "HTTP/1.1 200 OK\n{\"status\":\"success\"}",
				},
				Validated: true,
			},
			expectSuccess: true,
		},
		{
			name: "sql injection detection",
			task: &pb.ScanTask{
				Id: "TEST-SQLI-001",
				Target: &pb.Target{
					Url:    "http://test.com/users",
					Method: "GET",
				},
				Payload: &pb.Payload{
					InsertionPoints: []*pb.InsertionPoint{
						{
							Location: "query",
							Type:     "parameter",
							Payloads: []string{"' OR '1'='1"},
						},
					},
				},
				RuleContext: &pb.RuleContext{
					Id:       "BYAKUGAN-SQLI-001",
					Category: "injection",
					Severity: "high",
				},
			},
			mockEvidence: &pb.Evidence{
				Data: map[string]string{
					"request":  "GET /users?id=' OR '1'='1",
					"response": "HTTP/1.1 500\nSQLite error:",
				},
				Validated: true,
			},
			expectSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskChan := make(chan *pb.ScanTask, 1)
			resultChan := make(chan *pb.ScanResult, 1)
			logger := &mockLogger{}

			detector := &mockDetector{
				evidence: tt.mockEvidence,
				err:      tt.mockError,
			}

			w := NewWorker(1, taskChan, resultChan, detector, logger)
			result, err := w.ProcessTask(context.Background(), tt.task)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.task.Id, result.TaskId)
			assert.Equal(t, tt.expectSuccess, result.Success)

			if tt.expectSuccess {
				assert.NotNil(t, result.Evidence)
				assert.Equal(t, tt.mockEvidence.Data, result.Evidence.Data)
				assert.True(t, result.Evidence.Validated)
			}
		})
	}
}

func TestWorker_Metrics(t *testing.T) {
	w := NewWorker(1, nil, nil, &mockDetector{}, &mockLogger{})

	task := CreateTestTask()
	_, err := w.ProcessTask(context.Background(), task)
	require.NoError(t, err)

	assert.Equal(t, int64(1), w.metrics.TasksProcessed)
	assert.Greater(t, w.metrics.ProcessingTimeNs, int64(0))
}

func TestWorker_Start(t *testing.T) {
	tests := []struct {
		name        string
		taskID      string
		expectError bool
	}{
		{
			name:        "process task with string id",
			taskID:      "TEST-001",
			expectError: false,
		},
		{
			name:        "process task with uuid",
			taskID:      "123e4567-e89b-12d3-a456-426614174000",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskChan := make(chan *pb.ScanTask, 1)
			resultChan := make(chan *pb.ScanResult, 1)
			logger := &mockLogger{}

			w := NewWorker(1, taskChan, resultChan, &mockDetector{}, logger)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Start worker in goroutine
			go w.Start(ctx)

			// Send test task
			task := &pb.ScanTask{
				Id: tt.taskID,
				Target: &pb.Target{
					Url:    "http://test.com",
					Method: "GET",
				},
			}
			taskChan <- task

			// Wait for result
			select {
			case result := <-resultChan:
				assert.Equal(t, tt.taskID, result.TaskId)
				assert.False(t, tt.expectError)
			case <-time.After(time.Second):
				if !tt.expectError {
					t.Error("Timeout waiting for result")
				}
			}
		})
	}
}

package worker

import (
	"testing"

	pb "github.com/haniz/byakugan/scanner/proto"
	"github.com/stretchr/testify/assert"
)

func TestValidateTask(t *testing.T) {
	tests := []struct {
		name        string
		task        *pb.ScanTask
		expectError bool
	}{
		{
			name: "valid task",
			task: &pb.ScanTask{
				Id: "TEST-001",
				Target: &pb.Target{
					Url: "http://test.com",
				},
			},
			expectError: false,
		},
		{
			name:        "nil task",
			task:        nil,
			expectError: true,
		},
		{
			name: "missing target",
			task: &pb.ScanTask{
				Id: "TEST-002",
			},
			expectError: true,
		},
		{
			name: "missing url",
			task: &pb.ScanTask{
				Id:     "TEST-003",
				Target: &pb.Target{},
			},
			expectError: true,
		},
		{
			name: "missing insertion points",
			task: &pb.ScanTask{
				Id: "TEST-004",
				Target: &pb.Target{
					Url: "http://test.com",
				},
				Payload: &pb.Payload{},
			},
			expectError: true,
		},
		{
			name: "missing validation",
			task: &pb.ScanTask{
				Id: "TEST-005",
				Target: &pb.Target{
					Url: "http://test.com",
				},
				Payload: CreateTestPayload(),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTask(tt.task)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTaskConversion(t *testing.T) {
	pbTask := &pb.ScanTask{
		Id: "TEST-001",
		Target: &pb.Target{
			Url:    "http://test.com",
			Method: "POST",
		},
		AuthContext: &pb.AuthContext{
			Type:  "bearer",
			Token: "test.token",
		},
	}

	task, err := protoToTask(pbTask)
	assert.NoError(t, err)
	assert.Equal(t, pbTask.Id, task.ID)
	assert.Equal(t, pbTask.Target, task.Target)
	assert.Equal(t, pbTask.AuthContext, task.AuthContext)
	assert.NotNil(t, task.Evidence)
	assert.NotNil(t, task.Metadata)

	result := taskToResult(task)
	assert.Equal(t, task.ID, result.TaskId)
	assert.Equal(t, task.Evidence, result.Evidence)
	assert.Equal(t, task.Metadata, result.Metadata)
}

func CreateTestPayload() *pb.Payload {
	return &pb.Payload{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		InsertionPoints: []*pb.InsertionPoint{
			{
				Location: "body",
				Type:     "json",
				Payloads: []string{`{"test": "value"}`},
			},
		},
	}
}

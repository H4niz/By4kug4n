package worker

import (
	"fmt"
	"time"

	pb "github.com/haniz/byakugan/scanner/proto"
)

// Convert protobuf task to internal task
func protoToTask(pbTask *pb.ScanTask) (*Task, error) {
	if pbTask == nil {
		return nil, fmt.Errorf("nil task")
	}

	task := &Task{
		ID:          pbTask.Id,
		Target:      pbTask.Target,
		AuthContext: pbTask.AuthContext,
		Payload:     pbTask.Payload,
		Config:      pbTask.Config,
		Evidence: &pb.Evidence{
			Data: make(map[string]string),
		},
		Findings: make([]*pb.Finding, 0),
		Metadata: make(map[string]string),
	}

	return task, nil
}

// Convert internal task to protobuf result
func taskToResult(task *Task) *pb.ScanResult {
	return &pb.ScanResult{
		TaskId:    task.ID,
		Success:   task.Success,
		Findings:  task.Findings,
		Evidence:  task.Evidence,
		Timestamp: task.EndTime.Unix(),
		Metadata:  task.Metadata,
	}
}

func (t *Task) addEvidence(key string, value string) {
	if t.Evidence.Data == nil {
		t.Evidence.Data = make(map[string]string)
	}
	t.Evidence.Data[key] = value
}

func (t *Task) addFinding(finding *pb.Finding) {
	t.Findings = append(t.Findings, finding)
}

func (t *Task) markSuccess() {
	t.Success = true
}

func (t *Task) markFailure() {
	t.Success = false
}

// TaskConverter handles conversion between protobuf and internal types
type TaskConverter struct {
	logger WorkerLogger
}

func NewTaskConverter(l WorkerLogger) *TaskConverter {
	return &TaskConverter{logger: l}
}

func (c *TaskConverter) ProtoToInternal(task *pb.ScanTask) (*Task, error) {
	if task == nil {
		return nil, fmt.Errorf("nil task")
	}

	// Convert to internal task using protobuf types directly
	internalTask := &Task{
		ID:          task.Id,
		Target:      task.Target,      // Use pb.Target directly
		AuthContext: task.AuthContext, // Use pb.AuthContext directly
		Payload:     task.Payload,     // Use pb.Payload directly
		Config:      task.Config,
		Evidence: &pb.Evidence{
			Data: make(map[string]string),
		},
		Findings: make([]*pb.Finding, 0),
		Metadata: make(map[string]string),
	}

	return internalTask, nil
}

func (t *Task) ToResult() *pb.ScanResult {
	// Create scan result
	result := &pb.ScanResult{
		TaskId:    t.ID,
		Success:   t.Success,
		Findings:  t.Findings,
		Evidence:  t.Evidence, // Use Evidence directly
		Timestamp: time.Now().Unix(),
		Metadata:  t.Metadata,
	}
	return result
}

func (t *ScanTask) ToResult(evidence *pb.Evidence, err error) *pb.ScanResult {
	result := &pb.ScanResult{
		TaskId:    t.ID,
		Success:   err == nil,
		Evidence:  evidence, // Single Evidence
		Timestamp: time.Now().Unix(),
	}

	if err != nil {
		result.Metadata = map[string]string{
			"error": err.Error(),
		}
	}

	return result
}

// TaskHandler handles task operations
type TaskHandler struct {
	logger WorkerLogger
}

func NewTaskHandler(logger WorkerLogger) *TaskHandler {
	return &TaskHandler{
		logger: logger,
	}
}

type ScanTask struct {
	ID             string
	Target         *pb.Target
	AuthContext    *pb.AuthContext
	InsertionPoint *pb.InsertionPoint
	Payload        string
	Validation     *pb.Validation
	CreatedAt      time.Time
}

func validateTask(task *pb.ScanTask) error {
	if task == nil {
		return fmt.Errorf("nil task")
	}

	if task.Id == "" {
		return fmt.Errorf("missing task ID")
	}

	if task.Target == nil {
		return fmt.Errorf("missing target")
	}

	if task.Target.Url == "" {
		return fmt.Errorf("missing target URL")
	}

	return nil
}

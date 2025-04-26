package worker

import (
	"context"
	"sync"
	"time"

	"github.com/haniz/byakugan/scanner/engine/detector"
	pb "github.com/haniz/byakugan/scanner/proto"
)

// Task represents internal task state
type Task struct {
	ID          string
	Target      *pb.Target
	AuthContext *pb.AuthContext
	Payload     *pb.Payload
	Config      *pb.TaskConfig
	Evidence    *pb.Evidence
	Findings    []*pb.Finding
	Success     bool
	EndTime     time.Time
	Metadata    map[string]string
}

// PoolConfig defines worker pool configuration
type PoolConfig struct {
	Workers       int           `yaml:"workers"`
	QueueSize     int           `yaml:"queue_size"`
	MaxRetries    int           `yaml:"max_retries"`
	RetryDelay    time.Duration `yaml:"retry_delay"`
	MaxRetryDelay time.Duration `yaml:"max_retry_delay"`
	RetryBackoff  float64       `yaml:"retry_backoff"`
	Timeout       time.Duration `yaml:"timeout"`
}

// Pool manages worker pool
type Pool struct {
	workers    []*Worker
	taskChan   chan *pb.ScanTask
	resultChan chan *pb.ScanResult
	detector   detector.Detector
	logger     Logger
	metrics    *Metrics
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

// Worker represents a scanner worker
type Worker struct {
	id         int
	taskChan   <-chan *pb.ScanTask
	resultChan chan<- *pb.ScanResult
	detector   detector.Detector
	logger     Logger
	metrics    *Metrics
}

// Metrics tracks worker performance
type Metrics struct {
	TasksProcessed   int64
	TasksFailed      int64
	ProcessingTimeNs int64
}

// Logger defines worker logging interface
type Logger interface {
	Info(msg string, fields map[string]interface{})
	Error(msg string, err error, fields map[string]interface{})
	Debug(msg string, fields map[string]interface{})
}

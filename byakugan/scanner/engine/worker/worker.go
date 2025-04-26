package worker

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/haniz/byakugan/scanner/engine/detector"
	pb "github.com/haniz/byakugan/scanner/proto"
)

// NewWorker creates new worker
func NewWorker(id int, taskChan <-chan *pb.ScanTask, resultChan chan<- *pb.ScanResult, det detector.Detector, logger Logger) *Worker {
	return &Worker{
		id:         id,
		taskChan:   taskChan,
		resultChan: resultChan,
		detector:   det,
		logger:     logger,
		metrics:    &Metrics{},
	}
}

// Start starts the worker processing loop
func (w *Worker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-w.taskChan:
			if !ok {
				return
			}
			result, err := w.ProcessTask(ctx, task)
			if err != nil {
				// Fix: Use %s for string type
				w.logger.Error(fmt.Sprintf("Failed to process task %s", task.Id), err, map[string]interface{}{
					"task_id": task.Id,
				})
				continue
			}
			w.resultChan <- result
		}
	}
}

func (w *Worker) ProcessTask(ctx context.Context, task *pb.ScanTask) (*pb.ScanResult, error) {
	if err := validateTask(task); err != nil {
		return nil, err
	}

	startTime := time.Now()

	// Create detection context using detector package types
	detectionCtx := &detector.DetectionContext{
		Target:         task.Target,
		AuthContext:    task.AuthContext,
		InsertionPoint: task.Payload.InsertionPoints[0],
		Payload:        task.Payload.InsertionPoints[0].Payloads[0],
		Validation:     task.Validation,
	}

	// Execute detection
	evidence, err := w.detector.DetectVulnerability(ctx, detectionCtx)
	if err != nil {
		atomic.AddInt64(&w.metrics.TasksFailed, 1)
		w.logger.Error("Task execution failed", err, map[string]interface{}{
			"task_id": task.Id,
		})
		return &pb.ScanResult{
			TaskId:  task.Id,
			Success: false,
			Metadata: map[string]string{
				"error": err.Error(),
			},
		}, nil
	}

	result := &pb.ScanResult{
		TaskId:    task.Id,
		Success:   true,
		Evidence:  evidence,
		Timestamp: time.Now().Unix(),
		Metadata: map[string]string{
			"duration_ms": fmt.Sprintf("%d", time.Since(startTime).Milliseconds()),
		},
	}

	atomic.AddInt64(&w.metrics.TasksProcessed, 1)
	atomic.AddInt64(&w.metrics.ProcessingTimeNs, time.Since(startTime).Nanoseconds())

	return result, nil
}

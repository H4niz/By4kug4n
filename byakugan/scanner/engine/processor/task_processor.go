package processor

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/haniz/byakugan/scanner/engine/detector"
	"github.com/haniz/byakugan/scanner/engine/logger"
	"github.com/haniz/byakugan/scanner/engine/metrics"
	pb "github.com/haniz/byakugan/scanner/proto"
)

// TaskProcessor handles scan task processing
type TaskProcessor struct {
	config   *ProcessorConfig
	logger   logger.Logger
	detector detector.Detector // Changed from *detector.Detector
	metrics  *metrics.Metrics
}

// NewTaskProcessor creates a new task processor
func NewTaskProcessor(cfg *ProcessorConfig, logger logger.Logger) *TaskProcessor {
	if cfg == nil {
		cfg = DefaultProcessorConfig()
	}

	// Initialize detector config
	detectorCfg := detector.DefaultConfig()

	// Create detector instance
	d := detector.NewDetector(detectorCfg, logger)
	if d == nil {
		logger.Error("Failed to create detector", nil, nil)
		return nil
	}

	return &TaskProcessor{
		config:   cfg,
		logger:   logger,
		detector: d, // Use interface type
		metrics:  metrics.NewMetrics("task_processor"),
	}
}

// ProcessTask processes a scan task and returns results
func (p *TaskProcessor) ProcessTask(ctx context.Context, task *pb.ScanTask) (*pb.ScanResult, error) {
	if p.detector == nil {
		return nil, fmt.Errorf("detector not initialized")
	}

	// Start metrics tracking
	startTime := time.Now()
	defer func() {
		p.metrics.RecordDuration(time.Since(startTime))
		p.metrics.IncrementCounter("tasks_processed")
	}()

	if err := p.validateTask(task); err != nil {
		return nil, err
	}

	p.logger.Info("Processing scan task", map[string]interface{}{
		"task_id":   task.Id,
		"target":    task.Target.Url,
		"method":    task.Target.Method,
		"auth_type": task.AuthContext.Type,
	})

	findings := make([]*pb.Finding, 0)
	for _, point := range task.Payload.InsertionPoints {
		finding, err := p.processInsertionPoint(ctx, task, point)
		if err != nil {
			p.logger.Error("Failed to process insertion point", err, map[string]interface{}{
				"location": point.Location,
				"type":     point.Type,
			})
			continue
		}
		if finding != nil {
			findings = append(findings, finding)
		}
	}

	result := &pb.ScanResult{
		TaskId:    task.Id,
		Success:   len(findings) > 0,
		Findings:  findings,
		Timestamp: time.Now().Unix(),
		Metadata: map[string]string{
			"scanner_version": "1.0.0",
			"scan_type":       task.RuleContext.Category,
		},
	}

	return result, nil
}

func (p *TaskProcessor) processInsertionPoint(ctx context.Context, task *pb.ScanTask, point *pb.InsertionPoint) (*pb.Finding, error) {
	if point.Type == "jwt_none" {
		evidence, err := p.detector.DetectVulnerability(ctx, &detector.DetectionContext{
			Target:         task.Target,
			AuthContext:    task.AuthContext,
			InsertionPoint: point,
			Validation:     task.Validation,
		})
		if err != nil {
			p.logger.Error("JWT none detection failed", err, map[string]interface{}{
				"task_id":  task.Id,
				"location": point.Location,
			})
			return nil, err
		}

		if evidence != nil {
			return &pb.Finding{
				Id:       uuid.New().String(),
				RuleId:   task.RuleContext.Id,
				Severity: task.RuleContext.Severity,
				Title:    "JWT none algorithm vulnerability",
				Details:  "Successfully bypassed JWT signature verification using 'none' algorithm",
				Evidence: []*pb.Evidence{evidence}, // Fix: Wrap single evidence in slice
			}, nil
		}
	}
	return p.defaultProcessInsertionPoint(ctx, task, point)
}

func (p *TaskProcessor) defaultProcessInsertionPoint(ctx context.Context, task *pb.ScanTask, point *pb.InsertionPoint) (*pb.Finding, error) {
	for _, payload := range point.Payloads {
		evidence, err := p.detector.DetectVulnerability(ctx, &detector.DetectionContext{
			Target:         task.Target,
			AuthContext:    task.AuthContext,
			InsertionPoint: point,
			Payload:        payload,
			Validation:     task.Validation,
		})

		if err != nil {
			continue
		}

		if evidence.Validated {
			return &pb.Finding{
				Id:       generateUUID(),
				RuleId:   task.RuleContext.Id,
				Severity: task.RuleContext.Severity,
				Title:    fmt.Sprintf("Vulnerability found in %s", point.Location),
				Details:  fmt.Sprintf("Successfully exploited using payload in %s", point.Location),
				Evidence: []*pb.Evidence{evidence},
			}, nil
		}
	}

	return nil, nil
}

func (p *TaskProcessor) validateTask(task *pb.ScanTask) error {
	if task.Id == "" {
		return fmt.Errorf("task ID is required")
	}
	if task.Target == nil || task.Target.Url == "" {
		return fmt.Errorf("target URL is required")
	}
	if task.RuleContext == nil || task.RuleContext.Id == "" {
		return fmt.Errorf("rule context is required")
	}
	return nil
}

func generateUUID() string {
	return uuid.New().String()
}

// GetMetrics retrieves the metrics for the task processor
func (p *TaskProcessor) GetMetrics() map[string]interface{} {
	return p.metrics.GetMetrics()
}

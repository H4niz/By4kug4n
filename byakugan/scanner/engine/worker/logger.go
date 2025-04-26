package worker

import (
	"github.com/haniz/byakugan/scanner/engine/logger"
)

// WorkerLogger extends logger.Logger with worker-specific methods
type WorkerLogger interface {
	logger.Logger
	LogWorkerEvent(eventType string, fields map[string]interface{})
	LogTaskProgress(taskID string, progress float64, status string)
	LogTaskResult(taskID string, success bool, findings int)
}

type workerLogger struct {
	logger.Logger
}

func NewWorkerLogger(baseLogger logger.Logger) WorkerLogger {
	return &workerLogger{Logger: baseLogger}
}

func (l *workerLogger) LogWorkerEvent(eventType string, fields map[string]interface{}) {
	allFields := make(map[string]interface{})
	for k, v := range fields {
		allFields[k] = v
	}
	allFields["event_type"] = eventType
	l.Info("Worker event", allFields)
}

func (l *workerLogger) LogTaskProgress(taskID string, progress float64, status string) {
	l.Info("Task progress", map[string]interface{}{
		"task_id":  taskID,
		"progress": progress,
		"status":   status,
	})
}

func (l *workerLogger) LogTaskResult(taskID string, success bool, findings int) {
	l.Info("Task result", map[string]interface{}{
		"task_id":  taskID,
		"success":  success,
		"findings": findings,
	})
}

func (l *workerLogger) Error(msg string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["error"] = err.Error()
	l.Logger.Error(msg, err, fields)
}

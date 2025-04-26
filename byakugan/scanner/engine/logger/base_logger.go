package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type BaseLogger struct {
	mu     sync.RWMutex
	config LogConfig
	fields map[string]interface{}
}

func NewBaseLogger(cfg LogConfig) (Logger, error) {
	// Create logger directory if needed
	if err := os.MkdirAll(filepath.Dir(cfg.FilePath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	l := &BaseLogger{
		config: cfg,
		fields: make(map[string]interface{}),
	}

	l.fields["component"] = cfg.Component
	return l, nil
}

// Core logging methods
func (l *BaseLogger) Info(msg string, fields map[string]interface{}) {
	l.log(INFO, msg, nil, fields)
}

func (l *BaseLogger) Error(msg string, err error, fields map[string]interface{}) {
	l.log(ERROR, msg, err, fields)
}

func (l *BaseLogger) Debug(msg string, fields map[string]interface{}) {
	if l.config.Level == DEBUG {
		l.log(DEBUG, msg, nil, fields)
	}
}

func (l *BaseLogger) log(level LogLevel, msg string, err error, fields map[string]interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	entry := l.createEntry(level, msg, err, fields)
	l.output(entry)
}

func (l *BaseLogger) output(entry *LogEntry) {
	data, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling log entry: %v\n", err)
		return
	}

	// Write to file if configured
	if l.config.FilePath != "" {
		if f, err := os.OpenFile(l.config.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			defer f.Close()
			fmt.Fprintln(f, string(data))
		}
	}

	// Write to console if enabled
	if l.config.Console {
		fmt.Fprintln(os.Stdout, string(data))
	}
}

// Standard reporting methods
func (l *BaseLogger) LogRequest(method string, url string, headers, params map[string]interface{}) {
	l.Info("Request", map[string]interface{}{
		"method":  method,
		"url":     url,
		"headers": headers,
		"params":  params,
	})
}

func (l *BaseLogger) LogResponse(statusCode int, body map[string]interface{}, duration time.Duration) {
	l.Info("Response", map[string]interface{}{
		"status_code": statusCode,
		"body":        body,
		"duration":    duration.String(),
	})
}

func (l *BaseLogger) LogMetrics(component string, metrics map[string]interface{}) {
	l.Info("Metrics", mergeFields(map[string]interface{}{
		"component": component,
	}, metrics))
}

// Add missing methods to implement Logger interface
func (l *BaseLogger) WithField(key string, value interface{}) Logger {
	return l.WithFields(map[string]interface{}{key: value})
}

func (l *BaseLogger) WithFields(fields map[string]interface{}) Logger {
	l.mu.RLock()
	newFields := make(map[string]interface{}, len(l.fields)+len(fields))
	for k, v := range l.fields {
		newFields[k] = v
	}
	l.mu.RUnlock()

	for k, v := range fields {
		newFields[k] = v
	}

	return &BaseLogger{
		config: l.config,
		fields: newFields,
	}
}

func (l *BaseLogger) LogPayload(id string, payload interface{}) {
	l.Info("Payload processed", map[string]interface{}{
		"payload_id": id,
		"payload":    payload,
	})
}

func (l *BaseLogger) Close() error {
	// Add any cleanup logic here
	return nil
}

// Helper method to create structured log entry
func (l *BaseLogger) createEntry(level LogLevel, msg string, err error, fields map[string]interface{}) *LogEntry {
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		Error:     err,
		Component: l.config.Component,
		Fields:    make(map[string]interface{}, len(l.fields)+len(fields)),
	}

	// Add base fields
	for k, v := range l.fields {
		entry.Fields[k] = v
	}

	// Add additional fields
	for k, v := range fields {
		entry.Fields[k] = v
	}

	return entry
}

package logger

import (
	"sync"
	"testing"
	"time"
)

type TestLogger struct {
	t       *testing.T
	entries []*LogEntry
	mu      sync.RWMutex
}

// Interface for test logging operations
type TestLoggerInterface interface {
	Logger
	GetEntries() []*LogEntry
	FindEntry(level LogLevel, msg string) *LogEntry
	ClearEntries()
}

func NewTestLogger(t *testing.T) TestLoggerInterface {
	return &TestLogger{
		t:       t,
		entries: make([]*LogEntry, 0),
	}
}

func (l *TestLogger) Info(msg string, fields map[string]interface{}) {
	l.addEntry(INFO, msg, nil, fields)
}

func (l *TestLogger) Error(msg string, err error, fields map[string]interface{}) {
	l.addEntry(ERROR, msg, err, fields)
}

func (l *TestLogger) Debug(msg string, fields map[string]interface{}) {
	l.addEntry(DEBUG, msg, nil, fields)
}

func (l *TestLogger) addEntry(level LogLevel, msg string, err error, fields map[string]interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		Error:     err,
		Fields:    fields,
	}
	l.entries = append(l.entries, entry)
}

// Add log entry getters/finders
func (l *TestLogger) GetEntries() []*LogEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.entries
}

func (l *TestLogger) FindEntry(level LogLevel, msg string) *LogEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	for _, entry := range l.entries {
		if entry.Level == level && entry.Message == msg {
			return entry
		}
	}
	return nil
}

func (l *TestLogger) ClearEntries() {
	l.mu.Lock()
	l.entries = make([]*LogEntry, 0)
	l.mu.Unlock()
}

// Add new method for log sequence verification
func (l *TestLogger) AssertLogSequence(t *testing.T, expectedEntries []LogEntry) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if len(expectedEntries) > len(l.entries) {
		t.Errorf("Expected %d log entries but got %d", len(expectedEntries), len(l.entries))
		return
	}

	for i, expected := range expectedEntries {
		actual := l.entries[i]
		if expected.Level != actual.Level || expected.Message != actual.Message {
			t.Errorf("Log entry %d mismatch:\nexpected: %+v\nactual: %+v", i, expected, actual)
		}
	}
}

// Implement remaining Logger interface methods
func (l *TestLogger) WithField(key string, value interface{}) Logger {
	return l.WithFields(map[string]interface{}{key: value})
}

func (l *TestLogger) WithFields(fields map[string]interface{}) Logger {
	return l
}

func (l *TestLogger) LogMetrics(component string, metrics map[string]interface{}) {
	l.Info("Metrics", map[string]interface{}{"component": component, "metrics": metrics})
}

func (l *TestLogger) LogPayload(id string, payload interface{}) {
	l.Info("Payload", map[string]interface{}{"id": id, "payload": payload})
}

func (l *TestLogger) LogRequest(method string, url string, headers map[string]interface{}, params map[string]interface{}) {
	l.Info("Request", map[string]interface{}{
		"method":  method,
		"url":     url,
		"headers": headers,
		"params":  params,
	})
}

func (l *TestLogger) LogResponse(statusCode int, body map[string]interface{}, duration time.Duration) {
	l.Info("Response", map[string]interface{}{
		"status_code": statusCode,
		"body":        body,
		"duration":    duration,
	})
}

package logger

import (
	"time"
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Error     error                  `json:"error,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Component string                 `json:"component,omitempty"`
	TraceID   string                 `json:"trace_id,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"` // Add Data field
}

// LogConfig defines standardized logger configuration
type LogConfig struct {
	Component  string        `yaml:"component"`
	Level      LogLevel      `yaml:"level"`
	LogPath    string        `yaml:"log_path"`
	FilePath   string        `yaml:"file_path"`
	MaxSize    int64         `yaml:"max_size"`
	MaxAge     time.Duration `yaml:"max_age"`
	Compress   bool          `yaml:"compress"`
	BufferSize int           `yaml:"buffer_size"`
	Console    bool          `yaml:"console"`
	AsyncWrite bool          `yaml:"async_write"`
	FormatJSON bool          `yaml:"format_json"`
	QueueSize  int           `yaml:"queue_size"`
}

// Logger defines the standardized logging interface
type Logger interface {
	Info(msg string, fields map[string]interface{})
	Error(msg string, err error, fields map[string]interface{})
	Debug(msg string, fields map[string]interface{})
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger

	// Additional standardized methods
	LogMetrics(component string, metrics map[string]interface{})
	LogPayload(id string, payload interface{})
	LogRequest(method string, url string, headers map[string]interface{}, params map[string]interface{})
	LogResponse(statusCode int, body map[string]interface{}, duration time.Duration)
}

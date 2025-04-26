package logger

import (
	"fmt"
	"path/filepath"
)

// CreateLogger creates a new logger instance with proper error handling
func CreateLogger(component string, cfg *LogConfig) (Logger, error) {
	if cfg == nil {
		cfg = &LogConfig{
			Component:  component,
			Level:      INFO,
			LogPath:    "logs",
			FilePath:   filepath.Join("logs", component+".log"),
			BufferSize: 4096,
			Console:    true,
		}
	}

	baseLogger, err := NewBaseLogger(*cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create base logger: %w", err)
	}

	// Wrap with adapter to ensure interface compliance
	return NewLoggerAdapter(baseLogger), nil
}

// CreateLoggerWithDefaults creates a logger with default settings
func CreateLoggerWithDefaults(component string) (Logger, error) {
	cfg := &LogConfig{
		Component:  component,
		Level:      INFO,
		LogPath:    filepath.Join("logs", component),
		FilePath:   filepath.Join("logs", component, "service.log"),
		MaxSize:    100 * 1024 * 1024, // 100MB
		MaxAge:     24 * 60,           // 24 hours
		Compress:   true,
		BufferSize: 4096,
		Console:    true,
		AsyncWrite: true,
		FormatJSON: true,
	}

	return CreateLogger(component, cfg)
}

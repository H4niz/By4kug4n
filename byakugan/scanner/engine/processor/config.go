package processor

import (
	"fmt"
	"time"
)

// ProcessorConfig holds processor configuration
type ProcessorConfig struct {
	BatchSize      int           `yaml:"batch_size" json:"batch_size"`
	FlushInterval  time.Duration `yaml:"flush_interval" json:"flush_interval"`
	RetryAttempts  int           `yaml:"retry_attempts" json:"retry_attempts"`
	RetryDelay     time.Duration `yaml:"retry_delay" json:"retry_delay"`
	LogLevel       string        `yaml:"log_level" json:"log_level"`
	LogPath        string        `yaml:"log_path" json:"log_path"`
	MetricsEnabled bool          `yaml:"metrics_enabled" json:"metrics_enabled"`
}

// DefaultProcessorConfig returns default processor configuration
func DefaultProcessorConfig() *ProcessorConfig {
	return &ProcessorConfig{
		BatchSize:      100,
		FlushInterval:  5 * time.Second,
		RetryAttempts:  3,
		RetryDelay:     time.Second,
		LogLevel:       "info",
		LogPath:        "logs/processor.log",
		MetricsEnabled: true,
	}
}

// Validate validates processor configuration
func (c *ProcessorConfig) Validate() error {
	if c.BatchSize <= 0 {
		return fmt.Errorf("batch size must be positive")
	}
	if c.FlushInterval <= 0 {
		return fmt.Errorf("flush interval must be positive")
	}
	if c.RetryAttempts < 0 {
		return fmt.Errorf("retry attempts cannot be negative")
	}
	return nil
}

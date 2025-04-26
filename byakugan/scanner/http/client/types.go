package client

import "time"

// RetryOnConfig defines conditions for retry
type RetryOnConfig struct {
	StatusCodes   []int `json:"status_codes" yaml:"status_codes"`
	NetworkErrors bool  `json:"network_errors" yaml:"network_errors"`
	Timeouts      bool  `json:"timeouts" yaml:"timeouts"`
}

// RetryConfig represents retry settings
type RetryConfig struct {
	MaxAttempts       int           `json:"max_attempts" yaml:"max_attempts"`
	InitialDelay      time.Duration `json:"initial_delay" yaml:"initial_delay"`
	MaxDelay          time.Duration `json:"max_delay" yaml:"max_delay"`
	BackoffMultiplier float64       `json:"backoff_multiplier" yaml:"backoff_multiplier"`
	RetryOn           RetryOnConfig `json:"retry_on" yaml:"retry_on"`
}

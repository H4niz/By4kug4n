package worker

import "time"

// DefaultPoolConfig returns default worker pool configuration
func DefaultPoolConfig() *PoolConfig {
	return &PoolConfig{
		Workers:       5,
		QueueSize:     100,
		MaxRetries:    3,
		MaxRetryDelay: time.Second * 30,
		RetryBackoff:  2.0,
	}
}

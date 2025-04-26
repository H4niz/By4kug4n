package detector

import "time"

// Config represents detector configuration
type Config struct {
	MaxPatterns     int           `yaml:"max_patterns"`
	MaxBodySize     int64         `yaml:"max_body_size"`
	MaxHeaderSize   int64         `yaml:"max_header_size"`
	MaxRequestCount int64         `yaml:"max_request_count"`
	Timeout         time.Duration `yaml:"timeout"`
	RequestTimeout  time.Duration `yaml:"request_timeout"`
	MaxRetries      int           `yaml:"max_retries"`
	RetryDelay      time.Duration `yaml:"retry_delay"`
	FollowRedirects bool          `yaml:"follow_redirects"`
	MaxRedirects    int           `yaml:"max_redirects"`
	ValidateSSL     bool          `yaml:"validate_ssl"`
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		MaxPatterns:     100,
		MaxBodySize:     10 * 1024 * 1024, // 10MB
		MaxHeaderSize:   64 * 1024,        // 64KB
		MaxRequestCount: 1000,
		Timeout:         30 * time.Second,
		RequestTimeout:  30 * time.Second,
		MaxRetries:      3,
		RetryDelay:      time.Second,
		FollowRedirects: true,
		MaxRedirects:    5,
		ValidateSSL:     false,
	}
}

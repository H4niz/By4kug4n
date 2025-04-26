package client

import "time"

// Config represents HTTP client configuration
type Config struct {
	UserAgent      string        `json:"user_agent"`
	FollowRedirect bool          `json:"follow_redirects"`
	MaxRedirects   int           `json:"max_redirects"`
	VerifyCerts    bool          `json:"verify_certs"`
	Compression    bool          `json:"compression"`
	KeepAlive      bool          `json:"keep_alive"`
	Timeout        TimeoutConfig `json:"timeout"`
	TLS            TLSConfig     `json:"tls"`
	RetryWait      time.Duration `json:"retry_wait"`
	Retry          RetryConfig   `json:"retry"`
}

// TimeoutConfig represents timeout settings
type TimeoutConfig struct {
	Connect time.Duration `json:"connect"`
	Read    time.Duration `json:"read"`
	Write   time.Duration `json:"write"`
}

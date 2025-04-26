package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

// Config represents the main scanner configuration
type Config struct {
	Nodes         []NodeConfiguration `yaml:"nodes"`
	Communication CommunicationConfig `yaml:"communication"`
	Logging       LoggingConfig       `yaml:"logging"`
	Proxy         ProxyConfig         `yaml:"proxy"`
	Retry         RetryConfig         `yaml:"retry"`
	GRPC          struct {
		Port           int           `yaml:"port"`
		MaxMessageSize int           `yaml:"max_message_size"`
		MaxWorkers     int           `yaml:"max_workers"`
		Timeout        time.Duration `yaml:"timeout"`
	} `yaml:"grpc"`
	Engine struct {
		Workers    int           `yaml:"workers"`
		MaxRetries int           `yaml:"max_retries"`
		Timeout    time.Duration `yaml:"timeout"`
	} `yaml:"engine"`
	HTTP struct {
		MaxRedirects   int           `yaml:"max_redirects"`
		RequestTimeout time.Duration `yaml:"request_timeout"`
		RetryDelay     time.Duration `yaml:"retry_delay"`
	} `yaml:"http"`
}

// Add method to Config structidle_timeout
func (c *Config) LoadFromFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	return nil
}

// NodeConfiguration represents configuration for a single scanner node
type NodeConfiguration struct {
	Node         NodeConfig         `yaml:"node"`
	Performance  PerformanceConfig  `yaml:"performance"`
	WorkerPool   WorkerPoolConfig   `yaml:"worker_pool"`
	HTTPClient   HTTPClientConfig   `yaml:"http_client"`
	RateLimiting RateLimitingConfig `yaml:"rate_limiting"`
}

type CommunicationConfig struct {
	GRPC GRPCConfig `yaml:"grpc"`
}

// NodeConfig represents scanner node configuration
type NodeConfig struct {
	ID           string   `yaml:"id"`
	Name         string   `yaml:"name"`
	Region       string   `yaml:"region"`
	Tags         []string `yaml:"tags"`
	Capabilities []string `yaml:"capabilities"`
}

// PerformanceConfig represents performance limits
type PerformanceConfig struct {
	CPULimit    int           `yaml:"cpu_limit"`
	MemoryLimit int           `yaml:"memory_limit"`
	Network     NetworkConfig `yaml:"network"`
}

type NetworkConfig struct {
	BandwidthLimit   int `yaml:"bandwidth_limit"`
	ConnectionsLimit int `yaml:"connections_limit"`
}

// CalculateMaxWorkers calculates maximum workers based on resources
func (p *PerformanceConfig) CalculateMaxWorkers() int {
	// Base calculation on CPU and connection limits
	cpuBasedLimit := p.CPULimit / 2                   // Assume each worker uses ~2% CPU
	connBasedLimit := p.Network.ConnectionsLimit / 10 // Allow 10 connections per worker

	// Take the more restrictive limit
	maxWorkers := min(cpuBasedLimit, connBasedLimit)

	// Apply absolute limits
	switch {
	case maxWorkers > 100:
		maxWorkers = 100 // Hard upper limit
	case maxWorkers < 1:
		maxWorkers = 1 // Minimum of 1 worker
	}

	return maxWorkers
}

// CalculateQueueSizeLimits calculates min and max queue sizes based on resources
func (p *PerformanceConfig) CalculateQueueSizeLimits(workerCount int) (min, max int) {
	// Minimum queue size is 2x worker count
	min = workerCount * 2

	// Maximum queue size based on memory
	// Assume each queued task takes ~10KB memory
	max = (p.MemoryLimit * 1024) / 10

	// Apply absolute limits
	if max > 10000 {
		max = 10000 // Hard upper limit
	}
	if max < min {
		max = min // Ensure max is not less than min
	}

	return min, max
}

// helper function since Go 1.21's math.Min works only with float64
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// WorkerPoolConfig represents worker pool settings
type WorkerPoolConfig struct {
	Size                int           `yaml:"size"`
	QueueSize           int           `yaml:"queue_size"`
	BatchSize           int           `yaml:"batch_size"`
	RetryCount          int           `yaml:"retry_count"`
	RetryDelay          time.Duration `yaml:"retry_delay"`
	HealthCheckInterval time.Duration `yaml:"health_check_interval"`
	MetricsEnabled      bool          `yaml:"metrics_enabled"`
	ShutdownTimeout     time.Duration `yaml:"shutdown_timeout"`
}

// HTTPClientConfig represents HTTP client settings
type HTTPClientConfig struct {
	UserAgent        string `yaml:"user_agent"`
	FollowRedirects  bool   `yaml:"follow_redirects"`
	MaxRedirects     int    `yaml:"max_redirects"`
	KeepAlive        bool   `yaml:"keep_alive"`
	KeepAliveTimeout int    `yaml:"keep_alive_idle_timeout"`
	Timeout          struct {
		Connect time.Duration `yaml:"connect"`
		Read    time.Duration `yaml:"read"`
		Write   time.Duration `yaml:"write"`
	} `yaml:"timeout"`
	TLS struct {
		VerifyCerts bool   `yaml:"verify_certs"`
		MinVersion  string `yaml:"min_version"`
	} `yaml:"tls"`
	Compression    bool `yaml:"compression"`
	ConnectionPool struct {
		MaxIdleConns        int           `yaml:"max_idle_conns"`
		MaxConnsPerHost     int           `yaml:"max_conns_per_host"`
		IdleConnTimeout     time.Duration `yaml:"idle_conn_timeout"`
		TLSHandshakeTimeout time.Duration `yaml:"tls_handshake_timeout"`
	} `yaml:"connection_pool"`
	RetryPolicy struct {
		MaxAttempts     int           `yaml:"max_attempts"`
		InitialInterval time.Duration `yaml:"initial_interval"`
		MaxInterval     time.Duration `yaml:"max_interval"`
		Multiplier      float64       `yaml:"multiplier"`
	} `yaml:"retry_policy"`
	CircuitBreaker struct {
		Enabled   bool          `yaml:"enabled"`
		Threshold int           `yaml:"threshold"`
		Interval  time.Duration `yaml:"interval"`
		Timeout   time.Duration `yaml:"timeout"`
	} `yaml:"circuit_breaker"`
}

// RateLimitingConfig represents rate limiting settings
type RateLimitingConfig struct {
	Enabled        bool    `yaml:"enabled"`
	Strategy       string  `yaml:"strategy"`
	InitialRate    int     `yaml:"initial_rate"`
	MaxRate        int     `yaml:"max_rate"`
	MinRate        int     `yaml:"min_rate"`
	BackoffFactor  float64 `yaml:"backoff_factor"`
	RecoveryFactor float64 `yaml:"recovery_factor"`
}

// RetryConfig represents retry mechanism settings
type RetryConfig struct {
	MaxAttempts       int           `yaml:"max_attempts"`
	InitialDelay      time.Duration `yaml:"initial_delay"`
	MaxDelay          time.Duration `yaml:"max_delay"`
	BackoffMultiplier float64       `yaml:"backoff_multiplier"`
	RetryOn           struct {
		StatusCodes   []int `yaml:"status_codes"`
		NetworkErrors bool  `yaml:"network_errors"`
		Timeouts      bool  `yaml:"timeouts"`
	} `yaml:"retry_on"`
}

// ProxyConfig represents proxy settings
type ProxyConfig struct {
	Enabled            bool          `yaml:"enabled"`
	Proxies            []string      `yaml:"proxies"` // Array of proxy URLs
	CheckInterval      time.Duration `yaml:"check_interval"`
	Timeout            time.Duration `yaml:"timeout"`
	BlacklistThreshold int           `yaml:"blacklist_threshold"`
	BlacklistDuration  time.Duration `yaml:"blacklist_duration"`
	Rotation           struct {
		Strategy  string        `yaml:"strategy"` // random, round-robin
		Interval  time.Duration `yaml:"interval"`
		PerWorker bool          `yaml:"per_worker"` // Dedicated proxies per worker
	} `yaml:"rotation"`
}

// LoggingConfig defines configuration for logging
type LoggingConfig struct {
	Level          string                  `yaml:"level"`
	LogPath        string                  `yaml:"log_path"`
	MaxSize        int                     `yaml:"max_size"`
	MaxFiles       int                     `yaml:"max_files"`
	Compress       bool                    `yaml:"compress"`
	IncludeTraceID bool                    `yaml:"include_trace_id"`
	Rotation       LogRotationConfig       `yaml:"rotation"`
	Components     map[string]ComponentLog `yaml:"components"`
	BufferSize     int                     `yaml:"buffer_size"`
	AsyncWrite     bool                    `yaml:"async_write"`
	QueueSize      int                     `yaml:"queue_size"`
	FlushInterval  time.Duration           `yaml:"flush_interval"`
	Trace          struct {
		Enabled      bool    `yaml:"enabled"`
		SamplingRate float64 `yaml:"sampling_rate"`
	} `yaml:"trace"`
	Metrics struct {
		Enabled  bool          `yaml:"enabled"`
		Interval time.Duration `yaml:"interval"`
		Prefix   string        `yaml:"prefix"`
	} `yaml:"metrics"`
}

// LogRotationConfig defines log rotation settings
type LogRotationConfig struct {
	MaxSizeMB int64         `yaml:"max_size_mb"`
	MaxFiles  int           `yaml:"max_files"`
	Compress  bool          `yaml:"compress"`
	MaxAge    time.Duration `yaml:"max_age"`
}

// ComponentLog defines component specific logging settings
type ComponentLog struct {
	LogPath string `yaml:"log_path"`
	Level   string `yaml:"level"`
}

// TimeoutConfig represents timeout settings
type TimeoutConfig struct {
	Connect time.Duration `yaml:"connect"`
	Read    time.Duration `yaml:"read"`
	Write   time.Duration `yaml:"write"`
}

// DetectorConfig represents detector settings
type DetectorConfig struct {
	MaxBodySize      int64   `yaml:"max_body_size"`
	EnableScreenshot bool    `yaml:"enable_screenshots"`
	Confidence       float64 `yaml:"confidence_threshold"`
}

// PoolConfig represents connection pool settings
type PoolConfig struct {
	MaxSize     int           `yaml:"max_size"`
	MinIdle     int           `yaml:"min_idle"`
	MaxIdle     int           `yaml:"max_idle"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

// GRPCConfig represents gRPC settings
type GRPCConfig struct {
	Address        string          `yaml:"address"`
	MaxMessageSize int             `yaml:"max_message_size"`
	KeepAlive      KeepAliveConfig `yaml:"keep_alive"` // Note the yaml tag matches scanner.yaml
	ConnectionPool PoolConfig      `yaml:"connection_pool"`
	TLS            TLSConfig       `yaml:"tls"`
}

// KeepAliveConfig represents gRPC keepalive settings
type KeepAliveConfig struct {
	Enabled             bool          `yaml:"enabled"`
	Time                time.Duration `yaml:"time"`
	Timeout             time.Duration `yaml:"timeout"`
	IdleTimeout         time.Duration `yaml:"idle_timeout"`
	MaxAge              time.Duration `yaml:"max_age"`
	GracePeriod         time.Duration `yaml:"grace_period"`
	PermitWithoutStream bool          `yaml:"permit_without_stream"`
}

// TLSConfig represents gRPC TLS settings
type TLSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

// ProcessorConfig represents processor settings
type ProcessorConfig struct {
	BatchSize       int           `yaml:"batch_size"`
	FlushInterval   time.Duration `yaml:"flush_interval"`
	RetryAttempts   int           `yaml:"retry_attempts"`
	RetryDelay      time.Duration `yaml:"retry_delay"`
	MetricsEnabled  bool          `yaml:"metrics_enabled"`
	MetricsInterval time.Duration `yaml:"metrics_interval"`
}

// CommsConfig represents communication configuration
type CommsConfig struct {
	GRPC    GRPCConfig    `yaml:"grpc"`
	Logging LoggingConfig `yaml:"logging"`
}

var (
	loadOnce sync.Once
	instance *Config
	loadErr  error
)

// Load loads configuration from a YAML file
func Load(path string) (*Config, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config path: %w", err)
	}

	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Validate configuration
	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// LoadWithSync loads configuration synchronously
func LoadWithSync(path string) (*Config, error) {
	mu := &sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()

	return Load(path)
}

// GetInstance returns the singleton config instance
func GetInstance() *Config {
	return instance
}

// validateConfig performs validation of the configuration
func validateConfig(cfg *Config) error {
	for _, node := range cfg.Nodes {
		// Validate performance config
		if node.Performance.CPULimit <= 0 || node.Performance.CPULimit > 100 {
			return fmt.Errorf("invalid CPU limit: must be between 1-100, got %d",
				node.Performance.CPULimit)
		}
		if node.Performance.MemoryLimit <= 0 {
			return fmt.Errorf("invalid memory limit: must be greater than 0, got %d",
				node.Performance.MemoryLimit)
		}
	}
	return nil
}

// validate performs validation of the configuration
func (c *Config) validate() error {
	if len(c.Nodes) == 0 {
		return fmt.Errorf("missing required field: at least one node configuration required")
	}

	for i, node := range c.Nodes {
		if err := node.validate(); err != nil {
			return fmt.Errorf("invalid node config at index %d: %w", i, err)
		}
	}

	return nil
}

// validate performs validation of node configuration
func (nc *NodeConfiguration) validate() error {
	// Check required fields
	if nc.Node.ID == "" {
		return fmt.Errorf("missing required field: node.id")
	}

	// Validate worker pool
	if nc.WorkerPool.Size <= 0 {
		return fmt.Errorf("invalid worker pool size: must be positive")
	}

	// Validate performance limits
	if nc.Performance.CPULimit <= 0 || nc.Performance.CPULimit > 100 {
		return fmt.Errorf("invalid CPU limit: must be between 1-100")
	}

	if nc.Performance.MemoryLimit <= 0 {
		return fmt.Errorf("invalid memory limit: must be greater than 0")
	}

	return nil
}

// GetGRPCConfig retrieves the centralized GRPC configuration
func (c *Config) GetGRPCConfig() *GRPCConfig {
	return &c.Communication.GRPC
}

// validate performs deeper validation checks
func validateDeep(cfg *Config) []error {
	var errors []error

	// Validate node configurations
	for i, node := range cfg.Nodes {
		// Validate performance settings
		if node.Performance.CPULimit <= 0 || node.Performance.CPULimit > 100 {
			errors = append(errors, fmt.Errorf("node %d: CPU limit must be between 1-100", i))
		}

		if node.Performance.MemoryLimit <= 0 {
			errors = append(errors, fmt.Errorf("node %d: memory limit must be positive", i))
		}

		// Validate worker pool settings
		if node.WorkerPool.Size <= 0 {
			errors = append(errors, fmt.Errorf("node %d: worker pool size must be positive", i))
		}

		if node.WorkerPool.QueueSize <= 0 {
			errors = append(errors, fmt.Errorf("node %d: queue size must be positive", i))
		}
	}

	// Validate GRPC config
	if cfg.Communication.GRPC.Address == "" {
		errors = append(errors, fmt.Errorf("GRPC address must be specified"))
	}

	if cfg.Communication.GRPC.MaxMessageSize <= 0 {
		errors = append(errors, fmt.Errorf("GRPC max message size must be positive"))
	}

	// Validate logging config
	if cfg.Logging.LogPath == "" {
		errors = append(errors, fmt.Errorf("logging path must be specified"))
	}

	return errors
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

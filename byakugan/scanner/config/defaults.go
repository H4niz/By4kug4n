package config

import "time"

const (
	DefaultWorkerPoolSize   = 5
	DefaultQueueSize        = 100
	DefaultCPULimit         = 50  // 50% CPU
	DefaultMemoryLimit      = 512 // 512 MB
	DefaultConnectionLimit  = 200 // 200 concurrent connections
	DefaultBatchSize        = 50  // Default batch size for processing
	DefaultFlushInterval    = 5 * time.Second
	DefaultMaxMessageSize   = 10 * 1024 * 1024 // 10MB
	DefaultRetryAttempts    = 3
	DefaultRetryDelay       = time.Second
	DefaultBufferSize       = 4096
	DefaultMaxBatchSize     = 100
	DefaultFlushTimeout     = time.Second
	DefaultGRPCAddress      = "localhost:50051"
	DefaultKeepAliveTime    = 30 * time.Second
	DefaultKeepAliveTimeout = 10 * time.Second
	DefaultIdleTimeout      = 15 * time.Minute
	DefaultMaxAge           = 30 * time.Minute
	DefaultGracePeriod      = 5 * time.Second
	DefaultPoolMaxSize      = 50
	DefaultPoolMinIdle      = 5
	DefaultPoolMaxIdle      = 20
	DefaultPoolIdleTimeout  = 5 * time.Minute
)

// DefaultPerformanceConfig returns default performance configuration
func DefaultPerformanceConfig() *PerformanceConfig {
	return &PerformanceConfig{
		CPULimit:    DefaultCPULimit,
		MemoryLimit: DefaultMemoryLimit,
		Network: NetworkConfig{
			ConnectionsLimit: DefaultConnectionLimit,
		},
	}
}

// DefaultWorkerPoolConfig returns default worker pool configuration
func DefaultWorkerPoolConfig() WorkerPoolConfig {
	return WorkerPoolConfig{
		Size:                DefaultWorkerPoolSize,
		QueueSize:           DefaultQueueSize,
		BatchSize:           DefaultBatchSize,
		RetryCount:          DefaultRetryAttempts,
		RetryDelay:          DefaultRetryDelay,
		HealthCheckInterval: 30 * time.Second,
		MetricsEnabled:      true,
		ShutdownTimeout:     30 * time.Second,
	}
}

// DefaultHTTPClientConfig returns default HTTP client configuration
func DefaultHTTPClientConfig() HTTPClientConfig {
	return HTTPClientConfig{
		UserAgent:       "Byakugan Scanner/1.0",
		FollowRedirects: true,
		MaxRedirects:    10,
	}
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Nodes: []NodeConfiguration{
			{
				Node: NodeConfig{
					ID:           "scanner-01",
					Name:         "Primary Scanner",
					Region:       "us-east-1",
					Tags:         []string{"production"},
					Capabilities: []string{"rest", "graphql"},
				},
				Performance: PerformanceConfig{
					CPULimit:    80,
					MemoryLimit: 2048,
					Network: struct {
						BandwidthLimit   int `yaml:"bandwidth_limit"`
						ConnectionsLimit int `yaml:"connections_limit"`
					}{
						BandwidthLimit:   100,
						ConnectionsLimit: 500,
					},
				},
				WorkerPool: WorkerPoolConfig{
					Size:      10,
					QueueSize: 100,
				},
				HTTPClient: HTTPClientConfig{
					UserAgent:        "Byakugan Scanner/1.0",
					FollowRedirects:  true,
					MaxRedirects:     10,
					KeepAlive:        true,
					KeepAliveTimeout: 60,
					Timeout: struct {
						Connect time.Duration `yaml:"connect"`
						Read    time.Duration `yaml:"read"`
						Write   time.Duration `yaml:"write"`
					}{
						Connect: 5 * time.Second,
						Read:    30 * time.Second,
						Write:   5 * time.Second,
					},
					TLS: struct {
						VerifyCerts bool   `yaml:"verify_certs"`
						MinVersion  string `yaml:"min_version"`
					}{
						VerifyCerts: true,
						MinVersion:  "TLS1.2",
					},
					Compression: true,
				},
				RateLimiting: RateLimitingConfig{
					Enabled:        true,
					Strategy:       "adaptive",
					InitialRate:    50,
					MaxRate:        200,
					MinRate:        10,
					BackoffFactor:  1.5,
					RecoveryFactor: 1.2,
				},
			},
		},
		Logging: LoggingConfig{
			LogPath:        "byakugan/logs",
			Level:          "debug",
			IncludeTraceID: true,
			BufferSize:     DefaultBufferSize,
			AsyncWrite:     true,
			QueueSize:      1000,
			FlushInterval:  5 * time.Second,
			Rotation: LogRotationConfig{
				MaxSizeMB: 100,
				MaxFiles:  5,
				Compress:  true,
				MaxAge:    168 * time.Hour, // 7 days
			},
			Components: map[string]ComponentLog{
				"detector": {
					LogPath: "byakugan/logs/detector",
					Level:   "debug",
				},
				"payload": {
					LogPath: "byakugan/logs/payload",
					Level:   "debug",
				},
				"worker": {
					LogPath: "byakugan/logs/worker",
					Level:   "debug",
				},
			},
		},
		Proxy: ProxyConfig{
			Enabled:            false,
			Proxies:            []string{},
			CheckInterval:      5 * time.Minute,
			Timeout:            10 * time.Second,
			BlacklistThreshold: 3,
			BlacklistDuration:  30 * time.Minute,
			Rotation: struct {
				Strategy  string        `yaml:"strategy"`
				Interval  time.Duration `yaml:"interval"`
				PerWorker bool          `yaml:"per_worker"`
			}{
				Strategy:  "random",
				Interval:  5 * time.Minute,
				PerWorker: true,
			},
		},
		Retry: RetryConfig{
			MaxAttempts:       DefaultRetryAttempts,
			InitialDelay:      DefaultRetryDelay,
			MaxDelay:          10 * time.Second,
			BackoffMultiplier: 2.0,
			RetryOn: struct {
				StatusCodes   []int `yaml:"status_codes"`
				NetworkErrors bool  `yaml:"network_errors"`
				Timeouts      bool  `yaml:"timeouts"`
			}{
				StatusCodes:   []int{429, 500, 502, 503, 504},
				NetworkErrors: true,
				Timeouts:      true,
			},
		},
		Communication: CommunicationConfig{
			GRPC: GRPCConfig{
				Address:        DefaultGRPCAddress,
				MaxMessageSize: DefaultMaxMessageSize,
				KeepAlive: KeepAliveConfig{
					Enabled:             true,
					Time:                DefaultKeepAliveTime,
					Timeout:             DefaultKeepAliveTimeout,
					IdleTimeout:         DefaultIdleTimeout,
					MaxAge:              DefaultMaxAge,
					GracePeriod:         DefaultGracePeriod,
					PermitWithoutStream: true,
				},
				ConnectionPool: PoolConfig{
					MaxSize:     DefaultPoolMaxSize,
					MinIdle:     DefaultPoolMinIdle,
					MaxIdle:     DefaultPoolMaxIdle,
					IdleTimeout: DefaultPoolIdleTimeout,
				},
			},
		},
	}
}

package config

import (
	"fmt"
)

// ConfigValidator provides validation for configuration
type ConfigValidator struct {
	config *Config
}

// NewValidator creates a new configuration validator
func NewValidator(cfg *Config) *ConfigValidator {
	return &ConfigValidator{
		config: cfg,
	}
}

// Validate performs all validation checks
func (v *ConfigValidator) Validate() error {
	if v.config == nil {
		return fmt.Errorf("configuration is nil")
	}

	if err := v.validateBasic(); err != nil {
		return fmt.Errorf("basic validation failed: %w", err)
	}

	if errs := v.validateComponents(); len(errs) > 0 {
		return fmt.Errorf("component validation failed: %v", errs)
	}

	return nil
}

// validateBasic validates basic configuration requirements
func (v *ConfigValidator) validateBasic() error {
	if len(v.config.Nodes) == 0 {
		return fmt.Errorf("at least one node configuration is required")
	}

	if v.config.Communication.GRPC.Address == "" {
		return fmt.Errorf("GRPC address must be specified")
	}

	if v.config.Communication.GRPC.MaxMessageSize <= 0 {
		return fmt.Errorf("GRPC max message size must be positive")
	}

	return nil
}

// validateComponents validates all component configurations
func (v *ConfigValidator) validateComponents() []error {
	var errors []error

	// Validate GRPC configuration
	if err := v.validateGRPC(); err != nil {
		errors = append(errors, err)
	}

	// Validate nodes
	for i, node := range v.config.Nodes {
		if err := v.validateNode(i, node); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

// validateGRPC validates GRPC configuration
func (v *ConfigValidator) validateGRPC() error {
	grpc := v.config.Communication.GRPC

	if grpc.Address == "" {
		return fmt.Errorf("GRPC address must be specified")
	}

	if grpc.MaxMessageSize <= 0 {
		return fmt.Errorf("GRPC max message size must be positive")
	}

	// Validate KeepAlive settings
	if grpc.KeepAlive.Enabled {
		if grpc.KeepAlive.Time <= 0 {
			return fmt.Errorf("keepalive time must be positive")
		}
		if grpc.KeepAlive.Timeout <= 0 {
			return fmt.Errorf("keepalive timeout must be positive")
		}
		if grpc.KeepAlive.IdleTimeout <= 0 {
			return fmt.Errorf("idle timeout must be positive")
		}
		if grpc.KeepAlive.MaxAge <= 0 {
			return fmt.Errorf("max age must be positive")
		}
		if grpc.KeepAlive.GracePeriod <= 0 {
			return fmt.Errorf("grace period must be positive")
		}
	}

	// Validate connection pool
	if err := v.validateConnectionPool(grpc.ConnectionPool); err != nil {
		return fmt.Errorf("invalid connection pool config: %w", err)
	}

	return nil
}

func (v *ConfigValidator) validateConnectionPool(cfg PoolConfig) error {
	if cfg.MaxSize <= 0 {
		return fmt.Errorf("max size must be positive")
	}

	if cfg.MinIdle > cfg.MaxSize {
		return fmt.Errorf("min idle (%d) cannot exceed max size (%d)",
			cfg.MinIdle, cfg.MaxSize)
	}

	if cfg.MaxIdle > cfg.MaxSize {
		return fmt.Errorf("max idle (%d) cannot exceed max size (%d)",
			cfg.MaxIdle, cfg.MaxSize)
	}

	if cfg.IdleTimeout <= 0 {
		return fmt.Errorf("idle timeout must be positive")
	}

	return nil
}

// validateNode validates a single node configuration
func (v *ConfigValidator) validateNode(index int, node NodeConfiguration) error {
	if node.Node.ID == "" {
		return fmt.Errorf("node %d: ID is required", index)
	}

	if err := v.validatePerformance(index, node.Performance); err != nil {
		return err
	}

	if err := v.validateWorkerPool(index, node.WorkerPool); err != nil {
		return err
	}

	return nil
}

func (v *ConfigValidator) validatePerformance(index int, perf PerformanceConfig) error {
	if perf.CPULimit <= 0 || perf.CPULimit > 100 {
		return fmt.Errorf("node %d: CPU limit must be between 1-100", index)
	}

	if perf.MemoryLimit <= 0 {
		return fmt.Errorf("node %d: memory limit must be positive", index)
	}

	return nil
}

func (v *ConfigValidator) validateWorkerPool(index int, wp WorkerPoolConfig) error {
	if wp.Size <= 0 {
		return fmt.Errorf("node %d: worker pool size must be positive", index)
	}

	if wp.QueueSize <= 0 {
		return fmt.Errorf("node %d: queue size must be positive", index)
	}

	return nil
}

func (v *ConfigValidator) validateProxy() error {
	if len(v.config.Proxy.Proxies) == 0 {
		return fmt.Errorf("proxy list cannot be empty when proxy is enabled")
	}

	if v.config.Proxy.BlacklistThreshold <= 0 {
		return fmt.Errorf("proxy blacklist threshold must be positive")
	}

	return nil
}

func (v *ConfigValidator) validateLogging() error {
	if v.config.Logging.MaxSize <= 0 {
		return fmt.Errorf("log max size must be positive")
	}

	// Validate log levels for components
	for component, cfg := range v.config.Logging.Components {
		if cfg.LogPath == "" {
			return fmt.Errorf("log path for component %s must be specified", component)
		}
	}

	return nil
}

func (c *Config) Validate() error {
	if len(c.Nodes) == 0 {
		return fmt.Errorf("at least one node configuration required")
	}

	if err := c.Communication.GRPC.validate(); err != nil {
		return fmt.Errorf("invalid gRPC config: %w", err)
	}

	for i, node := range c.Nodes {
		if err := node.validate(); err != nil {
			return fmt.Errorf("invalid node config at index %d: %w", i, err)
		}
	}

	return nil
}

func (gc *GRPCConfig) validate() error {
	if gc.Address == "" {
		return fmt.Errorf("server address must be specified")
	}
	if gc.MaxMessageSize <= 0 {
		return fmt.Errorf("max message size must be positive")
	}
	return nil
}

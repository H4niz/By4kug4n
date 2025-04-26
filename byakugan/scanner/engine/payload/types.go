package payload

import (
	"sync"
	"sync/atomic"

	"github.com/haniz/byakugan/scanner/engine/logger"
)

// PayloadType represents the type of payload
type PayloadType string

const (
	TypeSQLi             PayloadType = "sqli"
	TypeXSS              PayloadType = "xss"
	TypePathTraversal    PayloadType = "path_traversal"
	TypeCommandInjection PayloadType = "cmd_injection"
	TypeJWT              PayloadType = "jwt"
	TypeUnknown          PayloadType = "unknown"
)

// VulnerabilityType represents the type of vulnerability
type VulnerabilityType string

// Payload represents a single attack payload
type Payload struct {
	ID          string                 `json:"id"`
	Type        PayloadType            `json:"type"`
	Value       string                 `json:"value"`
	Headers     map[string]string      `json:"headers"`
	Encoded     bool                   `json:"encoded"`
	Metadata    map[string]interface{} `json:"metadata"`
	Description string                 `json:"description"`
}

// PayloadTemplate represents a template for generating payloads
type PayloadTemplate struct {
	ID          string                 `yaml:"id"`
	Type        PayloadType            `yaml:"type"`
	Template    string                 `yaml:"template"`
	Templates   []string               `yaml:"templates,omitempty"` // Support multiple templates
	Pattern     string                 `yaml:"pattern,omitempty"`
	Variables   map[string][]string    `yaml:"variables"`
	Encodings   []string               `yaml:"encodings"`
	Tags        []string               `yaml:"tags"`
	Description string                 `yaml:"description"`
	Metadata    map[string]interface{} `yaml:"metadata"`
}

// Config represents payload generator configuration
type Config struct {
	Generator GeneratorConfig `json:"generator"`
	LogPath   string          `json:"log_path"` // Add LogPath field
}

// GeneratorConfig represents generator-specific configuration
type GeneratorConfig struct {
	TemplatePath           string              `json:"template_path"`
	CacheSize              int                 `json:"cache_size"`
	MaxPayloadSize         int                 `json:"max_payload_size"`
	TemplateTypes          []string            `json:"template_types"`
	Encodings              []string            `json:"encodings"`
	DefaultTransformations map[string][]string `json:"default_transformations"`
	Templating             struct {
		MaxRecursion    int      `json:"max_recursion"`
		EnableFunctions bool     `json:"enable_functions"`
		Functions       []string `json:"functions"`
	} `json:"templating"`
}

// Options contains configuration for the payload generator
type Options struct {
	TemplatePath string
	CacheSize    int
	MaxSize      int64
	LogPath      string          // Add log path
	LogLevel     logger.LogLevel // Add log level
}

// Cache implements a simple cache for payloads
type Cache struct {
	data   map[string][]Payload
	size   int
	mu     sync.RWMutex
	hits   int64
	misses int64
	logger logger.Logger
}

func NewCache(size int) *Cache {
	return &Cache{
		data:   make(map[string][]Payload),
		size:   size,
		hits:   0,
		misses: 0,
	}
}

// WithLogger adds a logger to the cache
func (c *Cache) WithLogger(l logger.Logger) *Cache {
	c.logger = l
	return c
}

// Get retrieves payloads from cache
func (c *Cache) Get(key string) ([]Payload, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	payloads, exists := c.data[key]
	if exists {
		atomic.AddInt64(&c.hits, 1)
		if c.logger != nil {
			c.logger.Debug("Cache hit", map[string]interface{}{
				"key":  key,
				"hits": c.hits,
			})
		}
	} else {
		atomic.AddInt64(&c.misses, 1)
		if c.logger != nil {
			c.logger.Debug("Cache miss", map[string]interface{}{
				"key":    key,
				"misses": c.misses,
			})
		}
	}
	return payloads, exists
}

// Set stores payloads in cache with eviction policy
func (c *Cache) Set(key string, payloads []Payload) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.data) >= c.size {
		// Evict random entry if cache is full
		for k := range c.data {
			if c.logger != nil {
				c.logger.Debug("Cache eviction", map[string]interface{}{
					"evicted_key": k,
					"cache_size":  len(c.data),
					"max_size":    c.size,
				})
			}
			delete(c.data, k)
			break
		}
	}
	c.data[key] = payloads
}

// Delete removes an entry from cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// Clear empties the cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string][]Payload)
}

// GetMetrics retrieves cache metrics
func (c *Cache) GetMetrics() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := float64(c.hits + c.misses)
	hitRate := 0.0
	if total > 0 {
		hitRate = float64(c.hits) / total
	}

	return map[string]interface{}{
		"size":     len(c.data),
		"max_size": c.size,
		"hits":     c.hits,
		"misses":   c.misses,
		"hit_rate": hitRate,
	}
}

// PayloadGenerator defines the interface for payload generation
type PayloadGenerator interface {
	LoadTemplates() error
	GetTemplate(id string) (*PayloadTemplate, bool)
	Generate(templateID string, data map[string]interface{}) ([]Payload, error)
}

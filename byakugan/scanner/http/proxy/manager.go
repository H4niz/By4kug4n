package proxy

import (
	"fmt"
	"math/rand"
	"net/url"
	"sync"
	"time"
)

// ProxyManager manages a list of proxies
type ProxyManager struct {
	proxies   []*ProxyEntry
	config    *Config
	blacklist map[string]time.Time
	mu        sync.RWMutex
	workers   map[int][]*ProxyEntry // Map worker IDs to dedicated proxies
}

// ProxyEntry represents a proxy with its metadata
type ProxyEntry struct {
	URL       *url.URL
	Failures  int
	LastUsed  time.Time
	LastCheck time.Time
}

// Config represents the configuration for the proxy manager
type Config struct {
	Enabled            bool          `yaml:"enabled"`
	Proxies            []string      `yaml:"proxies"`
	CheckInterval      time.Duration `yaml:"check_interval"`
	Timeout            time.Duration `yaml:"timeout"`
	BlacklistThreshold int           `yaml:"blacklist_threshold"`
	BlacklistDuration  time.Duration `yaml:"blacklist_duration"`
	Rotation           struct {
		Strategy  string        `yaml:"strategy"`
		Interval  time.Duration `yaml:"interval"`
		PerWorker bool          `yaml:"per_worker"`
	} `yaml:"rotation"`
}

// Remove duplicate Manager type and use ProxyManager consistently
type Manager = ProxyManager // Type alias for backward compatibility

// NewProxyManager creates a new proxy manager
func NewProxyManager(config *Config) (*ProxyManager, error) {
	if config == nil {
		return nil, fmt.Errorf("proxy config cannot be nil")
	}

	// Set default check interval if not specified
	if config.CheckInterval <= 0 {
		config.CheckInterval = 5 * time.Minute
	}

	pm := &ProxyManager{
		config:    config,
		blacklist: make(map[string]time.Time),
		workers:   make(map[int][]*ProxyEntry),
	}

	// Skip loading if disabled or no proxies configured
	if !config.Enabled || len(config.Proxies) == 0 {
		return pm, nil
	}

	// Load proxies from config
	validProxies := false
	for _, proxyStr := range config.Proxies {
		proxyURL, err := url.Parse(proxyStr)
		if err != nil || !isValidProxyURL(proxyURL) {
			continue
		}
		pm.AddProxy(proxyURL)
		validProxies = true
	}

	if config.Enabled && !validProxies {
		return nil, fmt.Errorf("no valid proxies configured")
	}

	// Start proxy health checker if enabled
	if config.Enabled {
		go pm.checkProxies()
	}

	return pm, nil
}

// Use NewProxyManager for both creation methods
func NewManager(config *Config) (*Manager, error) {
	return NewProxyManager(config)
}

// AddProxy adds a new proxy to the list
func (pm *ProxyManager) AddProxy(proxyURL *url.URL) {
	if proxyURL == nil {
		return
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Check if proxy already exists
	for _, existing := range pm.proxies {
		if existing.URL.String() == proxyURL.String() {
			return
		}
	}

	pm.proxies = append(pm.proxies, &ProxyEntry{
		URL:       proxyURL,
		LastCheck: time.Now(),
	})
}

// GetProxy returns the next available proxy
func (pm *ProxyManager) GetProxy() (*url.URL, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	for _, entry := range pm.proxies {
		if entry.Failures < 3 {
			entry.LastUsed = time.Now()
			return entry.URL, nil
		}
	}

	return nil, fmt.Errorf("no available proxies")
}

// ReportFailure reports a proxy failure
func (pm *ProxyManager) ReportFailure(proxyURL *url.URL) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for _, entry := range pm.proxies {
		if entry.URL.String() == proxyURL.String() {
			entry.Failures++
			if entry.Failures >= 3 {
				pm.blacklist[proxyURL.String()] = time.Now()
			}
			break
		}
	}
}

// checkProxies periodically checks proxy health
func (pm *ProxyManager) checkProxies() {
	// Ensure minimum check interval
	checkInterval := pm.config.CheckInterval
	if checkInterval < time.Minute {
		checkInterval = time.Minute
	}

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		pm.mu.Lock()
		now := time.Now()
		for i := len(pm.proxies) - 1; i >= 0; i-- {
			entry := pm.proxies[i]
			if entry.Failures >= pm.config.BlacklistThreshold {
				if blacklistTime, exists := pm.blacklist[entry.URL.String()]; exists {
					if now.Sub(blacklistTime) > pm.config.BlacklistDuration {
						entry.Failures = 0
						delete(pm.blacklist, entry.URL.String())
					}
				}
			}
		}
		pm.mu.Unlock()
	}
}

// AssignProxiesToWorker assigns dedicated proxies to a worker
func (pm *ProxyManager) AssignProxiesToWorker(workerID int) []*ProxyEntry {
	if !pm.config.Enabled {
		return nil
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Return existing assignment if already done
	if proxies, exists := pm.workers[workerID]; exists {
		return proxies
	}

	// Randomly assign ~30% of proxies to this worker
	workerProxies := make([]*ProxyEntry, 0)
	for _, proxy := range pm.proxies {
		if rand.Float32() < 0.3 {
			workerProxies = append(workerProxies, proxy)
		}
	}

	// Ensure at least one proxy is assigned
	if len(workerProxies) == 0 && len(pm.proxies) > 0 {
		workerProxies = append(workerProxies, pm.proxies[rand.Intn(len(pm.proxies))])
	}

	pm.workers[workerID] = workerProxies
	return workerProxies
}

func isValidProxyURL(u *url.URL) bool {
	if u == nil {
		return false
	}
	return u.Host != "" && (u.Scheme == "http" || u.Scheme == "https" || u.Scheme == "socks5")
}

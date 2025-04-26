package proxy

import (
	"fmt"
	"math/rand"
	"net/url"
	"sync"
	"time"
)

// Rotator represents a proxy rotation interface
type Rotator interface {
	GetNext() (*url.URL, error)
	ReportFailure(proxyURL *url.URL)
}

// ProxyRotator implements the Rotator interface
type ProxyRotator struct {
	manager    *ProxyManager
	workerID   int
	proxies    []*ProxyEntry
	current    *url.URL
	lastRotate time.Time
	mu         sync.RWMutex
}

// NewProxyRotator creates a new proxy rotator
func NewProxyRotator(manager *ProxyManager, workerID int) *ProxyRotator {
	if manager == nil {
		return nil
	}

	r := &ProxyRotator{
		manager:    manager,
		workerID:   workerID,
		lastRotate: time.Now(),
	}

	// Get dedicated proxies for this worker
	if proxies := manager.AssignProxiesToWorker(workerID); len(proxies) > 0 {
		r.proxies = proxies
	}

	return r
}

// GetNext returns the next proxy in rotation
func (r *ProxyRotator) GetNext() (*url.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.manager.config.Enabled || len(r.proxies) == 0 {
		return nil, nil
	}

	// Get available proxies (not blacklisted)
	available := make([]*ProxyEntry, 0)
	for _, proxy := range r.proxies {
		if proxy.Failures < r.manager.config.BlacklistThreshold {
			available = append(available, proxy)
		}
	}

	if len(available) == 0 {
		return nil, fmt.Errorf("no available proxies")
	}

	// Select proxy based on strategy
	var selected *ProxyEntry
	switch r.manager.config.Rotation.Strategy {
	case "random":
		selected = available[rand.Intn(len(available))]
	default: // round-robin
		for _, proxy := range available {
			if r.current == nil || proxy.URL.String() != r.current.String() {
				selected = proxy
				break
			}
		}
		if selected == nil {
			selected = available[0] // Fallback to first available
		}
	}

	r.current = selected.URL
	r.lastRotate = time.Now()
	selected.LastUsed = time.Now()

	return r.current, nil
}

// ReportFailure reports a proxy failure
func (r *ProxyRotator) ReportFailure(proxyURL *url.URL) {
	if proxyURL == nil {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Find and update the failed proxy
	for _, proxy := range r.proxies {
		if proxy.URL.String() == proxyURL.String() {
			proxy.Failures++
			if proxy.Failures >= r.manager.config.BlacklistThreshold {
				r.manager.blacklist[proxy.URL.String()] = time.Now()
			}
			// Force rotation on next GetNext() call
			if r.current != nil && r.current.String() == proxyURL.String() {
				r.current = nil
			}
			break
		}
	}
}

func (r *ProxyRotator) shouldRotate() bool {
	return r.current == nil || time.Since(r.lastRotate) > 5*time.Minute
}

// Proxy represents a proxy configuration
type Proxy struct {
	URL      string
	Protocol string
	Failures int
	LastUsed time.Time
}

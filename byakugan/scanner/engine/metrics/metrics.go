package metrics

import (
	"sync"
	"time"
)

type Metrics struct {
	mu        sync.RWMutex
	name      string
	counters  map[string]int64
	durations []time.Duration
	startTime time.Time
}

func NewMetrics(name string) *Metrics {
	return &Metrics{
		name:      name,
		counters:  make(map[string]int64),
		durations: make([]time.Duration, 0),
		startTime: time.Now(),
	}
}

func (m *Metrics) IncrementCounter(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters[name]++
}

func (m *Metrics) RecordDuration(d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.durations = append(m.durations, d)
}

func (m *Metrics) GetMetrics() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	metrics := make(map[string]interface{})

	// Basic metrics
	metrics["name"] = m.name
	metrics["uptime"] = time.Since(m.startTime)

	// Counter metrics
	for name, count := range m.counters {
		metrics[name] = count
	}

	// Duration metrics
	if len(m.durations) > 0 {
		var total time.Duration
		for _, d := range m.durations {
			total += d
		}
		metrics["avg_duration"] = total / time.Duration(len(m.durations))
		metrics["total_duration"] = total
	}

	return metrics
}

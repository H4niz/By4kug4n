package metrics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {
	m := NewMetrics("test_metrics")

	// Test counter
	m.IncrementCounter("test_counter")
	m.IncrementCounter("test_counter")

	// Test duration
	m.RecordDuration(100 * time.Millisecond)
	m.RecordDuration(200 * time.Millisecond)

	// Get metrics
	metrics := m.GetMetrics()

	// Verify metrics
	assert.Equal(t, "test_metrics", metrics["name"])
	assert.Equal(t, int64(2), metrics["test_counter"])
	assert.Equal(t, 150*time.Millisecond, metrics["avg_duration"])
	assert.Equal(t, 300*time.Millisecond, metrics["total_duration"])
}

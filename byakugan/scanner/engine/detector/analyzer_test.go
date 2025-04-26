package detector

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalyzer_CollectEvidence(t *testing.T) {
	config := DefaultConfig()
	analyzer := NewAnalyzer(config, nil)

	t.Run("collect evidence from request/response", func(t *testing.T) {
		// Create test server with delay
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Millisecond) // Force delay
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok"}`))
		}))
		defer server.Close()

		// Make actual request
		req, err := http.NewRequest("GET", server.URL, nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)

		// Collect evidence
		evidence, err := analyzer.CollectEvidence(req, resp, body)
		require.NoError(t, err)
		require.NotNil(t, evidence)

		// Verify response time is recorded
		assert.NotZero(t, evidence.ResponseTime)
		assert.Greater(t, evidence.ResponseTime, 10*time.Millisecond)
		assert.True(t, evidence.Validated)
	})

	t.Run("evidence collection with timing", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(20 * time.Millisecond) // Predictable delay
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok"}`))
		}))
		defer server.Close()

		req, err := http.NewRequest("GET", server.URL, nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)

		evidence, err := analyzer.CollectEvidence(req, resp, body)
		require.NoError(t, err)
		require.NotNil(t, evidence)

		// Verify timing accuracy
		assert.NotZero(t, evidence.ResponseTime)
		assert.GreaterOrEqual(t, evidence.ResponseTime, 20*time.Millisecond)
		assert.True(t, evidence.Validated)
	})
}

func TestAnalyzer_Stats(t *testing.T) {
	config := DefaultConfig()
	analyzer := NewAnalyzer(config, nil)

	t.Run("update stats", func(t *testing.T) {
		duration := 100 * time.Millisecond
		findings := 2

		// Initial values should be 0
		assert.Equal(t, int64(0), analyzer.stats.TotalScans)
		assert.Equal(t, int64(0), analyzer.stats.TotalFindings)

		// Update stats
		analyzer.updateStats(duration, findings)

		// Check updated values
		assert.Equal(t, int64(1), analyzer.stats.TotalScans)
		assert.Equal(t, int64(2), analyzer.stats.TotalFindings)
		assert.Equal(t, 1, len(analyzer.stats.scanDurations))
		assert.Equal(t, duration, analyzer.stats.scanDurations[0])
	})
}

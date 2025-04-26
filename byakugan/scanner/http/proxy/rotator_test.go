package proxy

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestProxyRotation(t *testing.T) {
	cfg := &Config{
		Enabled: true,
		Proxies: []string{
			"http://proxy1.test:8080",
			"http://proxy2.test:8080",
		},
		CheckInterval:      5 * time.Minute, // Set proper interval
		Timeout:            10 * time.Second,
		BlacklistThreshold: 3,
		BlacklistDuration:  30 * time.Minute,
		Rotation: struct {
			Strategy  string        `yaml:"strategy"`
			Interval  time.Duration `yaml:"interval"`
			PerWorker bool          `yaml:"per_worker"`
		}{
			Strategy:  "random",
			Interval:  time.Minute,
			PerWorker: true,
		},
	}

	manager, err := NewProxyManager(cfg)
	require.NoError(t, err)

	t.Run("Basic Rotation", func(t *testing.T) {
		rotator := NewProxyRotator(manager, 1)
		require.NotNil(t, rotator)

		// Test GetNext
		proxy1, err := rotator.GetNext()
		require.NoError(t, err)
		require.NotNil(t, proxy1)

		// Test rotation
		rotator.lastRotate = time.Now().Add(-2 * time.Minute)
		proxy2, err := rotator.GetNext()
		require.NoError(t, err)
		require.NotNil(t, proxy2)
	})

	t.Run("Failure Handling", func(t *testing.T) {
		rotator := NewProxyRotator(manager, 2)
		require.NotNil(t, rotator)

		// Get initial proxy
		proxy1, err := rotator.GetNext()
		require.NoError(t, err)
		require.NotNil(t, proxy1)

		// Report failure
		rotator.ReportFailure(proxy1)

		// Should get different proxy
		proxy2, err := rotator.GetNext()
		require.NoError(t, err)
		require.NotNil(t, proxy2)
		require.NotEqual(t, proxy1.String(), proxy2.String())
	})

	t.Run("Worker Specific Proxies", func(t *testing.T) {
		rotator1 := NewProxyRotator(manager, 3)
		rotator2 := NewProxyRotator(manager, 4)
		require.NotNil(t, rotator1)
		require.NotNil(t, rotator2)

		// Get proxies for each worker
		proxy1, err := rotator1.GetNext()
		require.NoError(t, err)
		require.NotNil(t, proxy1)

		proxy2, err := rotator2.GetNext()
		require.NoError(t, err)
		require.NotNil(t, proxy2)
	})

	t.Run("Disabled Proxies", func(t *testing.T) {
		disabledCfg := &Config{
			Enabled: false,
			Proxies: []string{
				"http://proxy1.test:8080",
			},
		}

		manager, err := NewProxyManager(disabledCfg)
		require.NoError(t, err)

		rotator := NewProxyRotator(manager, 5)
		require.NotNil(t, rotator)

		proxy, err := rotator.GetNext()
		require.NoError(t, err)
		require.Nil(t, proxy) // Should return nil when disabled
	})

	t.Run("Empty Proxy List", func(t *testing.T) {
		emptyCfg := &Config{
			Enabled: true,
			Proxies: []string{},
		}

		manager, err := NewProxyManager(emptyCfg)
		require.NoError(t, err)

		rotator := NewProxyRotator(manager, 6)
		require.NotNil(t, rotator)

		proxy, err := rotator.GetNext()
		require.NoError(t, err)
		require.Nil(t, proxy)
	})
}

func TestProxyRotatorNilManager(t *testing.T) {
	rotator := NewProxyRotator(nil, 1)
	require.Nil(t, rotator)
}

func TestProxyRotatorInvalidURLs(t *testing.T) {
	cfg := &Config{
		Enabled: true,
		Proxies: []string{
			"not-a-valid-url",
			"http://valid.test:8080",
			"ftp://invalid-scheme.test:8080",
		},
	}

	manager, err := NewProxyManager(cfg)
	require.NoError(t, err)
	require.NotNil(t, manager)

	rotator := NewProxyRotator(manager, 1)
	require.NotNil(t, rotator)

	// Should only get the valid proxy
	proxy, err := rotator.GetNext()
	require.NoError(t, err)
	require.NotNil(t, proxy)
	require.Equal(t, "valid.test:8080", proxy.Host)

	// Test failure handling with invalid URL
	rotator.ReportFailure(&url.URL{Host: "invalid"})
	proxy2, err := rotator.GetNext()
	require.NoError(t, err)
	require.NotNil(t, proxy2)
}

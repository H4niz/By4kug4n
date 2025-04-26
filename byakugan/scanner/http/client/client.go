package client

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/haniz/byakugan/scanner/config"
	"github.com/haniz/byakugan/scanner/http/proxy"
)

// TLSConfig represents TLS settings
type TLSConfig struct {
	MinVersion string `json:"min_version"`
	MaxVersion string `json:"max_version"`
}

// Client represents a custom HTTP client
type Client struct {
	client      *http.Client
	config      Config
	rateLimiter *RateLimiter
	retrier     *Retrier
	proxy       proxy.Rotator // Changed from *proxy.Rotator to proxy.Rotator
	transport   *http.Transport
}

// NewClient creates a new HTTP client
func NewClient(config Config) (*Client, error) {
	tlsConfig := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: !config.VerifyCerts,
	}

	transport := &http.Transport{
		TLSClientConfig:     tlsConfig,
		DisableCompression:  !config.Compression,
		DisableKeepAlives:   !config.KeepAlive,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   config.Timeout.Connect + config.Timeout.Read + config.Timeout.Write,
	}

	if !config.FollowRedirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return &Client{
		client: client,
		config: config,
		retrier: NewRetrier(RetryConfig{
			MaxAttempts:  3,
			InitialDelay: time.Second,
		}),
		rateLimiter: NewRateLimiter(RateLimitConfig{
			Enabled:       true,
			Strategy:      "fixed",
			InitialRate:   10,
			MaxRate:       20,
			MinRate:       5,
			BackoffFactor: 1.5,
		}),
	}, nil
}

// SetProxy sets the proxy rotator
func (c *Client) SetProxy(p proxy.Rotator) {
	c.proxy = p
}

// Do executes an HTTP request with retries and rate limiting
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c.rateLimiter != nil {
		c.rateLimiter.Wait()
	}

	if c.proxy != nil {
		proxyURL, err := c.proxy.GetNext()
		if err != nil {
			// Log warning but continue without proxy
			// TODO: Add proper logging
		} else if proxyURL != nil {
			req.URL.Host = proxyURL.Host
		}
	}

	req.Header.Set("User-Agent", c.config.UserAgent)

	return c.retrier.Do(func() (*http.Response, error) {
		resp, err := c.client.Do(req)
		if err != nil && c.proxy != nil {
			// Report proxy failure if request failed
			if proxyURL, _ := c.proxy.GetNext(); proxyURL != nil {
				c.proxy.ReportFailure(proxyURL)
			}
		}
		return resp, err
	})
}

// DoWithRetry executes an HTTP request with retries
func (c *Client) DoWithRetry(req *http.Request, maxAttempts int) (*http.Response, error) {
	return c.retrier.Do(func() (*http.Response, error) {
		if c.rateLimiter != nil {
			c.rateLimiter.Wait()
		}

		if c.proxy != nil {
			proxyURL, err := c.proxy.GetNext()
			if err != nil {
				// Log warning but continue without proxy
			} else if proxyURL != nil {
				req.URL.Host = proxyURL.Host
			}
		}

		resp, err := c.client.Do(req)
		if err != nil && c.proxy != nil {
			if proxyURL, _ := c.proxy.GetNext(); proxyURL != nil {
				c.proxy.ReportFailure(proxyURL)
			}
		}
		return resp, err
	})
}

// GetTimeout returns the client timeout configuration
func (c *Client) GetTimeout() *TimeoutConfig {
	return &c.config.Timeout
}

// GetRetryWait returns the retry wait duration
func (c *Client) GetRetryWait() time.Duration {
	return c.config.RetryWait
}

// GetUserAgent returns the user agent string
func (c *Client) GetUserAgent() string {
	return c.config.UserAgent
}

// GetRetryConfig returns the retry configuration matching config.RetryConfig format
func (c *Client) GetRetryConfig() config.RetryConfig {
	return config.RetryConfig{
		MaxAttempts:       c.config.Retry.MaxAttempts,
		InitialDelay:      c.config.Retry.InitialDelay,
		MaxDelay:          c.config.Retry.MaxDelay,
		BackoffMultiplier: c.config.Retry.BackoffMultiplier,
		RetryOn: struct {
			StatusCodes   []int `yaml:"status_codes"`
			NetworkErrors bool  `yaml:"network_errors"`
			Timeouts      bool  `yaml:"timeouts"`
		}{
			StatusCodes:   c.config.Retry.RetryOn.StatusCodes,
			NetworkErrors: c.config.Retry.RetryOn.NetworkErrors,
			Timeouts:      c.config.Retry.RetryOn.Timeouts,
		},
	}
}

// GetConfig returns the client configuration
func (c *Client) GetConfig() Config {
	return c.config
}

// GetTransport returns the HTTP transport
func (c *Client) GetTransport() http.RoundTripper {
	return c.getTransport()
}

func (c *Client) getTransport() *http.Transport {
	if c.transport == nil {
		c.transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !c.config.VerifyCerts,
			},
		}
	}
	return c.transport
}

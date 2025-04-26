package client

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net"
	"net/http"
	"time"
)

// Retrier handles retry logic for HTTP requests
type Retrier struct {
	config RetryConfig
}

// NewRetrier creates a new retrier with the given config
func NewRetrier(config RetryConfig) *Retrier {
	return &Retrier{
		config: config,
	}
}

// Do executes the given function with retry logic
func (r *Retrier) Do(fn func() (*http.Response, error)) (*http.Response, error) {
	var lastErr error
	attempts := 0

	for attempts < r.config.MaxAttempts {
		resp, err := fn()
		if err == nil {
			// Check if status code requires retry
			if !r.shouldRetryStatusCode(resp.StatusCode) {
				return resp, nil
			}
			// Close response body if we're going to retry
			resp.Body.Close()
		} else {
			// Check if error type requires retry
			if !r.shouldRetryError(err) {
				return nil, err
			}
			lastErr = err
		}

		attempts++
		if attempts == r.config.MaxAttempts {
			break
		}

		// Calculate backoff duration
		backoff := r.calculateBackoff(attempts)
		time.Sleep(backoff)
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("max retry attempts reached")
}

func (r *Retrier) shouldRetryStatusCode(code int) bool {
	for _, retryCode := range r.config.RetryOn.StatusCodes {
		if code == retryCode {
			return true
		}
	}
	return false
}

func (r *Retrier) shouldRetryError(err error) bool {
	if r.config.RetryOn.NetworkErrors {
		if _, ok := err.(net.Error); ok {
			return true
		}
	}

	if r.config.RetryOn.Timeouts {
		if err == context.DeadlineExceeded || err == context.Canceled {
			return true
		}
	}

	return false
}

func (r *Retrier) calculateBackoff(attempt int) time.Duration {
	// Calculate exponential backoff
	backoff := r.config.InitialDelay * time.Duration(math.Pow(r.config.BackoffMultiplier, float64(attempt-1)))

	// Cap at MaxDelay
	if backoff > r.config.MaxDelay {
		backoff = r.config.MaxDelay
	}

	// Add jitter (±20%)
	jitter := time.Duration(rand.Float64()*0.4-0.2) * backoff
	backoff += jitter

	return backoff
}

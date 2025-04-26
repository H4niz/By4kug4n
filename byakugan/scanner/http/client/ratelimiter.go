package client

import (
	"time"
)

// RateLimitConfig represents rate limiting configuration
type RateLimitConfig struct {
	Enabled       bool    `json:"enabled"`
	Strategy      string  `json:"strategy"`
	InitialRate   int     `json:"initial_rate"`
	MaxRate       int     `json:"max_rate"`
	MinRate       int     `json:"min_rate"`
	BackoffFactor float64 `json:"backoff_factor"`
}

// RateLimiter implements rate limiting for HTTP requests
type RateLimiter struct {
	config RateLimitConfig
	tokens chan struct{}
	stop   chan struct{}
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	r := &RateLimiter{
		config: config,
		tokens: make(chan struct{}, config.MaxRate),
		stop:   make(chan struct{}),
	}

	go r.refillTokens()
	return r
}

// Wait blocks until a token is available
func (rl *RateLimiter) Wait() {
	<-rl.tokens
}

// refillTokens continuously refills the token bucket
func (rl *RateLimiter) refillTokens() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for i := 0; i < rl.config.InitialRate; i++ {
				select {
				case rl.tokens <- struct{}{}:
				default:
					// Token bucket is full
					break
				}
			}
		case <-rl.stop:
			return
		}
	}
}

// UpdateRate updates the rate limit
func (rl *RateLimiter) UpdateRate(newRate int) {
	rl.config.InitialRate = newRate
}

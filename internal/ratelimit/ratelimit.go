// Package ratelimit provides a simple token-bucket rate limiter for
// throttling provider API calls during drift detection scans.
package ratelimit

import (
	"context"
	"fmt"
	"time"
)

// Limiter controls the rate of operations using a token-bucket approach.
type Limiter struct {
	tokens   chan struct{}
	rate     time.Duration
	quit     chan struct{}
}

// Config holds configuration for the rate limiter.
type Config struct {
	// RequestsPerSecond is the maximum number of requests allowed per second.
	RequestsPerSecond int
}

// New creates a new Limiter that allows up to cfg.RequestsPerSecond per second.
// Returns an error if RequestsPerSecond is less than 1.
func New(cfg Config) (*Limiter, error) {
	if cfg.RequestsPerSecond < 1 {
		return nil, fmt.Errorf("ratelimit: RequestsPerSecond must be >= 1, got %d", cfg.RequestsPerSecond)
	}

	rate := time.Second / time.Duration(cfg.RequestsPerSecond)
	l := &Limiter{
		tokens: make(chan struct{}, cfg.RequestsPerSecond),
		rate:   rate,
		quit:   make(chan struct{}),
	}

	// Pre-fill the bucket.
	for i := 0; i < cfg.RequestsPerSecond; i++ {
		l.tokens <- struct{}{}
	}

	// Refill tokens at the configured rate.
	go l.refill()

	return l, nil
}

// Wait blocks until a token is available or ctx is cancelled.
// Returns ctx.Err() if the context is done before a token is acquired.
func (l *Limiter) Wait(ctx context.Context) error {
	select {
	case <-l.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Stop shuts down the background refill goroutine.
func (l *Limiter) Stop() {
	close(l.quit)
}

func (l *Limiter) refill() {
	ticker := time.NewTicker(l.rate)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			select {
			case l.tokens <- struct{}{}:
			default:
				// Bucket is full; discard the token.
			}
		case <-l.quit:
			return
		}
	}
}

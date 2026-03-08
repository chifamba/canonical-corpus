package ratelimiter

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
)

// RateLimiter controls request rates using token bucket per host.
type RateLimiter struct {
	limiters   sync.Map
	maxPerHost float64
	semaphore  chan struct{}
}

// New creates a new RateLimiter.
func New(maxPerHost float64, globalConcurrency int) *RateLimiter {
	return &RateLimiter{
		maxPerHost: maxPerHost,
		semaphore:  make(chan struct{}, globalConcurrency),
	}
}

// getOrCreate returns the per-host limiter, creating one if needed.
func (r *RateLimiter) getOrCreate(host string) *rate.Limiter {
	val, loaded := r.limiters.Load(host)
	if loaded {
		return val.(*rate.Limiter)
	}
	limiter := rate.NewLimiter(rate.Limit(r.maxPerHost), int(r.maxPerHost)+1)
	actual, _ := r.limiters.LoadOrStore(host, limiter)
	return actual.(*rate.Limiter)
}

// Wait blocks until a token is available for the given host and a global slot is free.
func (r *RateLimiter) Wait(ctx context.Context, host string) error {
	// Acquire global concurrency slot.
	select {
	case r.semaphore <- struct{}{}:
	case <-ctx.Done():
		return ctx.Err()
	}
	// Wait for per-host token.
	limiter := r.getOrCreate(host)
	if err := limiter.Wait(ctx); err != nil {
		<-r.semaphore
		return err
	}
	return nil
}

// Release releases a global concurrency slot.
func (r *RateLimiter) Release() {
	<-r.semaphore
}

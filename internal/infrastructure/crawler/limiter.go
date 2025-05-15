package crawler

import (
	"context"
	"net/url"
	"sync"
	"time"
)

// RateLimiter implements rate limiting for crawling
type RateLimiter struct {
	delays map[string]time.Duration
	access map[string]time.Time
	mu     sync.Mutex
}

// NewRateLimiter creates a new RateLimiter
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		delays: make(map[string]time.Duration),
		access: make(map[string]time.Time),
	}
}

// SetDelay sets the delay for a domain
func (r *RateLimiter) SetDelay(domain string, delay time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.delays[domain] = delay
}

func (r *RateLimiter) Wait(ctx context.Context, urlStr string) error {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}
	domain := parsedURL.Hostname()

	r.mu.Lock()

	delay, exists := r.delays[domain]
	if !exists {
		delay = 1 * time.Second // Default delay if not set
	}

	lastAccess, exists := r.access[domain]

	now := time.Now()
	var waitTime time.Duration
	if exists {
		elapsed := now.Sub(lastAccess)
		if elapsed < delay {
			waitTime = delay - elapsed
		}
	}

	r.access[domain] = now.Add(waitTime)
	r.mu.Unlock()
	if waitTime > 0 {
		select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(waitTime):
				return nil
		}
	}

	return nil
}

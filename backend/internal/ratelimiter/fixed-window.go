package ratelimiter

import (
	"sync"
	"time"
)

// FixedWindowRateLimiter is a rate limiter that uses a fixed window algorithm.
type FixedWindowRateLimiter struct {
	sync.RWMutex
	clients map[string]int
	limit   int
	window  time.Duration
}

// NewFixedWindowLimiter creates a new FixedWindowRateLimiter with the given limit and window duration.
func NewFixedWindowLimiter(limit int, window time.Duration) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		clients: make(map[string]int),
		limit:   limit,
		window:  window,
	}
}

// Allow checks if a request from the given IP address is allowed based on the rate limit.
// It returns true if the request is allowed, and false otherwise.
// If the request is not allowed, it also returns the duration to wait before the next allowed request.
func (rl *FixedWindowRateLimiter) Allow(ip string) (bool, time.Duration) {
	rl.RLock()
	count, exists := rl.clients[ip]
	rl.RUnlock()

	if !exists || count < rl.limit {
		rl.Lock()
		if !exists {
			go rl.resetCount(ip)
		}

		rl.clients[ip]++
		rl.Unlock()
		return true, 0
	}

	return false, rl.window
}

// resetCount resets the request count for the given IP address after the window duration has passed.
func (rl *FixedWindowRateLimiter) resetCount(ip string) {
	time.Sleep(rl.window)
	rl.Lock()
	delete(rl.clients, ip)
	rl.Unlock()
}

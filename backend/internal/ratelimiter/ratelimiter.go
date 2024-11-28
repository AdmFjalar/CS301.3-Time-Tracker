package ratelimiter

import "time"

// Limiter is an interface for rate limiting.
type Limiter interface {
	// Allow checks if a request from the given IP address is allowed based on the rate limit.
	// It returns true if the request is allowed, and false otherwise.
	// If the request is not allowed, it also returns the duration to wait before the next allowed request.
	Allow(ip string) (bool, time.Duration)
}

// Config holds the configuration settings for the rate limiter.
type Config struct {
	RequestsPerTimeFrame int           // Number of requests allowed per time frame
	TimeFrame            time.Duration // Duration of the time frame
	Enabled              bool          // Flag to enable or disable the rate limiter
}

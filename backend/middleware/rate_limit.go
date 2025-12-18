package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a sliding window rate limiter
type RateLimiter struct {
	attempts map[string][]time.Time
	mu       sync.RWMutex
	limit    int           // Max requests allowed
	window   time.Duration // Time window for rate limiting
	stopCh   chan struct{} // Channel to signal cleanup goroutine to stop
}

// RateLimiterConfig holds configuration for rate limiter
type RateLimiterConfig struct {
	Limit  int           // Maximum number of requests
	Window time.Duration // Time window
}

// DefaultAuthRateLimiter returns a rate limiter configured for auth endpoints
// 5 attempts per minute (strict for login/register to prevent brute force)
func DefaultAuthRateLimiter() *RateLimiter {
	return NewRateLimiter(5, time.Minute)
}

// DefaultAPIRateLimiter returns a rate limiter for general API endpoints
// 100 requests per minute (more lenient for normal API usage)
func DefaultAPIRateLimiter() *RateLimiter {
	return NewRateLimiter(100, time.Minute)
}

// NewRateLimiter creates a new rate limiter with specified limit and window
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		attempts: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
		stopCh:   make(chan struct{}),
	}

	// Start background cleanup goroutine
	go rl.cleanup()

	return rl
}

// cleanup periodically removes expired entries to prevent memory leaks
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			now := time.Now()
			for ip, attempts := range rl.attempts {
				valid := make([]time.Time, 0)
				for _, t := range attempts {
					if now.Sub(t) < rl.window {
						valid = append(valid, t)
					}
				}
				if len(valid) == 0 {
					delete(rl.attempts, ip)
				} else {
					rl.attempts[ip] = valid
				}
			}
			rl.mu.Unlock()
		case <-rl.stopCh:
			return
		}
	}
}

// Stop gracefully stops the rate limiter cleanup goroutine
func (rl *RateLimiter) Stop() {
	close(rl.stopCh)
}

// isAllowed checks if a request from the given key is allowed
func (rl *RateLimiter) isAllowed(key string) (bool, int, time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Filter out expired attempts
	valid := make([]time.Time, 0)
	for _, t := range rl.attempts[key] {
		if t.After(windowStart) {
			valid = append(valid, t)
		}
	}
	rl.attempts[key] = valid

	remaining := rl.limit - len(valid)
	if remaining <= 0 {
		// Calculate retry time (when oldest request expires)
		if len(valid) > 0 {
			retryAfter := valid[0].Add(rl.window).Sub(now)
			return false, 0, retryAfter
		}
		return false, 0, rl.window
	}

	// Add current attempt
	rl.attempts[key] = append(rl.attempts[key], now)
	return true, remaining - 1, 0
}

// RateLimit returns a Gin middleware that rate limits requests by IP
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use client IP as the rate limit key
		key := c.ClientIP()

		allowed, remaining, retryAfter := rl.isAllowed(key)

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(rl.limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))

		if !allowed {
			c.Header("Retry-After", retryAfter.String())
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success":             false,
				"error":               "Too many requests. Please try again later.",
				"retry_after_seconds": int(retryAfter.Seconds()),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitByUserID returns a middleware that rate limits by user ID (for authenticated endpoints)
func (rl *RateLimiter) RateLimitByUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get user ID from context, fall back to IP
		key := c.ClientIP()
		if userID, exists := c.Get("user_id_string"); exists {
			key = userID.(string)
		}

		allowed, remaining, retryAfter := rl.isAllowed(key)

		c.Header("X-RateLimit-Limit", strconv.Itoa(rl.limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))

		if !allowed {
			c.Header("Retry-After", retryAfter.String())
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success":             false,
				"error":               "Too many requests. Please try again later.",
				"retry_after_seconds": int(retryAfter.Seconds()),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitByKey returns a middleware that rate limits by a custom key extracted from the request
func (rl *RateLimiter) RateLimitByKey(keyFunc func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := keyFunc(c)
		if key == "" {
			key = c.ClientIP()
		}

		allowed, remaining, retryAfter := rl.isAllowed(key)

		c.Header("X-RateLimit-Limit", strconv.Itoa(rl.limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))

		if !allowed {
			c.Header("Retry-After", retryAfter.String())
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success":             false,
				"error":               "Too many requests. Please try again later.",
				"retry_after_seconds": int(retryAfter.Seconds()),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

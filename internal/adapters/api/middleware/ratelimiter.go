package middleware

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"
)

// RateLimiter represents the rate limiting interface
type RateLimiter interface {
	Allow(key string) bool
	AllowWithContext(ctx context.Context, key string) bool
}

// TokenBucket represents a token bucket rate limiter
type TokenBucket struct {
	rate       float64
	bucketSize float64
	mu         sync.RWMutex
	tokens     map[string]*bucket
}

type bucket struct {
	tokens     float64
	lastAccess time.Time
}

// NewTokenBucket creates a new token bucket rate limiter
func NewTokenBucket(rate float64, bucketSize float64) *TokenBucket {
	return &TokenBucket{
		rate:       rate,
		bucketSize: bucketSize,
		tokens:     make(map[string]*bucket),
	}
}

// Allow checks if a request should be allowed
func (tb *TokenBucket) Allow(key string) bool {
	return tb.AllowWithContext(context.Background(), key)
}

// AllowWithContext checks if a request should be allowed with context
func (tb *TokenBucket) AllowWithContext(ctx context.Context, key string) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	b, exists := tb.tokens[key]

	if !exists {
		tb.tokens[key] = &bucket{
			tokens:     tb.bucketSize,
			lastAccess: now,
		}
		b = tb.tokens[key]
	}

	// Calculate tokens to add based on time passed
	timePassed := now.Sub(b.lastAccess).Seconds()
	tokensToAdd := timePassed * tb.rate

	b.tokens = min(b.tokens+tokensToAdd, tb.bucketSize)
	b.lastAccess = now

	if b.tokens >= 1 {
		b.tokens--
		return true
	}

	return false
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func RateLimitMiddleware(limiter RateLimiter) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if !limiter.Allow(ip) {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

package middleware

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/ouz/goauthboilerplate/internal/adapters/api/response"
	"github.com/ouz/goauthboilerplate/pkg/errors"
)

const (
	cleanupInterval = 5 * time.Minute
	bucketExpiry    = 1 * time.Hour
)

type RateLimiter interface {
	Allow(key string) bool
	AllowWithContext(ctx context.Context, key string) bool
	Cleanup()
}

type TokenBucket struct {
	rate       float64
	bucketSize float64
	mu         sync.RWMutex
	tokens     map[string]*bucket
	done       chan struct{}
}

type bucket struct {
	tokens     float64
	lastAccess time.Time
}

func NewTokenBucket(rate float64, bucketSize float64) *TokenBucket {
	tb := &TokenBucket{
		rate:       rate,
		bucketSize: bucketSize,
		tokens:     make(map[string]*bucket),
		done:       make(chan struct{}),
	}

	go tb.startCleanup()
	return tb
}

func (tb *TokenBucket) Allow(key string) bool {
	return tb.AllowWithContext(context.Background(), key)
}

func (tb *TokenBucket) AllowWithContext(ctx context.Context, key string) bool {
	select {
	case <-ctx.Done():
		return false
	default:
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
}

func (tb *TokenBucket) Cleanup() {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	for key, bucket := range tb.tokens {
		if now.Sub(bucket.lastAccess) > bucketExpiry {
			delete(tb.tokens, key)
		}
	}
}

func (tb *TokenBucket) startCleanup() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tb.Cleanup()
		case <-tb.done:
			return
		}
	}
}

func (tb *TokenBucket) Stop() {
	close(tb.done)
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
				response.Error(w, errors.InternalError("Failed to parse remote address", err))
				return
			}

			if !limiter.AllowWithContext(r.Context(), ip) {
				response.Error(w, errors.NewAppError(
					errors.ErrCodeTooManyRequests,
					"RATE_LIMIT_EXCEEDED",
					"Too many requests, please try again later",
					nil,
					http.StatusTooManyRequests,
				))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

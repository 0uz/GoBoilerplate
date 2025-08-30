package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ouz/goauthboilerplate/internal/config"
	"github.com/ouz/goauthboilerplate/internal/observability/metrics"
)

// HTTP metrics are now defined in the metrics package

func Logging(logger *config.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/v1/metrics" {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			metrics.HTTPRequestsInFlight.Inc()
			defer metrics.HTTPRequestsInFlight.Dec()

			wrapper := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(wrapper, r)
			duration := time.Since(start)

			// Sanitize endpoint for better grouping
			endpoint := sanitizeEndpoint(r.URL.Path)
			statusCode := strconv.Itoa(wrapper.status)

			// Prometheus metrics
			metrics.HTTPRequestsTotal.WithLabelValues(r.Method, endpoint, statusCode).Inc()
			metrics.HTTPRequestDuration.WithLabelValues(r.Method, endpoint).Observe(duration.Seconds())
			metrics.HTTPResponseSize.WithLabelValues(r.Method, endpoint).Observe(float64(wrapper.written))

			// Create log entry with additional fields
			entry := logger.WithFields(map[string]any{
				"method":     r.Method,
				"status":     wrapper.status,
				"path":       r.URL.EscapedPath(),
				"ip":         r.RemoteAddr,
				"user_agent": r.UserAgent(),
				"duration":   duration.String(),
				"size":       wrapper.written,
				"referer":    r.Referer(),
			})

			// Log based on status code
			switch {
			case wrapper.status >= 500:
				entry.Error("Server error occurred")
			case wrapper.status >= 400:
				entry.Warn("Client error occurred")
			case wrapper.status >= 300:
				entry.Info("Redirection occurred")
			default:
				entry.Info("Request completed successfully")
			}
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status  int
	written int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += int64(n)
	return n, err
}

// sanitizeEndpoint removes dynamic parts from URL path for better grouping
func sanitizeEndpoint(path string) string {
	if len(path) == 0 {
		return "/"
	}

	// Remove query parameters
	if idx := len(path); idx > 0 {
		for i, c := range path {
			if c == '?' {
				idx = i
				break
			}
		}
		path = path[:idx]
	}

	// Group similar endpoints
	switch {
	case path == "/":
		return "/"
	case path == "/live":
		return "/live"
	case path == "/ready":
		return "/ready"
	case path == "/metrics":
		return "/metrics"
	default:
		// For API endpoints, group by base path
		if len(path) > 8 && path[:8] == "/api/v1/" {
			rest := path[8:]
			if len(rest) > 0 {
				// Extract first segment after /api/v1/
				for i, c := range rest {
					if c == '/' && i > 0 {
						return "/api/v1/" + rest[:i] + "/*"
					}
				}
				return "/api/v1/" + rest
			}
		}
		return path
	}
}

package middleware

import (
	// "fmt"
	"net/http"
	"time"

	"github.com/ouz/goauthboilerplate/internal/config"
	// "github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/promauto"
)

// var (
// 	httpRequestsTotal = promauto.NewCounterVec(
// 		prometheus.CounterOpts{
// 			Name: "http_requests_total",
// 			Help: "Total number of HTTP requests",
// 		},
// 		[]string{"method", "path", "status", "error"},
// 	)

// 	httpRequestDuration = promauto.NewHistogramVec(
// 		prometheus.HistogramOpts{
// 			Name:    "http_request_duration_seconds",
// 			Help:    "Duration of HTTP requests",
// 			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
// 		},
// 		[]string{"method", "path", "status"},
// 	)

// 	httpResponseSize = promauto.NewHistogramVec(
// 		prometheus.HistogramOpts{
// 			Name:    "http_response_size_bytes",
// 			Help:    "Size of HTTP responses in bytes",
// 			Buckets: prometheus.ExponentialBuckets(100, 10, 8), // 100B to 10GB
// 		},
// 		[]string{"method", "path"},
// 	)

// 	activeConnections = promauto.NewGauge(
// 		prometheus.GaugeOpts{
// 			Name: "http_active_connections",
// 			Help: "Number of active HTTP connections",
// 		},
// 	)
// )

func Logging(logger *config.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			// activeConnections.Inc()
			// defer activeConnections.Dec()

			wrapper := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(wrapper, r)
			duration := time.Since(start)

			// Determine error status
			// isError := wrapper.status >= 400
			// errorLabel := fmt.Sprintf("%v", isError)

			// Prometheus metrics
			// httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, fmt.Sprintf("%d", wrapper.status), errorLabel).Inc()
			// httpRequestDuration.WithLabelValues(r.Method, r.URL.Path, fmt.Sprintf("%d", wrapper.status)).Observe(duration.Seconds())
			// httpResponseSize.WithLabelValues(r.Method, r.URL.Path).Observe(float64(wrapper.written))

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

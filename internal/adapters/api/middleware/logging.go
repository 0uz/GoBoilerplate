package middleware

import (
	"net/http"
	"time"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	activeConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_active_connections",
			Help: "Number of active HTTP connections",
		},
	)
)

func Logging(logger *logrus.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			activeConnections.Inc()
			defer activeConnections.Dec()

			wrapper := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(wrapper, r)
			duration := time.Since(start)

			// Prometheus metrics
			httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, fmt.Sprintf("%d", wrapper.status)).Inc()
			httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration.Seconds())

			// Log entry oluştur
			entry := logger.WithFields(logrus.Fields{
				"method":     r.Method,
				"status":     wrapper.status,
				"path":       r.URL.EscapedPath(),
				"ip":         r.RemoteAddr,
				"user_agent": r.UserAgent(),
				"duration":   duration,
			})

			switch {
			case wrapper.status >= 500:
				entry.Error("Request completed")
			case wrapper.status >= 400:
				entry.Warn("Request completed")
			default:
				entry.Info("Request completed")
			}
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

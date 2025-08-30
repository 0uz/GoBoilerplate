package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// HTTP metrics - middleware tarafından kullanılacak
var (
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status_code"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	HTTPResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Size of HTTP responses in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "endpoint"},
	)

	HTTPRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of HTTP requests currently being processed",
		},
	)
)

// Application business metrics
var (
	AuthAttemptsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_attempts_total",
			Help: "Total number of authentication attempts",
		},
		[]string{"type", "status"},
	)

	ActiveSessionsGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_sessions",
			Help: "Number of active user sessions",
		},
	)

	DatabaseConnectionsGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_connections_active",
			Help: "Number of active database connections",
		},
	)

	CacheOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_operations_total",
			Help: "Total number of cache operations",
		},
		[]string{"operation", "result", "cache_type"},
	)
)

// Business logic helper functions
func RecordAuthAttempt(authType, status string) {
	AuthAttemptsTotal.WithLabelValues(authType, status).Inc()
}

func UpdateActiveSessions(count float64) {
	ActiveSessionsGauge.Set(count)
}

func UpdateDatabaseConnections(count float64) {
	DatabaseConnectionsGauge.Set(count)
}

func RecordCacheOperation(operation, result, cacheType string) {
	CacheOperationsTotal.WithLabelValues(operation, result, cacheType).Inc()
}

// Cache-specific helper functions
func RecordCacheHit(cacheType string) {
	RecordCacheOperation("get", "hit", cacheType)
}

func RecordCacheMiss(cacheType string) {
	RecordCacheOperation("get", "miss", cacheType)
}

func RecordCacheSet(cacheType string) {
	RecordCacheOperation("set", "success", cacheType)
}

func RecordCacheDelete(cacheType string) {
	RecordCacheOperation("delete", "success", cacheType)
}

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

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
		[]string{"operation", "result", "cache_type", "cache_key"},
	)
)

func RecordAuthAttempt(authType, status string) {
	AuthAttemptsTotal.WithLabelValues(authType, status).Inc()
}

func UpdateActiveSessions(count float64) {
	ActiveSessionsGauge.Set(count)
}

func UpdateDatabaseConnections(count float64) {
	DatabaseConnectionsGauge.Set(count)
}

func RecordCacheOperation(operation, result, cacheType, cacheKey string) {
	CacheOperationsTotal.WithLabelValues(operation, result, cacheType, cacheKey).Inc()
}

func RecordCacheHit(cacheType, cacheKey string) {
	RecordCacheOperation("get", "hit", cacheType, cacheKey)
}

func RecordCacheMiss(cacheType, cacheKey string) {
	RecordCacheOperation("get", "miss", cacheType, cacheKey)
}

func RecordCacheSet(cacheType, cacheKey string) {
	RecordCacheOperation("set", "success", cacheType, cacheKey)
}

func RecordCacheDelete(cacheType, cacheKey string) {
	RecordCacheOperation("delete", "success", cacheType, cacheKey)
}

func RecordCacheHitLegacy(cacheType string) {
	RecordCacheOperation("get", "hit", cacheType, "unknown")
}

func RecordCacheMissLegacy(cacheType string) {
	RecordCacheOperation("get", "miss", cacheType, "unknown")
}

func RecordCacheSetLegacy(cacheType string) {
	RecordCacheOperation("set", "success", cacheType, "unknown")
}

func RecordCacheDeleteLegacy(cacheType string) {
	RecordCacheOperation("delete", "success", cacheType, "unknown")
}

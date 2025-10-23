package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total HTTP requests",
	}, []string{"method", "endpoint", "status"})

	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "HTTP request duration",
	}, []string{"method", "endpoint"})

	cacheHits = promauto.NewCounter(prometheus.CounterOpts{
		Name: "cache_hits_total",
		Help: "Total cache hits",
	})

	cacheMisses = promauto.NewCounter(prometheus.CounterOpts{
		Name: "cache_misses_total",
		Help: "Total cache misses",
	})
)

// RegisterMetrics registers /metrics endpoint
func RegisterMetrics() {
	http.Handle("/metrics", promhttp.Handler())
}

// RecordRequest records request metrics
func RecordRequest(method, endpoint string, status int, duration float64) {
	requestsTotal.WithLabelValues(method, endpoint, fmt.Sprintf("%d", status)).Inc()
	requestDuration.WithLabelValues(method, endpoint).Observe(duration)
}

// RecordCacheHit increments hit counter
func RecordCacheHit() {
	cacheHits.Inc()
}

// RecordCacheMiss increments miss counter
func RecordCacheMiss() {
	cacheMisses.Inc()
}

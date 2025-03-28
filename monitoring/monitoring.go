package monitoring

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	TrafficMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "route", "status_code"},
	)

	ActiveRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_active_requests",
			Help: "Number of active HTTP requests.",
		},
		[]string{"method", "route"},
	)

	LatencyMetric = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Latency of HTTP requests.",
		},
		[]string{"method", "route"},
	)
)

// IncrementTraffic increments the traffic metric
func IncrementTraffic(method, route string, status int) {
	TrafficMetric.WithLabelValues(method, route, strconv.Itoa(status)).Inc()
}

// IncrementActiveRequests increments the active requests metric
func IncrementActiveRequests(method, route string) {
	ActiveRequests.WithLabelValues(method, route).Inc()
}

// DecrementActiveRequests decrements the active requests metric
func DecrementActiveRequests(method, route string) {
	ActiveRequests.WithLabelValues(method, route).Dec()
}

// ObserveLatency observes the latency metric
func ObserveLatency(method, route string, latency float64) {
	LatencyMetric.WithLabelValues(method, route).Observe(latency)
}

// init registers all metrics in the default registerer
func init() {
	prometheus.MustRegister(LatencyMetric)
	prometheus.MustRegister(TrafficMetric)
	prometheus.MustRegister(ActiveRequests)
}

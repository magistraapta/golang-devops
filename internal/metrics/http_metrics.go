package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "golang-devops",
			Subsystem: "http",
			Name:      "http_requests_total",
		},
		[]string{"method", "endpoint", "status_code"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "golang-devops",
			Subsystem: "http",
			Name:      "http_request_duration_seconds",
			Help:      "Duration of HTTP requests in seconds",
		}, []string{"method", "endpoint", "status_code"},
	)
)

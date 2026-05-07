package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Kicks off any background collectors.
func RunMetrics(interval time.Duration) {
	RunCPUMetrics(interval)
}

// Exposes every collector to the registry main.go creates.
func RegisterAllMetrics(reg *prometheus.Registry) {
	reg.MustRegister(
		ProcessCPUUsagePercent,
		HTTPRequestsTotal,
		HTTPRequestDuration,
	)
}

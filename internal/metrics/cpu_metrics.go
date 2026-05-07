package metrics

import (
	"fmt"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/process"
)

var (
	ProcessCPUUsagePercent = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "golang_devops", // ✅ underscore, not hyphen
			Subsystem: "process",
			Name:      "cpu_usage_percent", // ✅ no duplicate "process"
			Help:      "Current CPU utilisation percentage of the process",
		},
	)
)

// RunCPUMetrics starts a background poller. Returns an error instead of panicking.
func RunCPUMetrics(interval time.Duration) error {
	proc, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return fmt.Errorf("metrics: create process handle: %w", err) // ✅ no panic
	}

	_, _ = proc.Percent(0) // ✅ prime the baseline so first real reading is accurate

	go func() {
		for {
			if pct, err := proc.Percent(interval); err == nil { // ✅ blocks for interval
				ProcessCPUUsagePercent.Set(pct)
			}
		}
	}()

	return nil
}

package balancer

import (
	"Go/config"
	"time"
)

func AnalyzeSystemState(target string) string {
	backends := config.MetricsMap
	healthyCount := 0
	highLatencyCount := 0
	totalBackends := len(backends)

	for _, backend := range backends {
		if backend.Metrics.IsHealthy && backend.Metrics.ConsecutiveFails <= 5 {
			healthyCount++
		}
		if backend.Metrics.LastLatency >= 500*time.Millisecond || backend.Metrics.AvgLatency >= 500*time.Millisecond {
			highLatencyCount++
		}
	}

	if healthyCount == 0 {
		return "AllFailed"
	}
	if float64(healthyCount)/float64(totalBackends) < 0.5 {
		return "ManyFailed"
	}
	if highLatencyCount > 1 {
		return "HighLatency"
	}

	return "Stable"
}

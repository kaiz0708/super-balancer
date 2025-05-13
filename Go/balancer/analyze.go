package balancer

import (
	"Go/config"
)

func AnalyzeSystemState() string {
	backends := config.MetricsMap
	healthyCount := 0
	highLatencyCount := 0
	totalBackends := len(backends)

	for _, backend := range backends {
		if backend.Metrics.IsHealthy && backend.Metrics.ConsecutiveFails <= config.ConfigSystem.ConsecutiveFails {
			healthyCount++
		}
		if backend.Metrics.TimeoutBreak >= config.ConfigSystem.TimeOutRate {
			highLatencyCount++
		}
	}

	if healthyCount == 0 {
		return "AllFailed"
	}
	if float64(healthyCount)/float64(totalBackends) <= 0.5 {
		return "ManyFailed"
	}
	if highLatencyCount >= 1 {
		return "HighLatency"
	}

	return "Stable"
}

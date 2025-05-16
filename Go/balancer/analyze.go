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
		config.StateSystem = config.AllFailed
		return config.AllFailed
	}
	if healthyCount <= totalBackends/2 {
		config.StateSystem = config.ManyFailed
		return config.ManyFailed
	}
	if highLatencyCount >= totalBackends/2 {
		config.StateSystem = config.HighLatency
		return config.HighLatency
	}
	config.StateSystem = config.Stable

	return config.Stable
}

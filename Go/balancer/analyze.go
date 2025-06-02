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
		if backend.Metrics.TimeoutRate >= config.ConfigSystem.TimeOutRate {
			highLatencyCount++
		}
	}

	if healthyCount == 0 {
		config.ConfigDefaultSystem.StateSystem = config.AllFailed
		return config.AllFailed
	}
	if healthyCount <= totalBackends/2 {
		config.ConfigDefaultSystem.StateSystem = config.ManyFailed
		return config.ManyFailed
	}
	if highLatencyCount >= totalBackends/2 {
		config.ConfigDefaultSystem.StateSystem = config.HighLatency
		return config.HighLatency
	}
	config.ConfigDefaultSystem.StateSystem = config.Stable
	config.ConfigSystem.Algorithm = config.ConfigDefaultSystem.AlgorithmUserChoose

	return config.Stable
}

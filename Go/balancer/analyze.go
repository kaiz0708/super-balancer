package balancer

import (
	"Go/config"
	"Go/factory"
	"net/http"
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
	if healthyCount == totalBackends {
		config.ConfigDefaultSystem.StateSystem = config.Stable
	}

	return config.Stable
}

func SelectServer(state string, r *http.Request) string {
	selected := ""
	if !config.ConfigSystem.SmartMode {
		selected = factory.Factory.GetStrategy(config.ConfigSystem.Algorithm).SelectServer(r)
		return selected
	}
	switch state {
	case config.ManyFailed:
		config.ConfigSystem.Algorithm = config.WeightedSuccessRateBalancer
		selected = factory.Factory.GetStrategy(config.WeightedSuccessRateBalancer).SelectServer(r)
	case config.HighLatency:
		selected = factory.Factory.GetStrategy(config.LowLatencyWeightedBalancer).SelectServer(r)
	default:
		config.ConfigSystem.Algorithm = config.ConfigDefaultSystem.AlgorithmUserChoose
		selected = factory.Factory.GetStrategy(config.ConfigDefaultSystem.AlgorithmUserChoose).SelectServer(r)
	}
	return selected
}

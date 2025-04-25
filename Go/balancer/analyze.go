package balancer

import (
	"Go/config"
	"time"
)

func CheckHealthyBackend(backend string) bool {
	target := config.MetricsMap[backend]

	failRate := float64(target.Metrics.FailureCount) / float64(target.Metrics.RequestCount+1)

	return target.Metrics.ConsecutiveFails >= 5 ||
		failRate > 0.5 ||
		target.Metrics.AvgLatency > 500*time.Millisecond
}

func AnalyzeSystemState(target string) string {
	checkHealthy := CheckHealthyBackend(target)
	config.MetricsMap[target].Mutex.Lock()
	defer config.MetricsMap[target].Mutex.Unlock()
	if checkHealthy {
		config.MetricsMap[target].Metrics.IsHealthy = false
	} else {
		config.MetricsMap[target].Metrics.IsHealthy = true
	}

	backends := config.MetricsMap
	healthyCount := 0
	highLatencyCount := 0
	totalBackends := len(backends)

	for _, backend := range backends {
		if backend.Metrics.IsHealthy && backend.Metrics.ConsecutiveFails < 3 {
			healthyCount++
		}
		if backend.Metrics.LastLatency > 500*time.Millisecond || backend.Metrics.AvgLatency > 500*time.Millisecond {
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

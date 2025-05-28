package algo_default

import (
	"Go/config"
	"math"
)

func WeightedRoundRobin() string {
	metrics := config.MetricsMap
	sumWeight := SumWeightMetrics()
	selected := ""
	maxCurrentWeight := int64(math.MinInt64)

	for backend, m := range metrics {
		target := m.Metrics
		if !target.IsHealthy {
			continue
		}
		currentWeight := target.CurrentWeight + target.Weight

		if maxCurrentWeight <= currentWeight {
			maxCurrentWeight = currentWeight
			selected = backend
		}
	}

	metrics[selected].Mutex.Lock()
	if selected != "" {
		config.MetricsMap[selected].Metrics.CurrentWeight += config.MetricsMap[selected].Metrics.Weight
		config.MetricsMap[selected].Metrics.CurrentWeight -= int64(sumWeight)
	}
	metrics[selected].Mutex.Unlock()

	return selected
}

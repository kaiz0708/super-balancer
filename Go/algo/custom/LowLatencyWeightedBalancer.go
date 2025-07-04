package custom

import (
	"Go/config"
	"Go/utils"
	"math"
	"net/http"
	"time"
)

type LowLatencyWeightedBalancerStrategy struct{}

func (r *LowLatencyWeightedBalancerStrategy) SelectServer(t *http.Request) string {
	metrics := config.MetricsMap
	selected := ""
	minLatency := time.Duration(math.MaxInt64)
	sumWeight := utils.SumWeightMetrics()
	maxCurrentWeight := int64(math.MinInt64)

	for backend, target := range metrics {
		if !target.Metrics.IsHealthy {
			continue
		}

		currentWeight := target.Metrics.CurrentWeight + target.Metrics.Weight

		if target.Metrics.AvgLatency < minLatency || maxCurrentWeight <= currentWeight {
			maxCurrentWeight = currentWeight
			minLatency = target.Metrics.AvgLatency
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

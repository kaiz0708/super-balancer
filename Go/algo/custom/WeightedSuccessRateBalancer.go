package custom

import (
	"Go/config"
	"math"
)

func WeightedSuccessRateBalancer() string {
	metrics := config.MetricsMap
	selected := ""
	highestScore := -1.0
	lowestActiveConn := int64(math.MaxInt64)
	sumWeight := int64(0)

	for _, m := range metrics {
		if m.Metrics.IsHealthy {
			sumWeight += m.Metrics.Weight
		}
	}

	for backend, m := range metrics {
		if !m.Metrics.IsHealthy {
			continue
		}

		totalRequest := m.Metrics.RequestCount
		successRequest := m.Metrics.SuccessCount
		activeConn := m.Metrics.ActiveConnections
		weight := m.Metrics.Weight

		if totalRequest == 0 {
			continue
		}

		successRate := float64(successRequest) / float64(totalRequest)

		normalizedWeight := float64(weight) / float64(sumWeight)

		score := successRate * normalizedWeight

		if (score > highestScore) || (score == highestScore && activeConn < lowestActiveConn) {
			highestScore = score
			lowestActiveConn = activeConn
			selected = backend
		}
	}

	if selected != "" {
		metrics[selected].Mutex.Lock()
		config.MetricsMap[selected].Metrics.CurrentWeight += config.MetricsMap[selected].Metrics.Weight
		config.MetricsMap[selected].Metrics.CurrentWeight -= sumWeight
		metrics[selected].Mutex.Unlock()
	}

	return selected
}

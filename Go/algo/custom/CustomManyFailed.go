package custom

import (
	"Go/config"
	"math"
)

func CustomManyFailed() string {
	metrics := config.MetricsMap
	selected := ""
	highestSuccessRate := -1.0
	lowestActiveConn := int64(math.MaxInt64)

	for backend, m := range metrics {
		if !m.Metrics.IsHealthy {
			continue
		}

		totalRequest := m.Metrics.RequestCount
		successRequest := m.Metrics.SuccessCount
		activeConn := m.Metrics.ActiveConnections

		if totalRequest == 0 {
			continue
		}

		successRate := float64(successRequest) / float64(totalRequest)

		if (successRate > highestSuccessRate) || (successRate == highestSuccessRate && activeConn < lowestActiveConn) {
			highestSuccessRate = successRate
			lowestActiveConn = activeConn
			selected = backend
		}
	}

	return selected
}

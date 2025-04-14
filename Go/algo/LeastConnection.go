package algo

import (
	"Go/config"
	"math"
)

func LeastConnection() string {
	metrics := config.MetricsMap
	selected := ""
	minConnection := int64(math.MaxInt64)

	for backend, m := range metrics {
		target := m.Metrics

		if !target.IsHealthy {
			continue
		}

		if target.ActiveConnections <= minConnection && target.IsHealthy {
			minConnection = target.ActiveConnections
			selected = backend
		}
	}
	return selected
}

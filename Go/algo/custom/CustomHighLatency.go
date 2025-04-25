package custom

import (
	"Go/config"
	"math"
	"time"
)

func CustomHighLatency() string {
	metrics := config.MetricsMap
	selected := ""
	minLatency := time.Duration(math.MaxInt64)

	for backend, target := range metrics {
		if !target.Metrics.IsHealthy {
			continue
		}

		if target.Metrics.AvgLatency < minLatency {
			minLatency = target.Metrics.AvgLatency
			selected = backend
		}

	}
	return selected
}

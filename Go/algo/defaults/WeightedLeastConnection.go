package algo_default

import (
	"Go/config"
	"math"
	"net/http"
)

type WeightedLeastConnectionStrategy struct{}

func (r *WeightedLeastConnectionStrategy) SelectServer(t *http.Request) string {
	metrics := config.MetricsMap
	maxLoadBackend := math.MaxFloat64
	selected := ""

	for backend, m := range metrics {
		target := m.Metrics

		if !target.IsHealthy {
			continue
		}

		load := float64(target.ActiveConnections) / float64(m.Metrics.Weight)

		if load <= maxLoadBackend {
			maxLoadBackend = load
			selected = backend
		}
	}

	return selected
}

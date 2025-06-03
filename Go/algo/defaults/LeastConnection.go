package algo_default

import (
	"Go/config"
	"math"
	"net/http"
)

type LeastConnectionStrategy struct{}

func (r *LeastConnectionStrategy) SelectServer(t *http.Request) string {
	metrics := config.MetricsMap
	selected := ""
	minConnection := int64(math.MaxInt64)

	for backend, m := range metrics {
		target := m.Metrics

		if !target.IsHealthy {
			continue
		}

		if target.ActiveConnections <= minConnection {
			minConnection = target.ActiveConnections
			selected = backend
		}
	}
	return selected
}

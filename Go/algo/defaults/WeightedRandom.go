package algo_default

import (
	"Go/config"
	"Go/utils"
	"math/rand"
	"net/http"
	"time"
)

type WeightedRandomStrategy struct{}

func (r *WeightedRandomStrategy) SelectServer(t *http.Request) string {
	backends := config.MetricsMap
	totalWeight := utils.SumWeightMetrics()
	selected := ""
	rand.Seed(time.Now().UnixNano())
	value := rand.Int63n(totalWeight)

	for backend, m := range backends {
		if !m.Metrics.IsHealthy {
			continue
		}

		value -= m.Metrics.Weight
		if value < 0 {
			selected = backend
			break
		}
	}
	return selected
}

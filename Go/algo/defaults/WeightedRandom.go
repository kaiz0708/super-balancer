package algo_default

import (
	"Go/config"
	"Go/utils"
	"math/rand"
	"time"
)

func WeightedRandom() string {

	backends := config.MetricsMap
	totalWeight := utils.SumWeightMetrics()
	selected := ""
	rand.Seed(time.Now().UnixNano())
	r := rand.Int63n(totalWeight)

	for backend, m := range backends {
		if !m.Metrics.IsHealthy {
			continue
		}

		r -= m.Metrics.Weight
		if r < 0 {
			selected = backend
			break
		}
	}
	return selected
}

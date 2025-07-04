package algo_default

import (
	"Go/config"
	"Go/utils"
	"hash/fnv"
	"net/http"
)

type IpHashStrategy struct{}

func (r *IpHashStrategy) SelectServer(t *http.Request) string {
	ip := utils.GetClientIP(t)
	backends := config.ConfigSystem.Servers
	metrics := config.MetricsMap
	selected := ""

	healthyBackends := make([]config.BackendConfig, 0)
	for _, backend := range backends {
		if metrics[backend.UrlConfig].Metrics.IsHealthy {
			healthyBackends = append(healthyBackends, backend)
		}
	}

	hasher := fnv.New32a()
	hasher.Write([]byte(ip))
	hashValue := hasher.Sum32()

	totalWeight := utils.SumWeightMetrics()

	hashValue = hashValue % uint32(totalWeight)
	currentWeight := int64(0)

	for _, backend := range healthyBackends {
		currentWeight += backend.WeightConfig
		if uint32(currentWeight) > hashValue {
			selected = backend.UrlConfig
			break
		}
	}
	return selected
}

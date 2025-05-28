package algo_default

import (
	"Go/config"
	"hash/fnv"
	"net/http"
)

func IpHash(r *http.Request) string {
	ip := getClientIP(r)
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

	totalWeight := SumWeightMetrics()

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

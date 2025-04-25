package algo_default

import (
	"Go/config"
	"math/rand"
)

func Random() string {
	backend := config.BackendServers
	metrics := config.MetricsMap
	selected := ""

	index := rand.Intn(len(backend))

	for !metrics[backend[index].UrlConfig].Metrics.IsHealthy {
		index = rand.Intn(len(backend))
	}

	selected = backend[index].UrlConfig

	return selected
}

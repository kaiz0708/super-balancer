package algo_default

import (
	"Go/config"
	"sync/atomic"
)

var countRequest uint64

func RoundRobin() string {
	metrics := config.MetricsMap
	backend := config.BackendServers
	n := len(backend)

	if n == 0 {
		return ""
	}

	index := atomic.AddUint64(&countRequest, 1)

	for !metrics[backend[int(index)%n].UrlConfig].Metrics.IsHealthy {
		index = atomic.AddUint64(&countRequest, 1)
	}

	selected := backend[int(index)%n].UrlConfig
	return selected
}

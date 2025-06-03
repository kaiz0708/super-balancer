package algo_default

import (
	"Go/config"
	"net/http"
	"sync/atomic"
)

var countRequest uint64

type RoundRobinStrategy struct{}

func (r *RoundRobinStrategy) SelectServer(t *http.Request) string {
	metrics := config.MetricsMap
	backend := config.ConfigSystem.Servers
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

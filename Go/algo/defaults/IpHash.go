package algo_default

import (
	"Go/config"
	"hash/fnv"
	"net"
	"net/http"
	"strings"
)

func getClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

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

	totalWeight := int64(0)
	for _, backend := range healthyBackends {
		totalWeight += backend.WeightConfig
	}

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

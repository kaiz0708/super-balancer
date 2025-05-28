package algo_default

import (
	"Go/config"
	"net"
	"net/http"
	"strings"
)

func SumWeightMetrics() int64 {
	metrics := config.MetricsMap
	sumWeight := 0

	for _, m := range metrics {
		target := m.Metrics

		if target.IsHealthy {
			sumWeight += int(target.Weight)
		}
	}
	return int64(sumWeight)
}

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

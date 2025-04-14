package config

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Metrics struct {
	RequestCount      uint64
	SuccessCount      uint64
	FailureCount      uint64
	TotalLatency      time.Duration
	LastLatency       time.Duration
	AvgLatency        time.Duration
	LastChecked       time.Time
	ConsecutiveFails  uint64
	IsHealthy         bool
	LastStatus        int
	ActiveConnections int64
	Weight            int64
	CurrentWeight     int64
}

type BackendConfig struct {
	UrlConfig    string
	WeightConfig int64
}

type BackendMetrics struct {
	Mutex   sync.Mutex
	Metrics *Metrics
}

var MetricsMap = map[string]*BackendMetrics{}

var LoadBalancerDefault string

var BackendServers = []BackendConfig{
	{UrlConfig: "http://localhost:3001/", WeightConfig: 1},
	{UrlConfig: "http://localhost:3002/", WeightConfig: 5},
}

func InitServer() {
	LoadBalancerDefault = os.Getenv("ROUND_ROBIN")
	urls := BackendServers
	fmt.Println("Default : ", LoadBalancerDefault)
	for _, url := range urls {
		MetricsMap[url.UrlConfig] = &BackendMetrics{
			Metrics: &Metrics{
				IsHealthy: true,
				Weight:    url.WeightConfig,
			},
		}
	}
}

package config

import (
	"sync"
	"time"
)

const (
	RoundRobinAlgo          = "ROUND_ROBIN"
	LeastConnectionAlgo     = "LEAST_CONNECTION"
	WeightedLeastConnection = "WEIGHTED_LEAST_CONNECTION"
	WeightedRoundRobin      = "WEIGHTED_ROUND_ROBIN"
	RandomAlgo              = "RANDOM"
)

type Metrics struct {
	RequestCount       uint64
	SuccessCount       uint64
	FailureCount       uint64
	TotalLatency       time.Duration
	LastLatency        time.Duration
	AvgLatency         time.Duration
	LastChecked        time.Time
	ConsecutiveFails   uint64
	ConsecutiveSuccess uint64
	TimeoutBreak       uint64
	IsHealthy          bool
	LastStatus         int
	ActiveConnections  int64
	Weight             int64
	CurrentWeight      int64
}

type BackendConfig struct {
	UrlConfig        string `json:"url" validate:"required"`
	WeightConfig     int64  `json:"weight"`
	HealthPathConfig string `json:"healthPath"`
}

type BackendMetrics struct {
	Mutex      sync.Mutex
	Metrics    *Metrics
	HealthPath string
}

type Config struct {
	DefaultProxy string          `json:"defaultProxy" validate:"required"`
	Algorithm    string          `json:"algorithm" validate:"required"`
	Servers      []BackendConfig `json:"servers" validate:"required,dive"`
}

var MetricsMap = map[string]*BackendMetrics{}

var LoadBalancerDefault string

var ConsecutiveFails uint64

var ConsecutiveSuccess uint64

var FailRate float64

var TimeOutRate uint64

var BackendServers = []BackendConfig{}

func InitServer() {
	urls := BackendServers
	for _, url := range urls {
		MetricsMap[url.UrlConfig] = &BackendMetrics{
			Metrics: &Metrics{
				IsHealthy: true,
				Weight:    url.WeightConfig,
			},
			HealthPath: url.HealthPathConfig,
		}
	}
}

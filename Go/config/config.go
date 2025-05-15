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
	WeightedRandom          = "WEIGHTED_RANDOM"
	Unhealthy               = "Unhealthy"
	Recovery                = "Recovery"
	Healthy                 = "Healthy"
	AllFailed               = "AllFailed"
	ManyFailed              = "ManyFailed"
	HighLatency             = "HighLatency"
	Stable                  = "Stable"
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

type AuthBasic struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
	DefaultProxy       string
	Algorithm          string
	Servers            []BackendConfig
	ConsecutiveFails   uint64
	ConsecutiveSuccess uint64
	FailRate           float64
	TimeOutRate        uint64
	TimeOutDelay       uint64
	ActiveLogin        bool
	AuthConfig         AuthBasic
}

var MetricsMap = map[string]*BackendMetrics{}

var ConfigSystem Config

func InitServer() {
	urls := ConfigSystem.Servers
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

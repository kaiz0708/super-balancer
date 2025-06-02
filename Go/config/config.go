package config

import (
	"sync"
	"time"
)

const (
	RoundRobinAlgo              = "ROUND_ROBIN"
	LeastConnectionAlgo         = "LEAST_CONNECTION"
	WeightedLeastConnection     = "WEIGHTED_LEAST_CONNECTION"
	WeightedRoundRobin          = "WEIGHTED_ROUND_ROBIN"
	RandomAlgo                  = "RANDOM"
	WeightedRandom              = "WEIGHTED_RANDOM"
	WeightedSuccessRateBalancer = "WEIGHTED_SUCCESS_RATE_BALANCER"
	LowLatencyWeightedBalancer  = "LOW_LATENCY_WEIGHTED_BALANCER"
	IpHash                      = "IP_HASH"
	Unhealthy                   = "Unhealthy"
	Recovery                    = "Recovery"
	Healthy                     = "Healthy"
	AllFailed                   = "AllFailed"
	ManyFailed                  = "ManyFailed"
	HighLatency                 = "HighLatency"
	Stable                      = "Stable"
	ErrorConnectionRefused      = "connection refused"
	ErrorNoSuchHost             = "no such host"
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
	TimeoutRate        uint64
	IsHealthy          bool
	LastStatus         int
	ActiveConnections  int64
	Weight             int64
	CurrentWeight      int64
}

type BackendMetrics struct {
	Mutex      sync.Mutex
	Metrics    *Metrics
	HealthPath string
}

type Config struct {
	Algorithm          string          `yaml:"algorithm"`
	Servers            []BackendConfig `yaml:"backends"`
	ConsecutiveFails   uint64          `yaml:"consecutiveFails"`
	ConsecutiveSuccess uint64          `yaml:"consecutiveSuccess"`
	FailRate           float64         `yaml:"failRate"`
	TimeOutRate        uint64          `yaml:"timeOutBreak"`
	TimeOutDelay       uint64          `yaml:"timeOutDelay"`
	AuthBasic          AuthConfig      `yaml:"auth"`
	SmartMode          bool            `yaml:"smartMode"`
	RateLimit          int64           `yaml:"rateLimit"`
}

type BackendConfig struct {
	UrlConfig        string `yaml:"url"`
	WeightConfig     int64  `yaml:"weight"`
	HealthPathConfig string `yaml:"healthPath"`
}

type AuthConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

var MetricsMap = map[string]*BackendMetrics{}

var ConfigSystem Config

var StateSystem string

var ActiveLogin bool

func InitServer() {
	StateSystem = Stable
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

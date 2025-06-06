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
	DefaultConsecutiveFails     = 10
	DefaultConsecutiveSuccess   = 10
	DefaultFailRate             = 0.1
	DefaultTimeOutRate          = 10
	DefaultTimeOutDelay         = 5
	DefaultHealthCheckInterval  = 1
	DefaultRateLimit            = 100
	DefaultWeight               = 1
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
	TimeOutRate        uint64          `yaml:"timeOutRate"`
	TimeOutDelay       uint64          `yaml:"timeOutDelay"`
	AuthBasic          AuthConfig      `yaml:"auth"`
	SmartMode          bool            `yaml:"smartMode"`
	RateLimit          int64           `yaml:"rateLimit"`
}

type DefaultSystem struct {
	StateSystem         string
	ActiveLogin         bool
	AlgorithmUserChoose string
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

var ConfigDefaultSystem DefaultSystem

func InitDefaultConfig() {
	ConfigSystem = Config{
		ConsecutiveFails:   DefaultConsecutiveFails,
		ConsecutiveSuccess: DefaultConsecutiveSuccess,
		FailRate:           DefaultFailRate,
		TimeOutRate:        DefaultTimeOutRate,
		TimeOutDelay:       DefaultTimeOutDelay,
		RateLimit:          DefaultRateLimit,
		SmartMode:          false,
	}
}

func SetupDefaultValues(cfg *Config) {
	if cfg.ConsecutiveFails == 0 {
		cfg.ConsecutiveFails = DefaultConsecutiveFails
	}
	if cfg.ConsecutiveSuccess == 0 {
		cfg.ConsecutiveSuccess = DefaultConsecutiveSuccess
	}
	if cfg.FailRate == 0 {
		cfg.FailRate = DefaultFailRate
	}
	if cfg.TimeOutRate == 0 {
		cfg.TimeOutRate = DefaultTimeOutRate
	}
	if cfg.TimeOutDelay == 0 {
		cfg.TimeOutDelay = DefaultTimeOutDelay
	}
	if cfg.RateLimit == 0 {
		cfg.RateLimit = DefaultRateLimit
	}
	for i := range cfg.Servers {
		if cfg.Servers[i].WeightConfig == 0 {
			cfg.Servers[i].WeightConfig = DefaultWeight
		}
	}
}

func InitServer() {
	ConfigDefaultSystem.StateSystem = Stable
	urls := ConfigSystem.Servers
	for _, url := range urls {
		weight := url.WeightConfig
		if weight == 0 {
			weight = DefaultWeight
		}
		MetricsMap[url.UrlConfig] = &BackendMetrics{
			Metrics: &Metrics{
				IsHealthy: true,
				Weight:    weight,
			},
			HealthPath: url.HealthPathConfig,
		}
	}
}

package main

import (
	"Go/balancer"
	"Go/config"
	"Go/factory"
	"Go/middleware"
	"Go/response"
	"fmt"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

func setupDefaultValues(cfg *config.Config) {
	if cfg.ConsecutiveFails == 0 {
		cfg.ConsecutiveFails = config.DefaultConsecutiveFails
	}
	if cfg.ConsecutiveSuccess == 0 {
		cfg.ConsecutiveSuccess = config.DefaultConsecutiveSuccess
	}
	if cfg.FailRate == 0 {
		cfg.FailRate = config.DefaultFailRate
	}
	if cfg.TimeOutRate == 0 {
		cfg.TimeOutRate = config.DefaultTimeOutRate
	}
	if cfg.TimeOutDelay == 0 {
		cfg.TimeOutDelay = config.DefaultTimeOutDelay
	}
	if cfg.HealthCheckInterval <= 0 {
		cfg.HealthCheckInterval = config.DefaultHealthCheckInterval
	}
	if cfg.RateLimit <= 0 {
		cfg.RateLimit = config.DefaultRateLimit
	}
	for i := range cfg.Servers {
		if cfg.Servers[i].WeightConfig == 0 {
			cfg.Servers[i].WeightConfig = config.DefaultWeight
		}
	}
}

func main() {
	configFile := "config.yaml"
	if envConfig := os.Getenv("CONFIG_FILE"); envConfig != "" {
		configFile = envConfig
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		os.Exit(1)
	}

	var cfg config.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		fmt.Printf("Error parsing config file: %v\n", err)
		os.Exit(1)
	}

	if cfg.Algorithm == "" || len(cfg.Servers) == 0 {
		fmt.Println("Missing required configuration fields")
		os.Exit(1)
	}

	setupDefaultValues(&cfg)

	config.ConfigSystem.Servers = make([]config.BackendConfig, len(cfg.Servers))
	for i, b := range cfg.Servers {
		config.ConfigSystem.Servers[i] = config.BackendConfig{
			UrlConfig:        b.UrlConfig,
			WeightConfig:     b.WeightConfig,
			HealthPathConfig: b.HealthPathConfig,
		}
	}

	config.ConfigSystem.Algorithm = cfg.Algorithm
	config.ConfigDefaultSystem.AlgorithmUserChoose = cfg.Algorithm
	config.ConfigSystem.ConsecutiveFails = cfg.ConsecutiveFails
	config.ConfigSystem.ConsecutiveSuccess = cfg.ConsecutiveSuccess
	config.ConfigSystem.FailRate = cfg.FailRate
	config.ConfigSystem.TimeOutRate = cfg.TimeOutRate
	config.ConfigSystem.TimeOutDelay = cfg.TimeOutDelay
	config.ConfigDefaultSystem.ActiveLogin = false
	config.ConfigSystem.SmartMode = cfg.SmartMode
	config.ConfigSystem.AuthBasic = config.AuthConfig{
		Username: cfg.AuthBasic.Username,
		Password: cfg.AuthBasic.Password,
	}
	config.ConfigSystem.HealthCheckInterval = cfg.HealthCheckInterval
	config.ConfigSystem.RateLimit = cfg.RateLimit

	config.NewDB(config.GetExecutableDir())
	config.InitServer()
	balancer.StartHealthCheck(time.Duration(config.ConfigSystem.HealthCheckInterval) * time.Second)
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimit)
	corsMiddleware := middleware.NewCORSMiddleware()
	factory.Factory = *factory.NewLoadBalancerFactory()
	http.Handle("/", rateLimiter.HandleRateLimit(corsMiddleware.HandleCORS(http.HandlerFunc(balancer.Handler))))
	http.Handle("/change-load-balancer", rateLimiter.HandleRateLimit(corsMiddleware.HandleCORS(http.HandlerFunc(balancer.ChangeAlgoLoadBalancer))))
	http.Handle("/metrics", rateLimiter.HandleRateLimit(corsMiddleware.HandleCORS(http.HandlerFunc(response.HandleStatusHTML))))
	http.Handle("/login-metrics", rateLimiter.HandleRateLimit(corsMiddleware.HandleCORS(http.HandlerFunc(balancer.Login))))
	http.Handle("/delete-error-history", rateLimiter.HandleRateLimit(corsMiddleware.HandleCORS(http.HandlerFunc(balancer.DeleteErrorHistory))))
	http.Handle("/error-history", rateLimiter.HandleRateLimit(corsMiddleware.HandleCORS(http.HandlerFunc(balancer.GetErrorHistory))))
	http.Handle("/reset-metrics", rateLimiter.HandleRateLimit(corsMiddleware.HandleCORS(http.HandlerFunc(balancer.ResetMetrics))))

	fmt.Printf("🚀 Load balancer running with rate limit of %d requests/second per IP\n", cfg.RateLimit)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("❌ Server failed to start:", err)
		os.Exit(1)
	}
}

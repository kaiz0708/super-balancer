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

	if cfg.Algorithm == "" || len(cfg.Servers) == 0 || (cfg.AuthBasic.Username == "" || cfg.AuthBasic.Password == "") {
		fmt.Println("Missing required configuration fields")
		os.Exit(1)
	}

	config.SetupDefaultValues(&cfg)
	config.ConfigSystem.Servers = make([]config.BackendConfig, len(cfg.Servers))
	for i, b := range cfg.Servers {
		config.ConfigSystem.Servers[i] = config.BackendConfig{
			UrlConfig:        b.UrlConfig,
			WeightConfig:     b.WeightConfig,
			HealthPathConfig: b.HealthPathConfig,
		}
	}

	config.ConfigSystem.AuthBasic = config.AuthConfig{
		Username: cfg.AuthBasic.Username,
		Password: cfg.AuthBasic.Password,
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
	config.ConfigSystem.RateLimit = cfg.RateLimit

	config.NewDB(config.GetExecutableDir())
	config.InitServer()
	balancer.StartHealthCheck(time.Duration(config.DefaultHealthCheckInterval) * time.Second)
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimit)
	corsMiddleware := middleware.NewCORSMiddleware()
	factory.Factory = *factory.NewLoadBalancerFactory()
	mux := http.NewServeMux()
	mux.Handle("/", rateLimiter.HandleRateLimit(corsMiddleware.HandleCORS(http.HandlerFunc(balancer.Handler))))
	mux.Handle("/change-load-balancer", rateLimiter.HandleRateLimit(corsMiddleware.HandleCORS(http.HandlerFunc(balancer.ChangeAlgoLoadBalancer))))
	mux.Handle("/metrics", rateLimiter.HandleRateLimit(corsMiddleware.HandleCORS(http.HandlerFunc(response.HandleStatusHTML))))
	mux.Handle("/login-metrics", rateLimiter.HandleRateLimit(corsMiddleware.HandleCORS(http.HandlerFunc(balancer.Login))))
	mux.Handle("/delete-state-history", rateLimiter.HandleRateLimit(corsMiddleware.HandleCORS(http.HandlerFunc(balancer.DeleteErrorHistory))))
	mux.Handle("/state-history", rateLimiter.HandleRateLimit(corsMiddleware.HandleCORS(http.HandlerFunc(balancer.GetErrorHistory))))
	mux.Handle("/reset-metrics", rateLimiter.HandleRateLimit(corsMiddleware.HandleCORS(http.HandlerFunc(balancer.ResetMetrics))))

	fmt.Printf("Load balancer running with rate limit of %d requests/second per IP\n", cfg.RateLimit)
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Server failed to start:", err)
		os.Exit(1)
	}
}

package main

import (
	test "Go/TestRequest"
	"Go/balancer"
	"Go/config"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

//go:embed static/*
var staticFiles embed.FS

func main() {

	validate := validator.New()

	r := gin.Default()

	subFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}
	r.StaticFS("/static", http.FS(subFS))

	r.POST("/api/backends", func(c *gin.Context) {
		var servers []config.BackendConfig
		if err := c.BindJSON(&servers); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		for _, server := range servers {
			if err := validate.Struct(server); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}
		for _, server := range servers {
			config.BackendServers = append(config.BackendServers, server)
			config.MetricsMap[server.UrlConfig] = &config.BackendMetrics{
				Metrics: &config.Metrics{
					Weight:    server.WeightConfig,
					IsHealthy: true,
				},
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "Backends added successfully"})
	})

	r.GET("/api/backends", func(c *gin.Context) {
		backends := make([]map[string]interface{}, 0)
		for _, server := range config.BackendServers {
			metric := config.MetricsMap[server.UrlConfig]
			if metric == nil {
				continue
			}
			metric.Mutex.Lock()
			failureRate := 0.0
			if metric.Metrics.RequestCount > 0 {
				failureRate = float64(metric.Metrics.FailureCount) / float64(metric.Metrics.RequestCount)
			}
			backends = append(backends, map[string]interface{}{
				"url":               server.UrlConfig,
				"healthy":           metric.Metrics.IsHealthy,
				"avgLatency":        metric.Metrics.AvgLatency.Milliseconds(),
				"requestCount":      metric.Metrics.RequestCount,
				"failureRate":       failureRate,
				"weight":            metric.Metrics.Weight,
				"activeConnections": metric.Metrics.ActiveConnections,
			})
			metric.Mutex.Unlock()
		}
		c.JSON(http.StatusOK, gin.H{
			"backends":  backends,
			"algorithm": config.LoadBalancerDefault,
		})
	})

	r.DELETE("/api/backends", func(c *gin.Context) {
		var req struct {
			Url string `json:"url"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newServers := []config.BackendConfig{}
		for _, server := range config.BackendServers {
			if server.UrlConfig != req.Url {
				newServers = append(newServers, server)
			}
		}
		config.BackendServers = newServers
		delete(config.MetricsMap, req.Url)
		c.JSON(http.StatusOK, gin.H{"message": "Backend removed"})
	})

	r.NoRoute(func(c *gin.Context) {
		balancer.Handler(c.Writer, c.Request)
	})

	r.POST("/change-load-balancer", func(c *gin.Context) {
		balancer.ChangeAlgoLoadBalancer(c.Writer, c.Request)
	})

	r.GET("/test", func(c *gin.Context) {
		test.SpamRequests(c)
	})

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	fmt.Println("🚀 Server running on :8080")
	if err := r.Run(":8080"); err != nil {
		fmt.Println("❌ Server failed to start:", err)
		os.Exit(1)
	}
}

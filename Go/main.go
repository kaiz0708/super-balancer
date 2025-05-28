package main

import (
	"Go/balancer"
	"Go/config"
	"Go/middleware"
	"Go/response"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	algorithm := flag.String("algorithm", "", "Load balancing algorithm")
	backends := flag.String("backends", "", "JSON array of backend servers")
	consecutiveFails := flag.String("consecutiveFails", "", "Set up amount consecutiveFails")
	consecutiveSuccess := flag.String("consecutiveSuccess", "", "Set up amount consecutiveSuccess")
	failRate := flag.String("failRate", "", "Set up rate fail")
	timeOutBreak := flag.String("timeOutBreak", "", "Set up timeout break")
	timeOutDelay := flag.String("timeOutDelay", "", "Set up timeout delay")
	auth := flag.String("auth", "", "Set up auth")
	smartMode := flag.String("smartMode", "", "Set up smart mode")
	rateLimit := flag.String("rateLimit", "1000", "Maximum requests per second per IP")

	flag.Parse()

	if *algorithm == "" ||
		*backends == "" ||
		*consecutiveFails == "" ||
		*consecutiveSuccess == "" ||
		*failRate == "" ||
		*timeOutBreak == "" ||
		*timeOutDelay == "" {
		fmt.Println("Missing one or more required flags")
		os.Exit(1)
	}

	rateLimitValue, err := strconv.Atoi(*rateLimit)
	if err != nil {
		fmt.Println("Invalid rate limit value:", err)
		os.Exit(1)
	}

	consecutiveFailsValue, err := strconv.ParseUint(*consecutiveFails, 10, 64)
	if err != nil {
		fmt.Println("Lỗi chuyển đổi:", err)
		return
	}

	consecutiveSuccessValue, err := strconv.ParseUint(*consecutiveSuccess, 10, 64)
	if err != nil {
		fmt.Println("Lỗi chuyển đổi:", err)
		return
	}

	timeOutDelayValue, err := strconv.ParseUint(*timeOutDelay, 10, 64)
	if err != nil {
		fmt.Println("Lỗi chuyển đổi:", err)
		return
	}

	faileRateValue, err := strconv.ParseFloat(*failRate, 64)
	if err != nil {
		fmt.Println("Lỗi chuyển đổi:", err)
		return
	}

	timeOutBreakValue, err := strconv.ParseUint(*timeOutBreak, 10, 64)
	if err != nil {
		fmt.Println("Lỗi chuyển đổi:", err)
		return
	}

	if err := json.Unmarshal([]byte(*backends), &config.ConfigSystem.Servers); err != nil {
		fmt.Println("Failed to parse backends:", err)
		os.Exit(1)
	}

	if err := json.Unmarshal([]byte(*auth), &config.ConfigSystem.AuthConfig); err != nil {
		fmt.Println("Failed to parse auth:", err)
		os.Exit(1)
	}

	smartModeValue, err := strconv.ParseBool(*smartMode)
	if err != nil {
		fmt.Println("Invalid value for smartMode: ", err)
	}

	config.ConfigSystem.Algorithm = *algorithm
	config.ConfigSystem.ConsecutiveFails = consecutiveFailsValue
	config.ConfigSystem.ConsecutiveSuccess = consecutiveSuccessValue
	config.ConfigSystem.FailRate = faileRateValue
	config.ConfigSystem.TimeOutRate = timeOutBreakValue
	config.ConfigSystem.TimeOutDelay = timeOutDelayValue
	config.ConfigSystem.ActiveLogin = false
	config.ConfigSystem.SmartMode = smartModeValue
	config.NewDB(config.GetExecutableDir())
	config.InitServer()
	balancer.StartHealthCheck(1 * time.Second)

	rateLimiter := middleware.NewRateLimiter(rateLimitValue)
	corsMiddleware := middleware.NewCORSMiddleware()
	http.Handle("/", rateLimiter.HandleRateLimit(corsMiddleware.HandleCORS(http.HandlerFunc(balancer.Handler))))
	http.Handle("/change-load-balancer", rateLimiter.HandleRateLimit(corsMiddleware.HandleCORS(http.HandlerFunc(balancer.ChangeAlgoLoadBalancer))))
	http.Handle("/metrics", rateLimiter.HandleRateLimit(corsMiddleware.HandleCORS(http.HandlerFunc(response.HandleStatusHTML))))
	http.Handle("/login-metrics", rateLimiter.HandleRateLimit(corsMiddleware.HandleCORS(http.HandlerFunc(balancer.Login))))
	http.Handle("/delete-error-history", rateLimiter.HandleRateLimit(corsMiddleware.HandleCORS(http.HandlerFunc(balancer.DeleteErrorHistory))))
	http.Handle("/error-history", rateLimiter.HandleRateLimit(corsMiddleware.HandleCORS(http.HandlerFunc(balancer.GetErrorHistory))))
	http.Handle("/reset-metrics", rateLimiter.HandleRateLimit(corsMiddleware.HandleCORS(http.HandlerFunc(balancer.ResetMetrics))))

	fmt.Printf("🚀 Load balancer running on :8080 with rate limit of %d requests/second per IP\n", rateLimitValue)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("❌ Server failed to start:", err)
		os.Exit(1)
	}
}

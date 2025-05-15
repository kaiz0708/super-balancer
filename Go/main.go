package main

import (
	test "Go/TestRequest"
	"Go/balancer"
	"Go/config"
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

	defaultProxy := flag.String("defaultProxy", "", "Default proxy path")
	algorithm := flag.String("algorithm", "", "Load balancing algorithm")
	backends := flag.String("backends", "", "JSON array of backend servers")
	consecutiveFails := flag.String("consecutiveFails", "", "Set up amount consecutiveFails")
	consecutiveSuccess := flag.String("consecutiveSuccess", "", "Set up amount consecutiveSuccess")
	failRate := flag.String("failRate", "", "Set up rate fail")
	timeOutBreak := flag.String("timeOutBreak", "", "Set up timeout break")
	timeOutDelay := flag.String("timeOutDelay", "", "Set up timeout delay")
	auth := flag.String("auth", "", "Set up auth")

	flag.Parse()

	if *defaultProxy == "" ||
		*algorithm == "" ||
		*backends == "" ||
		*consecutiveFails == "" ||
		*consecutiveSuccess == "" ||
		*failRate == "" ||
		*timeOutBreak == "" ||
		*timeOutDelay == "" {
		fmt.Println("Missing one or more required flags")
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

	config.ConfigSystem.DefaultProxy = *defaultProxy
	config.ConfigSystem.Algorithm = *algorithm
	config.ConfigSystem.ConsecutiveFails = consecutiveFailsValue
	config.ConfigSystem.ConsecutiveSuccess = consecutiveSuccessValue
	config.ConfigSystem.FailRate = faileRateValue
	config.ConfigSystem.TimeOutRate = timeOutBreakValue
	config.ConfigSystem.TimeOutDelay = timeOutDelayValue
	config.ConfigSystem.ActiveLogin = false
	config.NewDB(config.GetExecutableDir())
	config.InitServer()
	balancer.StartHealthCheck(1 * time.Second)
	http.HandleFunc(config.ConfigSystem.DefaultProxy, balancer.Handler)
	http.HandleFunc(config.ConfigSystem.DefaultProxy+"change-load-balancer", balancer.ChangeAlgoLoadBalancer)
	http.HandleFunc(config.ConfigSystem.DefaultProxy+"test", test.SpamRequests)
	http.HandleFunc("/metrics", response.HandleStatusHTML)
	http.HandleFunc("/login-metrics", balancer.Login)
	http.HandleFunc("/delete-error-history", balancer.DeleteErrorHistory)
	http.HandleFunc("/error-history", balancer.GetErrorHistory)
	fmt.Println("🚀 Load balancer running on :8080")
	err = http.ListenAndServe(config.ConfigSystem.DefaultProxy, nil)
	if err != nil {
		fmt.Println("❌ Server failed to start:", err)
		os.Exit(1)
	}
}

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

	"github.com/go-playground/validator/v10"
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
	passwordAccessMetrics := flag.String("passwordAccessMetrics", "", "Set up Password to access metrics")

	flag.Parse()

	fmt.Println("========== Flag Values =========")
	fmt.Println("defaultProxy:", *defaultProxy)
	fmt.Println("algorithm:", *algorithm)
	fmt.Println("backends:", *backends)
	fmt.Println("consecutiveFails:", *consecutiveFails)
	fmt.Println("consecutiveSuccess:", *consecutiveSuccess)
	fmt.Println("failRate:", *failRate)
	fmt.Println("timeOutBreak:", *timeOutBreak)
	fmt.Println("timeOutDelay:", *timeOutDelay)

	if *defaultProxy == "" ||
		*algorithm == "" ||
		*backends == "" ||
		*consecutiveFails == "" ||
		*consecutiveSuccess == "" ||
		*failRate == "" ||
		*timeOutBreak == "" ||
		*timeOutDelay == "" ||
		*passwordAccessMetrics == "" {
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

	var cfg config.Config
	cfg.DefaultProxy = *defaultProxy
	cfg.Algorithm = *algorithm

	if err := json.Unmarshal([]byte(*backends), &cfg.Servers); err != nil {
		fmt.Println("Failed to parse backends:", err)
		os.Exit(1)
	}

	validate := validator.New()
	err = validate.Struct(cfg)
	if err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			fmt.Printf("Invalid field: %s\n", fieldErr.StructNamespace())
		}
		os.Exit(1)
	}

	config.BackendServers = append(config.BackendServers, cfg.Servers...)
	config.LoadBalancerDefault = cfg.Algorithm
	config.ConsecutiveFails = consecutiveFailsValue
	config.ConsecutiveSuccess = consecutiveSuccessValue
	config.FailRate = faileRateValue
	config.TimeOutRate = timeOutBreakValue
	config.TimeOutDelay = timeOutDelayValue
	config.PasswordAccessMetrics = *passwordAccessMetrics
	if err != nil {
		fmt.Println("Lỗi chuyển đổi:", err)
		return
	}
	config.InitServer()

	http.HandleFunc(cfg.DefaultProxy, balancer.Handler)
	http.HandleFunc(cfg.DefaultProxy+"change-load-balancer", balancer.ChangeAlgoLoadBalancer)
	http.HandleFunc(cfg.DefaultProxy+"test", test.SpamRequests)
	http.HandleFunc("/metrics", response.HandleStatusHTML)
	fmt.Println("🚀 Load balancer running on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("❌ Server failed to start:", err)
		os.Exit(1)
	}
}

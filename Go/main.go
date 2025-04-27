package main

import (
	test "Go/TestRequest"
	"Go/balancer"
	"Go/config"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
)

func main() {

	defaultProxy := flag.String("defaultProxy", "", "Default proxy path")
	algorithm := flag.String("algorithm", "", "Load balancing algorithm")
	backends := flag.String("backends", "", "JSON array of backend servers")

	flag.Parse()

	if *defaultProxy == "" || *algorithm == "" || *backends == "" {
		fmt.Println("Missing one or more required flags")
		os.Exit(1)
	}

	var cfg config.Config
	cfg.DefaultProxy = *defaultProxy
	cfg.Algorithm = *algorithm

	if err := json.Unmarshal([]byte(*backends), &cfg.Servers); err != nil {
		fmt.Println("Failed to parse backends:", err)
		os.Exit(1)
	}

	validate := validator.New()
	err := validate.Struct(cfg)
	if err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			fmt.Printf("Invalid field: %s\n", fieldErr.StructNamespace())
		}
		os.Exit(1)
	}

	config.BackendServers = append(config.BackendServers, cfg.Servers...)
	config.LoadBalancerDefault = cfg.Algorithm
	config.InitServer()

	http.HandleFunc(cfg.DefaultProxy, balancer.Handler)
	http.HandleFunc(cfg.DefaultProxy+"change-load-balancer", balancer.ChangeAlgoLoadBalancer)
	http.HandleFunc(cfg.DefaultProxy+"test", test.SpamRequests)
	fmt.Println("🚀 Load balancer running on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("❌ Server failed to start:", err)
		os.Exit(1)
	}
}

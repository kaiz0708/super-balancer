package main

import (
	"Go/balancer"
	"Go/config"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config.InitServer()
	http.HandleFunc("/", balancer.Handler)
	http.HandleFunc("/change-load-balancer", balancer.ChangeAlgoLoadBalancer)
	println("Load balancer running on :8080")
	http.ListenAndServe(":8080", nil)
}

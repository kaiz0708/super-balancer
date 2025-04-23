package balancer

import (
	"Go/algo"
	"Go/config"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func updateMetricsBackend(backend string, start time.Time, resp *http.Response) {
	latency := time.Since(start)
	success := true
	if resp.StatusCode >= 500 {
		success = false
	}
	UpdateMetrics(backend, latency, success, resp.StatusCode)
}

func updateActiveConnection(backend string, state bool) {
	UpdateActiveConnectionMetrics(backend, state)
}

func ChangeAlgoLoadBalancer(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	config.LoadBalancerDefault = name

	json.NewEncoder(w).Encode(name)

}

func TestRequest(w http.ResponseWriter, r *http.Request) {

}

func analyzeSystemState() string {
	backends := config.MetricsMap
	healthyCount := 0
	highLatencyCount := 0
	totalBackends := len(backends)

	for _, backend := range backends {
		if backend.Metrics.IsHealthy && backend.Metrics.ConsecutiveFails < 3 {
			healthyCount++
		}
		if backend.Metrics.LastLatency > 500*time.Millisecond || backend.Metrics.AvgLatency > 500*time.Millisecond {
			highLatencyCount++
		}
	}

	if healthyCount == 0 {
		return "AllFailed"
	}
	if float64(healthyCount)/float64(totalBackends) < 0.5 {
		return "ManyFailed"
	}
	if highLatencyCount > 1 {
		return "HighLatency"
	}
	return "Stable"
}

func Handler(w http.ResponseWriter, r *http.Request) {
	algoChoice := config.LoadBalancerDefault
	target := algo.AlgoLoadBalancer(algoChoice)

	url, err := url.Parse(target)

	if err != nil {
		http.Error(w, "Invalid backend URL", http.StatusInternalServerError)
		return
	}

	start := time.Now()
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ModifyResponse = func(resp *http.Response) error {
		go func() {
			pickState := analyzeSystemState()
			algo.ChooseAlgorithm(pickState)
		}()
		go updateActiveConnection(url.String(), false)
		updateMetricsBackend(target, start, resp)
		return nil
	}
	r.Host = url.Host
	go updateActiveConnection(url.String(), true)
	proxy.ServeHTTP(w, r)
}

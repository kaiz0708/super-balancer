package balancer

import (
	"Go/algo"
	"Go/algo/custom"
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
			pickState := AnalyzeSystemState(url.String())
			if pickState == "AllFailed" {
				custom.CustomAllFailed(w)
			} else {
				algo.ChooseAlgorithm(pickState)
			}
		}()
		go updateActiveConnection(url.String(), false)
		go updateMetricsBackend(target, start, resp)
		return nil
	}
	r.Host = url.Host
	go updateActiveConnection(url.String(), true)
	proxy.ServeHTTP(w, r)
}

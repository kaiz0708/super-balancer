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

func updateMetricsBackend(backend string, start time.Time, statusCode int) {
	latency := time.Since(start)
	success := true
	if statusCode == 502 {
		success = false
	}
	UpdateMetrics(backend, latency, success, statusCode)
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
	target := ""
	pickState := AnalyzeSystemState()

	if pickState == "AllFailed" {
		custom.CustomAllFailed(w)
	} else {
		target = algo.ChooseAlgorithm(pickState)
	}

	url, err := url.Parse(target)

	if err != nil {
		http.Error(w, "Invalid backend URL", http.StatusInternalServerError)
		return
	}

	start := time.Now()
	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.ModifyResponse = func(resp *http.Response) error {
		go updateActiveConnection(url.String(), false)
		go updateMetricsBackend(target, start, resp.StatusCode)
		return nil
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		go updateActiveConnection(url.String(), false)
		go updateMetricsBackend(target, start, 502)
	}

	r.Host = url.Host
	updateActiveConnection(url.String(), true)
	proxy.ServeHTTP(w, r)
}

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

func UpdateMetricsBackend(backend string, start time.Time, statusCode int) {
	latency := time.Since(start)
	success := true
	if statusCode == 502 {
		success = false
	}
	UpdateMetrics(backend, latency, success, statusCode)
}

func UpdateActiveConnection(backend string, state bool) {
	UpdateActiveConnectionMetrics(backend, state)
}

func CheckUnhealthyBackend(w http.ResponseWriter, r *http.Request) {
	backends := config.MetricsMap

	for backend, m := range backends {
		if !m.Metrics.IsHealthy {
			url, _ := url.Parse(backend)
			start := time.Now()
			proxy := httputil.NewSingleHostReverseProxy(url)
			proxy.ModifyResponse = func(resp *http.Response) error {
				UpdateMetricsBackend(backend, start, resp.StatusCode)
				return nil
			}
			r.Host = url.Host
			proxy.ServeHTTP(w, r)
		}
	}
}

func HttpProxy(backend string, w http.ResponseWriter, r *http.Request) {
	url, err := url.Parse(backend)

	if err != nil {
		http.Error(w, "Invalid backend URL", http.StatusInternalServerError)
		return
	}

	start := time.Now()
	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.ModifyResponse = func(resp *http.Response) error {
		go UpdateActiveConnection(url.String(), false)
		go UpdateMetricsBackend(backend, start, resp.StatusCode)
		return nil
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		go UpdateActiveConnection(url.String(), false)
		go UpdateMetricsBackend(backend, start, 502)
	}

	r.Host = url.Host
	UpdateActiveConnection(url.String(), true)
	proxy.ServeHTTP(w, r)
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
	go CheckUnhealthyBackend(w, r)
	HttpProxy(target, w, r)
}

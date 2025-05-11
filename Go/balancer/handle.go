package balancer

import (
	"Go/algo"
	"Go/config"
	"Go/response"
	"encoding/json"
	"fmt"
	"net"
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

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == config.AuthConfig.Username && password == config.AuthConfig.Password {
		config.ActiveLogin = true
		http.Redirect(w, r, "/metrics", http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}

func UpdateActiveConnection(backend string, state bool) {
	UpdateActiveConnectionMetrics(backend, state)
}

func CheckUnhealthyBackend() {
	for backend, m := range config.MetricsMap {
		if !m.Metrics.IsHealthy && m.HealthPath != "" {
			go func(backend string, m *config.BackendMetrics) {
				url := backend + m.HealthPath
				start := time.Now()

				resp, err := http.Get(url)
				if err != nil {
					fmt.Printf("Failed health check for %s: %v\n", backend, err)
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode == 200 {
					m.Mutex.Lock()
					config.MetricsMap[backend].Metrics.ConsecutiveSuccess++
					m.Mutex.Unlock()
					UpdateMetricsBackend(backend, start, resp.StatusCode)
				}
			}(backend, m)
		}
	}
}

func HttpProxy(backend string, w http.ResponseWriter, r *http.Request) {
	url, err := url.Parse(backend)

	transport := &http.Transport{
		ResponseHeaderTimeout: time.Duration(config.TimeOutDelay) * time.Second,
	}

	if err != nil {
		http.Error(w, "Invalid backend URL", http.StatusInternalServerError)
		return
	}

	start := time.Now()
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = transport

	proxy.ModifyResponse = func(resp *http.Response) error {
		go UpdateActiveConnection(backend, false)
		go UpdateMetricsBackend(backend, start, resp.StatusCode)
		return nil
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			config.MetricsMap[backend].Mutex.Lock()
			config.MetricsMap[backend].Metrics.TimeoutBreak++
			config.MetricsMap[backend].Mutex.Unlock()
		}
		go UpdateActiveConnection(backend, false)
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

func StartHealthCheck(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			CheckUnhealthyBackend()
		}
	}()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	CheckUnhealthyBackend()
	pickState := AnalyzeSystemState()
	if pickState == "AllFailed" {
		response.CustomAllFailed(w)
		return
	}
	target := algo.ChooseAlgorithm(pickState)
	HttpProxy(target, w, r)
}

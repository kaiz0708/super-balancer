package balancer

import (
	"Go/algo"
	"Go/config"
	"Go/response"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == config.ConfigSystem.AuthConfig.Username && password == config.ConfigSystem.AuthConfig.Password {
		config.ConfigSystem.ActiveLogin = true
		http.Redirect(w, r, "/metrics", http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}

func CheckUnhealthyBackend() {
	for backend, m := range config.MetricsMap {
		if !m.Metrics.IsHealthy && m.HealthPath != "" {
			go func(backend string, m *config.BackendMetrics) {
				url := backend + m.HealthPath

				resp, err := http.Get(url)
				if err != nil {
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode == 200 {
					m.Mutex.Lock()
					config.MetricsMap[backend].Metrics.ConsecutiveSuccess++
					m.Mutex.Unlock()
					UpdateBackendRecovering(backend)
				}
			}(backend, m)
		}
	}
}

func HttpProxy(backend string, w http.ResponseWriter, r *http.Request) {
	url, err := url.Parse(backend)

	transport := &http.Transport{
		ResponseHeaderTimeout: time.Duration(config.ConfigSystem.TimeOutDelay) * time.Second,
	}

	if err != nil {
		http.Error(w, "Invalid backend URL", http.StatusInternalServerError)
		return
	}

	start := time.Now()
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = transport

	proxy.ModifyResponse = func(resp *http.Response) error {
		UpdateActiveConnectionMetrics(backend, false)
		UpdateMetrics(backend, time.Since(start), resp.StatusCode)
		return nil
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			config.MetricsMap[backend].Mutex.Lock()
			config.MetricsMap[backend].Metrics.TimeoutBreak++
			config.MetricsMap[backend].Mutex.Unlock()
		}
		UpdateActiveConnectionMetrics(backend, false)
		UpdateBackendUnhealthy(backend, 502)
	}

	r.Host = url.Host
	UpdateActiveConnectionMetrics(url.String(), true)
	proxy.ServeHTTP(w, r)
}

func ChangeAlgoLoadBalancer(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	config.ConfigSystem.Algorithm = name

	json.NewEncoder(w).Encode(name)
}

func DeleteErrorHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		ID int64 `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := config.GlobalDB.DeleteErrorHistory(request.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func GetErrorHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	errors := config.GlobalDB.ReadMetrics()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(errors)
}

func ResetMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		BackendID string `json:"backend_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, exists := config.MetricsMap[request.BackendID]
	if !exists {
		http.Error(w, "Backend not found", http.StatusNotFound)
		return
	}

	UpdateResetMetrics(request.BackendID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
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
	if pickState == config.AllFailed {
		response.CustomAllFailed(w)
		return
	}
	target := algo.ChooseAlgorithm(pickState, r)
	HttpProxy(target, w, r)
}

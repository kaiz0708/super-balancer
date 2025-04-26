package balancer

import (
	"Go/config"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

var countRequestLock sync.Mutex
var totalRequests uint64 = 0

const updateEvery = 1

func UpdateMetrics(backend string, latency time.Duration, success bool, status int) {
	backendMetric := config.MetricsMap[backend]
	backendMetric.Mutex.Lock()
	defer backendMetric.Mutex.Unlock()
	m := backendMetric.Metrics
	m.RequestCount++
	m.TotalLatency += latency
	m.LastLatency = latency
	m.LastChecked = time.Now()
	m.LastStatus = status

	if m.RequestCount > 0 {
		m.AvgLatency = m.TotalLatency / time.Duration(m.RequestCount)
	}

	if success {
		m.SuccessCount++
		if m.ConsecutiveFails >= 5 {
			m.AvgLatency = 0
			m.TotalLatency = 0
			m.RequestCount = 0
			m.FailureCount = 0
		}
		m.ConsecutiveFails = 0
		m.IsHealthy = true
	} else {
		m.FailureCount++
		m.ConsecutiveFails++
	}

	failRate := float64(m.FailureCount) / float64(m.RequestCount)
	if m.ConsecutiveFails >= 5 || failRate >= 0.5 || m.AvgLatency >= 500*time.Millisecond {
		m.IsHealthy = false
	}

	countRequestLock.Lock()
	totalRequests++
	countRequestLock.Unlock()

	if totalRequests%updateEvery == 0 {
		clearTerminal()
		logInforBackend()
	}
}

func logInforBackend() {
	data, _ := json.Marshal(config.MetricsMap)
	fmt.Println(string(data))
}

func UpdateActiveConnectionMetrics(backend string, state bool) {
	backendMetric := config.MetricsMap[backend]
	backendMetric.Mutex.Lock()
	defer backendMetric.Mutex.Unlock()
	m := backendMetric.Metrics
	if state {
		m.ActiveConnections++
	} else {
		m.ActiveConnections--
	}
}

func clearTerminal() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Print("\033[2J\033[H")
	}
}

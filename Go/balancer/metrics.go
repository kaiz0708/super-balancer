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

	if success {
		m.SuccessCount++
		m.ConsecutiveFails = 0
		m.IsHealthy = true
	} else {
		m.FailureCount++
		m.ConsecutiveFails++
		if m.ConsecutiveFails >= 3 {
			m.IsHealthy = false
		}
	}

	if m.RequestCount > 0 {
		m.AvgLatency = m.TotalLatency / time.Duration(m.RequestCount)
	}

	countRequestLock.Lock()
	totalRequests++
	countRequestLock.Unlock()

	if totalRequests%updateEvery == 0 {
		clearTerminal()
		logInforBackend()
	}
}

func ResetMetricsAfterHealthy() {

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

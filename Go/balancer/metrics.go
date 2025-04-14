package balancer

import (
	"Go/config"
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
		logMetrics()
	}
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
		fmt.Print("\033[2J\033[H") // ANSI escape code to clear screen
	}
}

func logMetrics() {
	fmt.Println("======== Backend Metrics ========")
	for backend, m := range config.MetricsMap {
		if _, ok := config.MetricsMap[backend]; ok {
			target := m.Metrics
			successRate := float64(target.SuccessCount) / float64(target.RequestCount)
			errorRate := float64(target.FailureCount) / float64(target.RequestCount)
			avgLatencyMs := float64(target.AvgLatency.Microseconds()) / 1000.0
			lastLatencyMs := float64(target.LastLatency.Microseconds()) / 1000.0
			fmt.Printf("Backend: %s\n", backend)
			fmt.Printf("  - Requests:        %d\n", target.RequestCount)
			fmt.Printf("  - Successes:       %d\n", target.SuccessCount)
			fmt.Printf("  - Failures:        %d\n", target.FailureCount)
			fmt.Printf("  - Success Rate:    %.2f%%\n", successRate*100)
			fmt.Printf("  - Error Rate:      %.2f%%\n", errorRate*100)
			fmt.Printf("  - Avg Latency:     %.2f ms\n", avgLatencyMs)
			fmt.Printf("  - Last Latency:    %.2f ms\n", lastLatencyMs)
			fmt.Printf("  - Healthy:         %v\n", target.IsHealthy)
			fmt.Printf("  - Last Status:     %d\n", target.LastStatus)
			fmt.Printf("  - Last Checked:    %s\n", target.LastChecked.Format(time.RFC3339))
			fmt.Printf("  - ActiveConnection:%d\n", target.ActiveConnections)
			fmt.Printf("  - Weight:        	 %d\n", target.Weight)
			fmt.Printf("  - CurrentWeight:   %d\n", target.CurrentWeight)
			fmt.Println("---------------------------------------")
		}
	}
}

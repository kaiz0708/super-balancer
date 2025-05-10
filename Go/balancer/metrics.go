package balancer

import (
	"Go/algo"
	"Go/config"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var countRequestLock sync.Mutex
var TotalRequests uint64 = 0

const updateEvery = 100

func UpdateMetrics(backend string, latency time.Duration, success bool, status int) {
	backendMetric := config.MetricsMap[backend]
	backendMetric.Mutex.Lock()
	defer backendMetric.Mutex.Unlock()
	m := backendMetric.Metrics
	if m.IsHealthy {
		m.RequestCount++
		m.TotalLatency += latency
		m.LastLatency = latency
		m.LastChecked = time.Now()
		m.LastStatus = status

		if m.RequestCount > 0 {
			m.AvgLatency = m.TotalLatency / time.Duration(m.RequestCount)
		}
	}

	if m.ConsecutiveSuccess >= config.ConsecutiveSuccess {
		m.IsHealthy = true
		m.ConsecutiveFails = 0
		m.AvgLatency = 0
		m.TotalLatency = 0
		m.LastLatency = 0
		m.TimeoutBreak = 0
		m.ConsecutiveSuccess = 0
	}

	if success {
		m.SuccessCount++
		m.ConsecutiveFails = 0
	} else {
		m.FailureCount++
		m.ConsecutiveFails++
	}

	failRate := float64(m.FailureCount) / float64(m.RequestCount)
	if m.ConsecutiveFails >= config.ConsecutiveFails || failRate >= config.FailRate || m.TimeoutBreak >= config.TimeOutRate {
		m.IsHealthy = false
	}

	countRequestLock.Lock()
	TotalRequests++
	countRequestLock.Unlock()

	if TotalRequests%updateEvery == 0 {
		logInforBackend()
		fmt.Println("Algo current : ", algo.AlgoCurrent)
	}
}

func logInforBackend() {
	data, _ := json.Marshal(config.MetricsMap)
	fmt.Println("Log backends infor : ", string(data))
}

func UpdateActiveConnectionMetrics(backend string, state bool) {
	backendMetric := config.MetricsMap[backend]
	m := backendMetric.Metrics
	if state {
		atomic.AddInt64(&m.ActiveConnections, 1)
	} else {
		atomic.AddInt64(&m.ActiveConnections, -1)
	}
}

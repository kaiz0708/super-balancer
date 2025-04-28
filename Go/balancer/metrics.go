package balancer

import (
	"Go/algo"
	"Go/config"
	"encoding/json"
	"fmt"
	"sync"
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

	if success {
		m.SuccessCount++
		if m.ConsecutiveFails >= config.ConsecutiveFails {
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
	if m.ConsecutiveFails >= config.ConsecutiveFails || failRate >= config.FailRate || m.AvgLatency >= 500*time.Millisecond || m.LastLatency >= 1*time.Second {
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
	backendMetric.Mutex.Lock()
	defer backendMetric.Mutex.Unlock()
	m := backendMetric.Metrics
	if state {
		m.ActiveConnections++
	} else {
		m.ActiveConnections--
	}
}

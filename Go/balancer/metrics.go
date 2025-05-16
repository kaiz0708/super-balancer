package balancer

import (
	"Go/config"
	"sync"
	"sync/atomic"
	"time"
)

var countRequestLock sync.Mutex
var TotalRequests uint64 = 0

func UpdateMetrics(backend string, latency time.Duration, status int) {
	backendMetric := config.MetricsMap[backend]
	backendMetric.Mutex.Lock()
	defer backendMetric.Mutex.Unlock()
	m := backendMetric.Metrics
	m.SuccessCount++
	m.ConsecutiveFails = 0
	m.RequestCount++
	m.TotalLatency += latency
	m.LastLatency = latency
	m.LastChecked = time.Now()
	m.LastStatus = status

	if m.RequestCount > 0 {
		m.AvgLatency = m.TotalLatency / time.Duration(m.RequestCount)
	}

	countRequestLock.Lock()
	TotalRequests++
	countRequestLock.Unlock()
}

func UpdateBackendUnhealthy(backend string, status int) {
	backendMetric := config.MetricsMap[backend]
	backendMetric.Mutex.Lock()
	defer backendMetric.Mutex.Unlock()
	m := backendMetric.Metrics
	m.FailureCount++
	m.ConsecutiveFails++
	m.LastStatus = status
	failRate := float64(m.FailureCount) / float64(m.RequestCount)
	if m.ConsecutiveFails >= config.ConfigSystem.ConsecutiveFails || failRate >= config.ConfigSystem.FailRate || m.TimeoutBreak >= config.ConfigSystem.TimeOutRate {
		m.IsHealthy = false
		config.GlobalDB.InsertMetrics(backend, config.Unhealthy, m)
	}
}

func UpdateBackendRecovering(backend string) {
	backendMetric := config.MetricsMap[backend]
	backendMetric.Mutex.Lock()
	defer backendMetric.Mutex.Unlock()
	m := backendMetric.Metrics
	if m.ConsecutiveSuccess >= config.ConfigSystem.ConsecutiveSuccess {
		m.IsHealthy = true
		m.ConsecutiveFails = 0
		m.AvgLatency = 0
		m.TotalLatency = 0
		m.LastLatency = 0
		m.TimeoutBreak = 0
		m.ConsecutiveSuccess = 0
		m.SuccessCount = 0
		m.RequestCount = 1
		m.FailureCount = 0
		config.GlobalDB.InsertMetrics(backend, config.Recovery, m)
	}
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

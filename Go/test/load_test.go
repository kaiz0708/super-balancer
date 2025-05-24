package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	// Số lượng request đồng thời - giảm xuống còn 50
	concurrentRequests := 50
	// Tổng số request - giảm xuống còn 1000
	totalRequests := 1000
	// URL của load balancer
	lbURL := "http://localhost:8080"

	var wg sync.WaitGroup
	results := make(chan string, totalRequests)
	
	// Thêm biến để theo dõi thống kê
	var (
		totalLatency time.Duration
		successCount int
		errorCount   int
		mu           sync.Mutex
	)

	startTime := time.Now()

	// Tạo worker pool
	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < totalRequests/concurrentRequests; j++ {
				start := time.Now()
				resp, err := http.Get(lbURL)
				if err != nil {
					mu.Lock()
					errorCount++
					mu.Unlock()
					results <- fmt.Sprintf("Error: %v", err)
					continue
				}
				defer resp.Body.Close()
				
				latency := time.Since(start)
				mu.Lock()
				totalLatency += latency
				successCount++
				mu.Unlock()

				results <- fmt.Sprintf("Worker %d, Request %d: Status %d, Latency %v", 
					workerID, j, resp.StatusCode, latency)
			}
		}(i)
	}

	// Đóng channel khi tất cả worker hoàn thành
	go func() {
		wg.Wait()
		close(results)
	}()

	// In kết quả
	for result := range results {
		fmt.Println(result)
	}

	// In thống kê
	totalTime := time.Since(startTime)
	fmt.Printf("\n=== Test Statistics ===\n")
	fmt.Printf("Total Requests: %d\n", totalRequests)
	fmt.Printf("Concurrent Requests: %d\n", concurrentRequests)
	fmt.Printf("Total Time: %v\n", totalTime)
	fmt.Printf("Successful Requests: %d\n", successCount)
	fmt.Printf("Failed Requests: %d\n", errorCount)
	if successCount > 0 {
		fmt.Printf("Average Latency: %v\n", totalLatency/time.Duration(successCount))
	}
	fmt.Printf("Requests per second: %.2f\n", float64(successCount)/totalTime.Seconds())
} 
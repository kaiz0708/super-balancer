package middleware

import (
	"Go/utils"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	requestsPerSecond int64
	windowSize        time.Duration
	requests          map[string][]time.Time
	mu                sync.Mutex
}

func NewRateLimiter(requestsPerSecond int64) *RateLimiter {
	return &RateLimiter{
		requestsPerSecond: requestsPerSecond,
		windowSize:        time.Second,
		requests:          make(map[string][]time.Time),
	}
}

func (rl *RateLimiter) HandleRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := utils.GetClientIP(r)

		rl.mu.Lock()
		now := time.Now()

		if requests, exists := rl.requests[clientIP]; exists {
			var validRequests []time.Time
			for _, reqTime := range requests {
				if now.Sub(reqTime) <= rl.windowSize {
					validRequests = append(validRequests, reqTime)
				}
			}
			rl.requests[clientIP] = validRequests
		}

		if int64(len(rl.requests[clientIP])) > rl.requestsPerSecond {
			rl.mu.Unlock()
			log.Println("Rate limit exceeded")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "Rate limit exceeded",
				"message": "Too many requests, please try again later",
			})
			return
		}

		rl.requests[clientIP] = append(rl.requests[clientIP], now)
		rl.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

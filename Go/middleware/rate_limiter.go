package middleware

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	requestsPerSecond int
	windowSize        time.Duration
	requests          map[string][]time.Time
	mu                sync.Mutex
}

func NewRateLimiter(requestsPerSecond int) *RateLimiter {
	return &RateLimiter{
		requestsPerSecond: requestsPerSecond,
		windowSize:        time.Second,
		requests:          make(map[string][]time.Time),
	}
}

func (rl *RateLimiter) HandleRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr

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

		if len(rl.requests[clientIP]) >= rl.requestsPerSecond {
			rl.mu.Unlock()
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		rl.requests[clientIP] = append(rl.requests[clientIP], now)
		rl.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

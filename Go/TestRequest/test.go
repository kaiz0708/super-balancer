package test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func CallRequest(id int) {
	resp, err := http.Get("http://localhost:8080/api/photo/list")
	if err != nil {
		fmt.Printf("Request %d lỗi: %v\n", id, err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Request %d trả về status: %s\n", id, resp.Status)
}

const (
	maxConcurrency = 20
	minDelay       = 400
	maxExtraDelay  = 200
)

func SpamRequests(w http.ResponseWriter, r *http.Request) {
	quantityStr := r.URL.Query().Get("quantity")
	num, err := strconv.Atoi(quantityStr)
	if err != nil || num <= 0 {
		// Trả về JSON với lỗi nếu quantity không hợp lệ
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Invalid Quantity"}`, http.StatusBadRequest)
		return
	}

	sem := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < num; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func(id int) {
			defer wg.Done()
			CallRequest(id)

			time.Sleep(time.Duration(minDelay+rand.Intn(maxExtraDelay)) * time.Millisecond)

			<-sem
		}(i)
	}

	wg.Wait()
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": fmt.Sprintf("Đã gửi xong %d request!", num),
	}
	json.NewEncoder(w).Encode(response)
}

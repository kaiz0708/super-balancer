package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	maxConcurrency = 30
	minDelay       = 600
	maxExtraDelay  = 400
)

func CallRequest(id int) {
	endpoints := []struct {
		URL    string
		Method string
		Body   string
	}{
		{"http://localhost:8080/api/photo/list", "GET", ""},
		{"http://localhost:8080/api/data", "POST", `{"example":"post data"}`},
		{"http://localhost:8080/api/data", "PUT", `{"example":"updated data"}`},
		{"http://localhost:8080/api/data", "DELETE", ""},
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(endpoints))
	selected := endpoints[randomIndex]

	var req *http.Request
	var err error

	if selected.Body != "" {
		req, err = http.NewRequest(selected.Method, selected.URL, bytes.NewBuffer([]byte(selected.Body)))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(selected.Method, selected.URL, nil)
	}

	if err != nil {
		fmt.Printf("Request %d tạo request lỗi: %v\n", id, err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Request %d lỗi: %v\n", id, err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Request %d (%s %s) -> Status: %d\n", id, selected.Method, selected.URL, resp.StatusCode)
}

func SpamRequests(w http.ResponseWriter, r *http.Request) {
	quantityStr := r.URL.Query().Get("quantity")
	num, err := strconv.Atoi(quantityStr)
	if err != nil || num <= 0 {
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

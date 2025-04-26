package test

import (
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

func SpamRequests(w http.ResponseWriter, r *http.Request) {
	quantityStr := r.URL.Query().Get("quantity")
	num, err := strconv.Atoi(quantityStr)
	if err != nil || num <= 0 {
		http.Error(w, "Số lượng không hợp lệ", http.StatusBadRequest)
		return
	}

	const maxConcurrency = 20
	sem := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < num; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func(id int) {
			defer wg.Done()
			CallRequest(id)

			time.Sleep(time.Duration(10+rand.Intn(20)) * time.Millisecond)

			<-sem
		}(i)
	}

	wg.Wait()
	fmt.Fprintln(w, "Đã gửi xong", num, "request!")
}

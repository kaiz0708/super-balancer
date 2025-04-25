package test

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

func CallRequest(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get("http://localhost:8080/api/photo/list")
	if err != nil {
		fmt.Printf("Request %d lỗi: %v\n", id, err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Request %d trả về status: %s\n", id, resp.Status)
}

func SpamRequests(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	quantity := r.URL.Query().Get("quantity")
	num, err := strconv.Atoi(quantity)
	wg.Add(num)

	if err != nil {
		fmt.Println("Lỗi chuyển đổi:", err)
	} else {
		fmt.Println("Kết quả:", num)
	}

	for i := 0; i < num; i++ {
		go CallRequest(i, &wg)
	}

	wg.Wait()
	fmt.Println("Gửi xong 1000 request!")
}

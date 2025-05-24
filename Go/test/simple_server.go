package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	port := "8081" // Có thể thay đổi port cho mỗi server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Thêm delay ngẫu nhiên để mô phỏng thời gian xử lý thực tế
		time.Sleep(time.Duration(100+time.Now().UnixNano()%900) * time.Millisecond)
		
		// Trả về thông tin server để dễ dàng phân biệt
		fmt.Fprintf(w, "Response from test server on port %s\n", port)
		fmt.Fprintf(w, "Request received at: %s\n", time.Now().Format(time.RFC3339))
	})

	fmt.Printf("Starting test server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
} 
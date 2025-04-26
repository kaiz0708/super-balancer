package test

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
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

func SpamRequests(c *gin.Context) {
	quantityStr := c.DefaultQuery("quantity", "0")
	num, err := strconv.Atoi(quantityStr)
	if err != nil || num <= 0 {
		c.JSON(400, gin.H{"error": "Số lượng không hợp lệ"})
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
	c.JSON(200, gin.H{"message": fmt.Sprintf("Đã gửi xong %d request!", num)})
}

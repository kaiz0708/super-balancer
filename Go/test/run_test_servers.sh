#!/bin/bash

# Compile test server
go build -o test_server simple_server.go

# Chạy 3 test server trên các port khác nhau
./test_server -port 8081 &
./test_server -port 8082 &
./test_server -port 8083 &

# Đợi các server khởi động
sleep 2

# Chạy load test
go run load_test.go

# Dọn dẹp
pkill -f test_server 
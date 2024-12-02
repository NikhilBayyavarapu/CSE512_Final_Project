package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type LoginRequest struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TransactionTest struct {
	SenderID int `json:"sender_id"`
}

var loginURLS = []string{
	"http://localhost:8080/login",
	"http://localhost:8081/login",
	"http://localhost:8082/login",
	"http://localhost:8083/login",
	"http://localhost:8084/login",
	"http://localhost:8085/login",
}

var transactionURLS = []string{
	"http://localhost:8080/transactions",
	"http://localhost:8081/transactions",
	"http://localhost:8082/transactions",
	"http://localhost:8083/transactions",
	"http://localhost:8084/transactions",
	"http://localhost:8085/transactions",
}

// sendRequest sends a POST request to one of the backend servers.
func sendLoginRequest(payload LoginRequest, results chan<- time.Duration, failureCount *int32, idx int) {
	// Serialize the payload
	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		atomic.AddInt32(failureCount, 1) // Increment failure count
		return
	}

	// Pick a server using a random strategy
	url := loginURLS[idx]

	// Record the start time
	start := time.Now()

	// Send the request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Request failed:", err)
		atomic.AddInt32(failureCount, 1) // Increment failure count
		return
	}
	defer resp.Body.Close()

	// Check for non-200 status codes as failures
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with status code: %d\n", resp.StatusCode)
		atomic.AddInt32(failureCount, 1) // Increment failure count
		return
	}

	// Record the duration
	duration := time.Since(start)

	// Send the duration to the results channel
	results <- duration
}

func getTransactions(payload TransactionTest, results chan<- time.Duration, failureCount *int32, idx int) {
	// Pick a server using a random strategy
	url := transactionURLS[idx]

	// Record the start time
	start := time.Now()

	// Send the request
	resp, err := http.Get(fmt.Sprintf("%s?sender_id=%d", url, payload.SenderID))
	if err != nil {
		fmt.Println("Request failed:", err)
		atomic.AddInt32(failureCount, 1) // Increment failure count
		return
	}
	defer resp.Body.Close()

	// Check for non-200 status codes as failures
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with status code: %d\n", resp.StatusCode)
		atomic.AddInt32(failureCount, 1) // Increment failure count
		return
	}

	// Record the duration
	duration := time.Since(start)

	// Send the duration to the results channel
	results <- duration
}

func main() {
	concurrentWorkers := 100 // Number of concurrent workers (requests sent simultaneously)
	totalRequests := 30000   // Total requests to send
	counter := 0

	payload := LoginRequest{
		UserID:   "106",
		Email:    "Joe.Wilderman@hotmail.com",
		Password: "r3h5_o0Z8K5lsQI",
	}

	var wg sync.WaitGroup
	var failureCount int32 // Counter for failed requests
	results := make(chan time.Duration, totalRequests)

	// Start workers
	startTest := time.Now()
	for i := 0; i < concurrentWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < totalRequests/concurrentWorkers; j++ {
				index := counter % 6 // Randomize server choice
				counter++
				sendLoginRequest(payload, results, &failureCount, index)
			}
		}(i)
	}

	// Wait for all workers to finish
	wg.Wait()
	close(results)

	// Calculate and display stats
	totalTime := time.Since(startTest)
	var totalDuration time.Duration
	for result := range results {
		totalDuration += result
	}
	avgDuration := totalDuration / time.Duration(totalRequests)

	// Print detailed results
	fmt.Println("Login Test")
	fmt.Printf("Total Login equests: %d\n", totalRequests)
	fmt.Printf("Concurrent requests: %d\n", concurrentWorkers)
	fmt.Printf("Total time taken: %v\n", totalTime)
	fmt.Printf("Average request time: %.6f ms\n", avgDuration.Seconds()*1000) // Milliseconds with microsecond precision
	fmt.Printf("Requests per second: %.2f\n", float64(totalRequests)/totalTime.Seconds())
	fmt.Printf("Failed requests: %d\n", failureCount)

	// Transaction test
	payload2 := TransactionTest{
		SenderID: 107,
	}

	var wg2 sync.WaitGroup
	var failureCount2 int32 // Counter for failed requests
	results2 := make(chan time.Duration, totalRequests)

	// Start workers
	startTest2 := time.Now()
	for i := 0; i < concurrentWorkers; i++ {
		wg2.Add(1)
		go func(workerID int) {
			defer wg2.Done()
			for j := 0; j < totalRequests/concurrentWorkers; j++ {
				index := counter % 6 // Randomize server choice
				counter++
				getTransactions(payload2, results2, &failureCount2, index)
			}
		}(i)
	}

	// Wait for all workers to finish
	wg2.Wait()
	close(results2)

	// Calculate and display stats
	totalTime2 := time.Since(startTest2)
	var totalDuration2 time.Duration
	for result := range results2 {
		totalDuration2 += result
	}

	avgDuration2 := totalDuration2 / time.Duration(totalRequests)

	// Print detailed results
	fmt.Println("\n\nTransaction Test")
	fmt.Printf("Total Get Transaction requests: %d\n", totalRequests)
	fmt.Printf("Concurrent requests: %d\n", concurrentWorkers)
	fmt.Printf("Total time taken: %v\n", totalTime2)
	fmt.Printf("Average request time: %.6f ms\n", avgDuration2.Seconds()*1000) // Milliseconds with microsecond precision
	fmt.Printf("Requests per second: %.2f\n", float64(totalRequests)/totalTime2.Seconds())
	fmt.Printf("Failed requests: %d\n", failureCount2)
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
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

var loginURLS = []string{
	"http://localhost:8080/login",
	"http://localhost:8081/login",
	"http://localhost:8082/login",
	"http://localhost:8083/login",
	"http://localhost:8084/login",
	"http://localhost:8085/login",
}

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

func testLogin() {
	concurrentWorkers := 100 // Number of concurrent workers (requests sent simultaneously)
	totalRequests := 1000    // Total requests to send
	counter := 0

	payloads := []LoginRequest{
		{
			UserID:   "100",
			Email:    "Patrick_Hackett31@gmail.com",
			Password: "WHeI1fEFjuDoi3o",
		},
		{
			UserID:   "101",
			Email:    "Cullen_Hilpert30@yahoo.com",
			Password: "UcPVrTS_Qw0MzZO",
		},
		{
			UserID:   "102",
			Email:    "Macy82@yahoo.com",
			Password: "kYpvO2tK8gNywwq",
		},
		{
			UserID:   "103",
			Email:    "Destinee.Hauck42@hotmail.com",
			Password: "ytw_ex3rJ4F2EzA",
		},
		{
			UserID:   "104",
			Email:    "Mathilde_Kertzmann38@hotmail.com",
			Password: "7JpH7EpQglV5EbL",
		},
		{
			UserID:   "105",
			Email:    "Ivory2@gmail.com",
			Password: "esNq3ZHTo2tFzwi",
		},
		{
			UserID:   "106",
			Email:    "Joe.Wilderman@hotmail.com",
			Password: "r3h5_o0Z8K5lsQI",
		},
		{
			UserID:   "107",
			Email:    "Javier_Weimann45@yahoo.com",
			Password: "1jS1F9Sf3vC80YS",
		},
		{
			UserID:   "108",
			Email:    "Damien72@gmail.com",
			Password: "6SSONH37u3xt2ut",
		},
		{
			UserID:   "109",
			Email:    "Cassandra.Kuhic@gmail.com",
			Password: "DaI6b8UntUmqKQf",
		},
		{
			UserID:   "110",
			Email:    "Thomas19@yahoo.com",
			Password: "qfKH89aXG9QFcOW",
		},
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
				index2 := counter % 11
				counter++
				sendLoginRequest(payloads[index2], results, &failureCount, index)
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

}

var monthlyURLs = []string{
	"http://localhost:8080/monthdata",
	"http://localhost:8081/monthdata",
	"http://localhost:8082/monthdata",
	"http://localhost:8083/monthdata",
	"http://localhost:8084/monthdata",
	"http://localhost:8085/monthdata",
}

func getMonthlyTransactions(userID int, month int, year int, results chan<- time.Duration, failureCount *int32, idx int) {
	// Select the server URL using the provided index
	url := monthlyURLs[idx]

	// Construct the request URL with the necessary query parameters
	requestURL := fmt.Sprintf("%s?user_id=%d&month=%d&year=%d", url, userID, month, year)

	// Record the start time for performance measurement
	start := time.Now()

	// Send the HTTP GET request
	resp, err := http.Get(requestURL)
	if err != nil {
		// Log the error and increment the failure counter
		atomic.AddInt32(failureCount, 1) // Increment failure count
		return
	}
	defer resp.Body.Close()

	// Check for non-200 status codes and treat them as failures
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request to %s failed with status code: %d\n", url, resp.StatusCode)
		atomic.AddInt32(failureCount, 1) // Increment failure count
		return
	}

	// Ignore the response body content
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		fmt.Printf("Failed to discard response body from %s: %v\n", url, err)
		atomic.AddInt32(failureCount, 1) // Increment failure count
		return
	}

	// Record the duration of the request
	duration := time.Since(start)

	// Send the duration to the results channel
	results <- duration
}

func testMonthlyTransactions() {
	concurrentWorkers := 100 // Number of concurrent workers (requests sent simultaneously)
	totalRequests := 1000    // Total requests to send
	counter := 0

	// Monthly transaction test
	var wg3 sync.WaitGroup
	var failureCount3 int32 // Counter for failed requests
	results3 := make(chan time.Duration, totalRequests)

	// Start workers
	startTest3 := time.Now()
	for i := 0; i < concurrentWorkers; i++ {
		wg3.Add(1)
		go func(workerID int) {
			defer wg3.Done()
			for j := 0; j < totalRequests/concurrentWorkers; j++ {
				index := counter % 6 // Randomize server choice
				counter++
				rand.Seed(time.Now().UnixNano())

				// Generate random values
				userID := rand.Intn(200000-100+1) + 100 // Between 100 and 200000 inclusive
				year := rand.Intn(2024-2020+1) + 2020   // Between 2020 and 2024 inclusive
				month := rand.Intn(12-1+1) + 1          // Between 1 and 12 inclusive
				getMonthlyTransactions(userID, month, year, results3, &failureCount3, index)
			}
		}(i)
	}

	// Wait for all workers to finish
	wg3.Wait()
	close(results3)

	// Calculate and display stats
	totalTime3 := time.Since(startTest3)
	var totalDuration3 time.Duration
	for result := range results3 {
		totalDuration3 += result
	}

	avgDuration3 := totalDuration3 / time.Duration(totalRequests)

	// Print detailed results
	fmt.Println("\n\nMonthly Transaction Test")
	fmt.Printf("Total Get Monthly Transaction requests: %d\n", totalRequests)
	fmt.Printf("Concurrent requests: %d\n", concurrentWorkers)
	fmt.Printf("Total time taken: %v\n", totalTime3)
	fmt.Printf("Average request time: %.6f ms\n", avgDuration3.Seconds()*1000) // Milliseconds with microsecond precision
	fmt.Printf("Requests per second: %.2f\n", float64(totalRequests)/totalTime3.Seconds())
	fmt.Printf("Failed requests: %d\n", failureCount3)
}

type TransactionTest struct {
	SenderID int `json:"sender_id"`
}

var transactionURLS = []string{
	"http://localhost:8080/transactions",
	"http://localhost:8081/transactions",
	"http://localhost:8082/transactions",
	"http://localhost:8083/transactions",
	"http://localhost:8084/transactions",
	"http://localhost:8085/transactions",
}

func getTransactions(payload TransactionTest, results chan<- time.Duration, failureCount *int32, idx int) {
	// Pick a server using a random strategy
	url := transactionURLS[idx]

	// Record the start time
	start := time.Now()

	// Send the request
	resp, err := http.Get(fmt.Sprintf("%s?sender_id=%d", url, payload.SenderID))
	if err != nil {
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

func testTransaction() {
	concurrentWorkers := 100 // Number of concurrent workers (requests sent simultaneously)
	totalRequests := 1000    // Total requests to send
	counter := 0

	senderid := rand.Intn(200000-100+1) + 100 // Between 100 and 200000 inclusive

	// Transaction test
	payload2 := TransactionTest{
		SenderID: senderid,
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

func main() {
	// testLogin()
	// testMonthlyTransactions()
	testTransaction()
}

package main

import (
	"cse512/handlers"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestGetTransaction(t *testing.T) {

	users := []int{106}

	results := []struct {
		status string
		amount int
	}{
		{
			status: "completed",
			amount: 1452,
		},
		{
			status: "completed",
			amount: -5969,
		},
	}

	for _, user := range users {

		res, err := http.Get(fmt.Sprintf("http://localhost:8080/transactions?sender_id=%d", user))
		if err != nil {
			t.Errorf("Request failed: %v", err)
		}

		defer res.Body.Close()

		var finalResponse []handlers.TransactionResponse

		err = json.NewDecoder(res.Body).Decode(&finalResponse)
		if err != nil {
			t.Errorf("Error decoding response: %v", err)
		}

		for idx, response := range finalResponse {
			if idx >= 2 {
				break
			}

			if response.Amount != results[idx].amount {
				t.Errorf("Expected amount %d, got %d", results[idx].amount, response.Amount)
			}

			if response.Status != results[idx].status {
				t.Errorf("Expected status %s, got %s", results[idx].status, response.Status)
			}

		}

	}
}

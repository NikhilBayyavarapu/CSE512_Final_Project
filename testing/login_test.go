package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
)

type LoginRequest struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TestUser struct {
	UserID   int    `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func TestLogin(t *testing.T) {
	users := []TestUser{
		{
			UserID:   106,
			Email:    "Joe.Wilderman@hotmail.com",
			Password: "r3h5_o0Z8K5lsQI",
		},
		{
			UserID:   110,
			Email:    "Thomas19@yahoo.com",
			Password: "qfKH89aXG9QFcOW",
		},
		{
			UserID:   3,
			Email:    "wronguser@gmail.com",
			Password: "password",
		},
	}

	results := []int{200, 200, 401}

	for idx, user := range users {
		payload := LoginRequest{
			UserID:   strconv.Itoa(user.UserID),
			Email:    user.Email,
			Password: user.Password,
		}

		data, err := json.Marshal(payload)
		if err != nil {
			t.Errorf("Error marshalling JSON: %v", err)
		}

		res, err := http.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer(data))
		if err != nil {
			t.Errorf("Request failed: %v", err)
		}

		defer res.Body.Close()

		if res.StatusCode != results[idx] {
			t.Errorf("Expected status code %d, got %d", results[idx], res.StatusCode)
		}

		if res.StatusCode == 200 {
			type Response struct {
				Status string `json:"status"`
			}

			var response Response

			err = json.NewDecoder(res.Body).Decode(&response)
			if err != nil {
				t.Errorf("Error decoding response: %v", err)
			}

			if response.Status != "success" {
				t.Errorf("Expected status 'success', got '%s'", response.Status)
			}
		}
	}
}

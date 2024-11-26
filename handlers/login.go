package handlers

import (
	"context"
	"cse512/db"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// Response structure for consistent frontend handling
type Response struct {
	Status  string `json:"status"`  // Status of the response ("success" or "error")
	Message string `json:"message"` // Detailed message for the response
	Data    any    `json:"data"`    // Optional field to return any additional data
}

// HandleLogin processes user login requests
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{
			Status:  "error",
			Message: "Invalid request method. Only POST is allowed.",
		})
		return
	}

	var credentials struct {
		UserID   string `json:"user_id"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Status:  "error",
			Message: "Failed to parse JSON.",
		})
		return
	}

	userID := credentials.UserID
	email := credentials.Email
	password := credentials.Password

	fmt.Println("Received login request for user:", userID, "with email:", email, "and password:", password)

	if userID == "" || email == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Status:  "error",
			Message: "Missing required fields (user_id, email, password).",
		})
		return
	}

	// Continue with database lookup and validation...
	client := db.GetClient()
	collection := client.Database("bank").Collection("users")

	// Fetch the user's hashed password from MongoDB
	var result bson.M
	user_id, _ := strconv.Atoi(userID)

	err = collection.FindOne(context.Background(), bson.M{"user_id": user_id}).Decode(&result)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{
			Status:  "error",
			Message: "Error fetching details. Please try again.",
		})
		return
	}

	// Validate email and hashed password
	storedPassword, _ := result["password"].(string)
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{
			Status:  "error",
			Message: "Invalid credentials. Please try again.",
		})
		return
	}

	storedEmail, _ := result["email"].(string)

	if storedEmail != email {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{
			Status:  "error",
			Message: "Invalid credentials. Please try again.",
		})
		return
	}

	// Successful login response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: "Login successful.",
		Data: map[string]any{
			"user_id": userID,
			"email":   email,
			"name":    result["first_name"].(string) + " " + result["last_name"].(string),
			"balance": result["current_balance"],
		},
	})
}

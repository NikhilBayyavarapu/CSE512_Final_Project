package main

import (
	"context"
	"cse512/handlers"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func authenticateUser(accountID int, inputPassword string) bool {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27151"))
	if err != nil {
		fmt.Printf("Failed to connect to MongoDB: %v\n", err)
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database("bank").Collection("users")
	// Fetch the user's hashed password from MongoDB
	var result bson.M
	err = collection.FindOne(ctx, bson.M{"user_id": accountID}).Decode(&result)
	if err != nil {
		fmt.Printf("Failed to find account: %v\n", err)
		return false
	}

	// Extract the hashed password from the database
	storedPassword, ok := result["password"].(string)
	if !ok {
		fmt.Println("Password not found in the database.")
		return false
	}

	// Compare the input password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(inputPassword))
	if err != nil {
		fmt.Println("Password does not match.")
		return false
	}

	fmt.Println("Authentication successful!")
	return true
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/login", handlers.HandleLogin).Methods("POST")

	fmt.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}

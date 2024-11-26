package handlers

import (
	"context"
	"cse512/db"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TransactionResponse represents the response structure for the transaction handler
type TransactionResponse struct {
	Status    string `json:"status"`
	Amount    int    `json:"amount"`
	TimeStamp int    `json:"dateTimeStamp"`
	Remarks   string `json:"remarks"`
}

// HandleTransaction handles requests for retrieving user transactions
func HandleTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(TransactionResponse{
			Status: "error",
			Amount: 0,
		})
		return
	}

	// Parse the request body
	var user struct {
		SenderID int `json:"sender_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(TransactionResponse{
			Status:  "error",
			Amount:  0,
			Remarks: "Failed to parse JSON.",
		})
		return
	}

	userID := user.SenderID

	// Database query setup
	client := db.GetClient()
	database := client.Database("bank")
	transactionCollection := database.Collection("transactions")

	// Find transactions where the user is either the sender or the receiver
	filter := bson.M{
		"$or": []bson.M{
			{"sender_id": userID},
			{"receiver_id": userID},
		},
	}
	opts := options.Find().
		SetSort(bson.D{{"dateTimeStamp", -1}}). // Sort by dateTimeStamp descending
		SetLimit(10)                            // Limit to 10 results

	// Execute the query
	cursor, err := transactionCollection.Find(context.Background(), filter, opts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TransactionResponse{
			Status:  "error",
			Remarks: "Failed to fetch transactions.",
		})
		return
	}
	defer cursor.Close(context.Background())

	// Decode results
	var transactions []TransactionResponse
	for cursor.Next(context.Background()) {
		var transaction struct {
			Status    string `bson:"status"`
			Amount    int    `bson:"amount"`
			TimeStamp int    `bson:"dateTimeStamp"`
			Remarks   string `bson:"remarks"`
		}
		if err := cursor.Decode(&transaction); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(TransactionResponse{
				Status:  "error",
				Remarks: "Failed to decode transactions.",
			})
			return
		}
		transactions = append(transactions, TransactionResponse{
			Status:    transaction.Status,
			Amount:    transaction.Amount,
			TimeStamp: transaction.TimeStamp,
			Remarks:   transaction.Remarks,
		})
	}

	// Return the results
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transactions)
}

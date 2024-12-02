package handlers

import (
	"context"
	"cse512/db"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type MonthlyTransaction struct {
	SenderID      int    `json:"sender_id"`
	ReceiverID    int    `json:"receiver_id"`
	Amount        int    `json:"amount"`
	Remarks       string `json:"remarks"`
	DateTimeStamp string `json:"dateTimeStamp"`
	Status        string `json:"status"`
}

func GetMonthData(w http.ResponseWriter, r *http.Request) {
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

	// Get user ID from query params
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	// Get month from query params
	month := r.URL.Query().Get("month")
	if month == "" {
		http.Error(w, "month is required", http.StatusBadRequest)
		return
	}

	// Get year from query params
	year := r.URL.Query().Get("year")
	if year == "" {
		http.Error(w, "year is required", http.StatusBadRequest)
		return
	}

	// Parse month and year into integers
	monthInt, err := strconv.Atoi(month)
	if err != nil || monthInt < 1 || monthInt > 12 {
		http.Error(w, "invalid month provided", http.StatusBadRequest)
		return
	}

	yearInt, err := strconv.Atoi(year)
	if err != nil || yearInt < 0 {
		http.Error(w, "invalid year provided", http.StatusBadRequest)
		return
	}

	// Calculate the start and end timestamps for the month
	startDate := time.Date(yearInt, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second) // Last second of the month

	startTimestamp := startDate.Unix()
	endTimestamp := endDate.Unix()

	// Query the MongoDB collection
	client := db.GetClient()
	db := client.Database("bank")
	collection := db.Collection("transactions") // Replace with actual DB and collection names

	user_id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "invalid user_id provided", http.StatusBadRequest)
		return
	}

	filter := bson.M{
		"$or": []bson.M{
			{"sender_id": user_id},
			{"receiver_id": user_id},
		},
		"dateTimeStamp": bson.M{
			"$gte": startTimestamp,
			"$lte": endTimestamp,
		},
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, fmt.Sprintf("error querying database: %v", err), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	// Parse the results into a slice
	var responses []MonthlyTransaction
	for cursor.Next(context.Background()) {
		var transaction struct {
			SenderID      int    `bson:"sender_id"`
			ReceiverID    int    `bson:"receiver_id"`
			Amount        int    `bson:"amount"`
			Remarks       string `bson:"remarks"`
			DateTimeStamp int    `bson:"dateTimeStamp"`
			Status        string `bson:"status"`
		}

		if err := cursor.Decode(&transaction); err != nil {
			http.Error(w, fmt.Sprintf("error decoding transaction: %v", err), http.StatusInternalServerError)
			return
		}

		// Convert timestamp to string format
		formattedDate := time.Unix(int64(transaction.DateTimeStamp), 0).Format("02 Jan 2006")

		// Create the response object
		responses = append(responses, MonthlyTransaction{
			SenderID:      transaction.SenderID,
			ReceiverID:    transaction.ReceiverID,
			Amount:        transaction.Amount,
			Remarks:       transaction.Remarks,
			DateTimeStamp: formattedDate,
			Status:        transaction.Status,
		})
	}

	// Check if no transactions were found
	if len(responses) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "No transactions found for the specified criteria"}`))
		return
	}

	// Set headers for file download
	w.Header().Set("Content-Disposition", "attachment; filename=transactions.csv")
	w.Header().Set("Content-Type", "text/csv")

	// Create a CSV writer
	writer := csv.NewWriter(w)

	// Write the header row
	err = writer.Write([]string{"Sender ID", "Receiver ID", "Amount", "Remarks", "DateTimeStamp", "Status"})
	if err != nil {
		http.Error(w, fmt.Sprintf("error writing CSV header: %v", err), http.StatusInternalServerError)
		return
	}

	// Write the data rows
	for _, transaction := range responses {
		err := writer.Write([]string{
			strconv.Itoa(transaction.SenderID),
			strconv.Itoa(transaction.ReceiverID),
			strconv.Itoa(transaction.Amount),
			transaction.Remarks,
			transaction.DateTimeStamp,
			transaction.Status,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("error writing CSV row: %v", err), http.StatusInternalServerError)
			return
		}
	}

	// Flush the writer to ensure all data is written
	writer.Flush()

	// Check if there were any errors during the write process
	if err := writer.Error(); err != nil {
		http.Error(w, fmt.Sprintf("error flushing CSV data: %v", err), http.StatusInternalServerError)
		return
	}
}

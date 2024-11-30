package handlers

import (
	"context"
	"cse512/datamodels"
	"cse512/db"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Response structure for sending API responses
type Transaction struct {
	Status         string `json:"status"`
	Message        string `json:"message"`
	UpdatedBalance int    `json:"updated_balance"`
}

// insertErrorTransaction inserts a failed transaction record into the database
func insertErrorTransaction(senderID, receiverID, amount int, remarks string, timestamp int64, status string) {
	client := db.GetClient()
	transactionsCollection := client.Database("bank").Collection("transactions")

	failedTransaction := datamodels.Transaction{
		SenderID:      senderID,
		ReceiverID:    receiverID,
		Amount:        amount,
		Remarks:       remarks,
		DateTimeStamp: timestamp,
		Status:        status,
	}

	transactionsCollection.InsertOne(context.Background(), failedTransaction)
	fmt.Println("Failed transaction logged.")
}

// PerformTransaction handles a transaction between sender and receiver (withdraw, deposit, or transfer)
func PerformTransaction(w http.ResponseWriter, r *http.Request) {
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
		json.NewEncoder(w).Encode(Transaction{
			Status:  "error",
			Message: "Invalid request method. Only POST is allowed.",
		})
		return
	}

	// Parse request body to get transaction details
	var transaction struct {
		SenderID      int    `json:"sender_id"`
		ReceiverID    int    `json:"receiver_id"`
		AccountNumber int    `json:"account_number"`
		Amount        int    `json:"amount"`
		Remarks       string `json:"remarks"`
		Timestamp     int64  `json:"dateTimeStamp"`
	}

	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Failed to parse JSON.",
		})
		return
	}

	senderID := transaction.SenderID
	receiverID := transaction.ReceiverID
	amount := transaction.Amount
	remarks := transaction.Remarks
	timestamp := transaction.Timestamp
	accountNumber := transaction.AccountNumber

	// Validate fields
	if amount == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transaction{
			Status:  "error",
			Message: "Amount is required.",
		})
		return
	}

	// Get MongoDB client and database
	client := db.GetClient()
	database := client.Database("bank")
	usersCollection := database.Collection("users")
	transactionsCollection := database.Collection("transactions")

	// Find sender's data including account number and balance
	var sender struct {
		Balance int `bson:"current_balance"`
	}

	err = usersCollection.FindOne(context.Background(), bson.M{"user_id": senderID}).Decode(&sender)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		if err == mongo.ErrNoDocuments {
			json.NewEncoder(w).Encode(Transaction{
				Status:         "error",
				Message:        "Sender not found.",
				UpdatedBalance: sender.Balance,
			})
			insertErrorTransaction(senderID, receiverID, amount, remarks, timestamp, "failed")
		} else {
			json.NewEncoder(w).Encode(Transaction{
				Status:         "error",
				Message:        "Failed to fetch sender's data.",
				UpdatedBalance: sender.Balance,
			})
			insertErrorTransaction(senderID, receiverID, amount, remarks, timestamp, "failed")
		}
		return
	}

	// Check if sender has enough balance for withdrawal
	if sender.Balance < amount && senderID != receiverID {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transaction{
			Status:         "error",
			Message:        "Insufficient balance.",
			UpdatedBalance: sender.Balance,
		})
		insertErrorTransaction(senderID, receiverID, amount, remarks, timestamp, "failed")
		return
	}

	// Find receiver's data including account number and balance
	var receiver struct {
		AccountNumber int `bson:"account_number"`
		Balance       int `bson:"current_balance"`
	}
	err = usersCollection.FindOne(context.Background(), bson.M{"user_id": receiverID}).Decode(&receiver)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		if err == mongo.ErrNoDocuments {
			json.NewEncoder(w).Encode(Transaction{
				Status:         "error",
				Message:        "Receiver not found.",
				UpdatedBalance: sender.Balance,
			})
			insertErrorTransaction(senderID, receiverID, amount, remarks, timestamp, "failed")
		} else {
			json.NewEncoder(w).Encode(Transaction{
				Status:         "error",
				Message:        "Failed to fetch receiver's data.",
				UpdatedBalance: sender.Balance,
			})
			insertErrorTransaction(senderID, receiverID, amount, remarks, timestamp, "failed")
		}
		return
	}

	// Check if receiver's account number matches
	if receiver.AccountNumber != accountNumber {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transaction{
			Status:  "error",
			Message: "Receiver's account number does not match.",
		})
		insertErrorTransaction(senderID, receiverID, amount, remarks, timestamp, "failed")
		return
	}

	// Start MongoDB session to ensure atomicity
	session, err := client.StartSession()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Transaction{
			Status:         "error",
			Message:        "Failed to start session.",
			UpdatedBalance: sender.Balance,
		})
		return
	}
	defer session.EndSession(context.Background())

	// Start transaction
	err = session.StartTransaction()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Transaction{
			Status:         "error",
			Message:        "Failed to start transaction.",
			UpdatedBalance: sender.Balance,
		})
		return
	}

	// Handle self transaction (withdrawal or deposit)
	if senderID == receiverID {
		if amount > 0 {
			// Deposit: Increase balance
			_, err := usersCollection.UpdateOne(
				context.Background(),
				bson.M{"user_id": receiverID},
				bson.M{"$inc": bson.M{"current_balance": amount}},
			)
			if err != nil {
				session.AbortTransaction(context.Background())
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(Transaction{
					Status:         "error",
					Message:        "Failed to update balance (deposit).",
					UpdatedBalance: sender.Balance,
				})
				insertErrorTransaction(senderID, receiverID, amount, remarks, timestamp, "failed")
				return
			}
		} else {
			// Withdrawal: Decrease balance
			_, err := usersCollection.UpdateOne(
				context.Background(),
				bson.M{"user_id": senderID},
				bson.M{"$inc": bson.M{"current_balance": amount}},
			)
			if err != nil {
				session.AbortTransaction(context.Background())
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(Transaction{
					Status:         "error",
					Message:        "Failed to update balance (withdrawal).",
					UpdatedBalance: sender.Balance,
				})
				insertErrorTransaction(senderID, receiverID, amount, remarks, timestamp, "failed")
				return
			}
		}
	} else {
		// Standard transfer: sender != receiver
		_, err := usersCollection.UpdateOne(context.Background(), bson.M{"user_id": senderID}, bson.M{"$inc": bson.M{"current_balance": -amount}})
		if err != nil {
			session.AbortTransaction(context.Background())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Transaction{
				Status:         "error",
				Message:        "Failed to update sender's balance.",
				UpdatedBalance: sender.Balance,
			})
			insertErrorTransaction(senderID, receiverID, amount, remarks, timestamp, "failed")
			return
		}

		_, err = usersCollection.UpdateOne(context.Background(), bson.M{"user_id": receiverID}, bson.M{"$inc": bson.M{"current_balance": amount}})
		if err != nil {
			session.AbortTransaction(context.Background())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Transaction{
				Status:         "error",
				Message:        "Failed to update receiver's balance.",
				UpdatedBalance: sender.Balance,
			})
			insertErrorTransaction(senderID, receiverID, amount, remarks, timestamp, "failed")
			return
		}
	}

	// Commit the transaction
	err = session.CommitTransaction(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Transaction{
			Status:         "error",
			Message:        "Failed to commit transaction.",
			UpdatedBalance: sender.Balance,
		})
		insertErrorTransaction(senderID, receiverID, amount, remarks, timestamp, "failed")
		return
	}

	// Log transaction in the transactions collection
	completedTransaction := datamodels.Transaction{
		SenderID:      senderID,
		ReceiverID:    receiverID,
		Amount:        amount,
		Remarks:       remarks,
		DateTimeStamp: timestamp,
		Status:        "success",
	}

	_, err = transactionsCollection.InsertOne(context.Background(), completedTransaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Transaction{
			Status:         "error",
			Message:        "Failed to log transaction.",
			UpdatedBalance: sender.Balance,
		})
		return
	}

	err = usersCollection.FindOne(context.Background(), bson.M{"user_id": senderID}).Decode(&sender)
	// Success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Transaction{
		Status:         "success",
		Message:        "Transaction completed successfully.",
		UpdatedBalance: sender.Balance,
	})
}

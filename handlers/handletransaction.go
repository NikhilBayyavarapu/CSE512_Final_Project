package handlers

import (
	"context"
	"cse512/db"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Response structure for sending API responses
type Transaction struct {
	Status  string `json:"status"`
	Message string `json:"message"`
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
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transaction{
			Status:  "error",
			Message: "Failed to parse JSON.",
		})
		return
	}

	senderID := transaction.SenderID
	receiverID := transaction.ReceiverID
	amount := transaction.Amount
	remarks := transaction.Remarks
	timestamp := transaction.Timestamp
	// accountNumber := transaction.AccountNumber

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
		// AccountNumber int `bson:"account_number"`
		Balance int `bson:"current_balance"`
	}

	err = usersCollection.FindOne(context.Background(), bson.M{"user_id": senderID}).Decode(&sender)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(Transaction{
				Status:  "error",
				Message: "Sender not found.",
			})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Transaction{
				Status:  "error",
				Message: "Failed to fetch sender's data.",
			})
		}
		return
	}

	// Check if sender's account number matches
	// if sender.AccountNumber != accountNumber {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	json.NewEncoder(w).Encode(Transaction{
	// 		Status:  "error",
	// 		Message: "Sender's account number does not match.",
	// 	})
	// 	return
	// }

	// Check if sender has enough balance for withdrawal
	if sender.Balance < amount && senderID != receiverID {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transaction{
			Status:  "error",
			Message: "Insufficient balance.",
		})
		return
	}

	// Find receiver's data including account number and balance
	var receiver struct {
		AccountNumber int `bson:"account_number"`
		Balance       int `bson:"current_balance"`
	}
	err = usersCollection.FindOne(context.Background(), bson.M{"user_id": receiverID}).Decode(&receiver)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(Transaction{
				Status:  "error",
				Message: "Receiver not found.",
			})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Transaction{
				Status:  "error",
				Message: "Failed to fetch receiver's data.",
			})
		}
		return
	}

	// Check if receiver's account number matches
	// if receiver.AccountNumber != accountNumber {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	json.NewEncoder(w).Encode(Transaction{
	// 		Status:  "error",
	// 		Message: "Receiver's account number does not match.",
	// 	})
	// 	return
	// }

	// Start MongoDB session to ensure atomicity
	session, err := client.StartSession()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Transaction{
			Status:  "error",
			Message: "Failed to start session.",
		})
		return
	}
	defer session.EndSession(context.Background())

	// Start transaction
	err = session.StartTransaction()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Transaction{
			Status:  "error",
			Message: "Failed to start transaction.",
		})
		return
	}

	// Handle self transaction (withdrawal or deposit)
	if senderID == receiverID {
		if amount > 0 {
			// Deposit: Increase balance
			receiverUpdate := bson.M{"$inc": bson.M{"current_balance": amount}}
			_, err := usersCollection.UpdateOne(
				context.Background(),
				bson.M{"user_id": receiverID},
				receiverUpdate,
			)
			if err != nil {
				session.AbortTransaction(context.Background())
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(Transaction{
					Status:  "error",
					Message: "Failed to update balance (deposit).",
				})
				return
			}
		} else {
			// Withdrawal: Decrease balance
			senderUpdate := bson.M{"$inc": bson.M{"current_balance": -amount}} // amount is negative for withdrawal
			_, err := usersCollection.UpdateOne(
				context.Background(),
				bson.M{"user_id": senderID},
				senderUpdate,
			)
			if err != nil {
				session.AbortTransaction(context.Background())
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(Transaction{
					Status:  "error",
					Message: "Failed to update balance (withdrawal).",
				})
				return
			}
		}
	} else {
		// Standard transfer: sender != receiver
		// Update sender's balance (decrease amount)
		senderUpdate := bson.M{"$inc": bson.M{"current_balance": -amount}}
		_, err := usersCollection.UpdateOne(context.Background(), bson.M{"user_id": senderID}, senderUpdate)
		if err != nil {
			session.AbortTransaction(context.Background())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Transaction{
				Status:  "error",
				Message: "Failed to update sender's balance.",
			})
			return
		}

		// Update receiver's balance (increase amount)
		receiverUpdate := bson.M{"$inc": bson.M{"current_balance": amount}}
		_, err = usersCollection.UpdateOne(context.Background(), bson.M{"user_id": receiverID}, receiverUpdate)
		if err != nil {
			session.AbortTransaction(context.Background())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Transaction{
				Status:  "error",
				Message: "Failed to update receiver's balance.",
			})
			return
		}
	}

	// Insert transaction record
	transactionRecord := bson.M{
		"transaction_id": 1, // You can increment this or use auto-generated ID
		"sender_id":      senderID,
		"receiver_id":    receiverID,
		"amount":         amount,
		"remarks":        remarks,
		"dateTimeStamp":  timestamp,
		"status":         "completed",
	}

	_, err = transactionsCollection.InsertOne(context.Background(), transactionRecord)
	if err != nil {
		session.AbortTransaction(context.Background())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Transaction{
			Status:  "error",
			Message: "Failed to insert transaction record.",
		})
		return
	}

	// Commit the transaction
	err = session.CommitTransaction(context.Background())
	if err != nil {
		session.AbortTransaction(context.Background())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Transaction{
			Status:  "error",
			Message: "Failed to commit transaction.",
		})
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Transaction{
		Status:  "success",
		Message: "Transaction completed successfully.",
	})
}

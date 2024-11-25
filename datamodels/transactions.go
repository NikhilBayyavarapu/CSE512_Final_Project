package datamodels

type Transaction struct {
	TransactionID int    `json:"transaction_id" bson:"transaction_id"` // Unique ID for the transaction
	SenderID      int    `json:"sender_id" bson:"sender_id"`           // ID of the sender
	Amount        int    `json:"amount" bson:"amount"`                 // Transaction amount, can be negative for withdrawal
	ReceiverID    int    `json:"receiver_id" bson:"receiver_id"`       // ID of the receiver
	Remarks       string `json:"remarks" bson:"remarks"`               // Description or notes about the transaction
	DateTimeStamp int64  `json:"dateTimeStamp" bson:"dateTimeStamp"`   // Timestamp for the transaction
	Status        string `json:"status" bson:"status"`                 // Status of the transaction, e.g., completed
}

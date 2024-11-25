package datamodels

type User struct {
	UserID    int    `json:"user_id" bson:"user_id"`                 // Unique ID for the user
	FirstName string `json:"first_name" bson:"first_name"`           // First name of the user
	LastName  string `json:"last_name" bson:"last_name"`             // Last name of the user
	Email     string `json:"email" bson:"email"`                     // Email of the user
	Balance   int    `json:"current_balance" bson:"current_balance"` // Current balance of the user
	PassHash  string `json:"password" bson:"password"`               // Hash of the user's password
}

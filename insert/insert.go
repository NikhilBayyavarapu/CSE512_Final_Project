package utilities

import (
	"context"
	"cse512/datamodels"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func HashedPassword(pwd string) (string, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPwd), nil
}

func InsertData() {
	uri := "mongodb://localhost:27151"

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Error while connecting to MongoDB: %v", err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalf("Error while disconnecting from MongoDB: %v", err)
		}
	}()

	fmt.Println("Connected to MongoDB Router!")

	db := client.Database("bank")
	usersCollection := db.Collection("users")
	transactionCollection := db.Collection("transactions")

	filePath := "mock_data_userInfo.json"

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error while opening file: %v", err)
	}
	defer file.Close()

	toBytes, _ := io.ReadAll(file)
	var users []datamodels.User

	if err = json.Unmarshal(toBytes, &users); err != nil {
		log.Fatalf("Error while unmarshalling data: %v", err)
	}

	for _, user := range users {
		user.PassHash, err = HashedPassword(user.PassHash)
		if err != nil {
			log.Fatalf("Error while hashing password: %v", err)
		}

		_, err = usersCollection.InsertOne(ctx, user)
		if err != nil {
			log.Fatalf("Error while inserting user data: %v", err)
		}
	}

	filePath2 := "mock_transactions.json"
	file2, err := os.Open(filePath2)
	if err != nil {
		log.Fatalf("Error while opening file: %v", err)
	}
	defer file2.Close()

	toBytes2, _ := io.ReadAll(file2)
	var transactions []datamodels.Transaction

	if err = json.Unmarshal(toBytes2, &transactions); err != nil {
		log.Fatalf("Error while unmarshalling data: %v", err)
	}

	for _, transaction := range transactions {
		_, err = transactionCollection.InsertOne(ctx, transaction)
		if err != nil {
			log.Fatalf("Error while inserting transaction data: %v", err)
		}
	}

	fmt.Println("Data inserted successfully!")
	fmt.Println("Total users: ", len(users))
	fmt.Println("Total transactions: ", len(transactions))
}

package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func connect() {
	temp, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27151,localhost:27152,localhost:27153"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v\n", err)
	}

	if err = temp.Ping(context.Background(), nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v\n", err)
	}

	client = temp
	return
}

func GetClient() *mongo.Client {
	if client == nil {
		connect()
	}
	return client
}

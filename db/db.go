package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client

func connect() {
	// Configure connection pool settings
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27151,localhost:27152,localhost:27153").
		SetMaxPoolSize(30000).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(5 * time.Minute).
		SetReadPreference(readpref.Secondary()).
		SetReadConcern(readconcern.Local())

	// Establish connection to the MongoDB server
	temp, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v\n", err)
	}

	// Ping to check connection
	if err = temp.Ping(context.Background(), nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v\n", err)
	}

	client = temp
}

func GetClient() *mongo.Client {
	if client == nil {
		connect()
	}
	return client
}

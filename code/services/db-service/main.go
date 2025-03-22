package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	defer client.Disconnect(context.Background())

	fmt.Println("Connected to the Threadit database server!")

	// TODO: Create a gRPC server and remove the loop below

	// Infinite loop to simulate listening
	for {
		// Simulate listening for gRPC requests
		fmt.Println("Listening for gRPC requests...")
		time.Sleep(10 * time.Second)
	}
}

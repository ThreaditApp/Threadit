package main

import (
	"context"
	server "db-service/src"
	"fmt"
	dbpd "gen/db-service/pb"
	"log"
	"net"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatalf("missing MONGO_URI env var")
	}

	// get env port
	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		log.Fatalf("missing SERVICE_PORT env var")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	// Database name in mongodb
	databaseName := "threadit"

	// Paths to JSON files inside the container
	threadsPath := "/dataset/threads.json"
	communitiesPath := "/dataset/communities.json"

	// Load data into MongoDB
	if err := loadThreads(client, databaseName, threadsPath); err != nil {
		log.Fatalf("Error loading threads: %v", err)
	}
	if err := loadCommunities(client, databaseName, communitiesPath); err != nil {
		log.Fatalf("Error loading communities: %v", err)
	}

	// start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	mongoDatabase := client.Database(databaseName)
	dbpd.RegisterDBServiceServer(grpcServer, &server.DBServer{
		Mongo: mongoDatabase,
	})

	log.Printf("gRPC server is listening on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

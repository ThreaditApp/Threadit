package main

import (
	"context"
	server "db-service/src"
	"fmt"
	dbpd "gen/db-service/pb"
	"log"
	"net"
	"os"
	"time"

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

	// retry connecting to mongo every 5 seconds
	for err != nil {
		time.Sleep(5 * time.Second)
		log.Println("Attempting to connect to MongoDB...")
		client, err = mongo.Connect(context.Background(), clientOptions)
	}

	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	// Load data into MongoDB
	mongoDatabaseName := "threadit"
	serviceAccountJsonPath := "/var/secret/gcp/gcs-key.json"
	var basePath string
	if _, err := os.Stat(serviceAccountJsonPath); err == nil {
		basePath = "gs://threadit-dataset"
	} else if os.IsNotExist(err) {
		basePath = "/dataset"
	} else {
		log.Fatalf("error checking for gcs-key.json: %v", err)
	}
	if err := loadThreadsDataset(client, mongoDatabaseName, basePath); err != nil {
		log.Fatalf("Error loading threads: %v", err)
	}
	if err := loadCommunitiesDataset(client, mongoDatabaseName, basePath); err != nil {
		log.Fatalf("Error loading communities: %v", err)
	}

	// start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(1024*1024*500), // 500MB
		grpc.MaxSendMsgSize(1024*1024*500), // 500MB
	)
	mongoDatabase := client.Database(mongoDatabaseName)
	dbpd.RegisterDBServiceServer(grpcServer, &server.DBServer{
		Mongo: mongoDatabase,
	})

	log.Printf("gRPC server is listening on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

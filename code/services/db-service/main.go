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
	mongoDatabaseName := "threadit"
	serviceAccountJsonPath := "/var/secret/gcp/gcs-key.json"
	if _, err := os.Stat(serviceAccountJsonPath); err == nil {
		loadDataSetsFromBucket(client, mongoDatabaseName)
	} else if os.IsNotExist(err) {
		loadDatasetsFromLocal(client, mongoDatabaseName)
	} else {
		log.Fatalf("error checking for gcs-key.json: %v", err)
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

// GCS Bucket Paths to JSON files inside the container
func loadDataSetsFromBucket(client *mongo.Client, mongoDatabaseName string) {
	basePath := "gs://threadit-dataset"

	// Load data into MongoDB
	if err := loadThreadsFromBucket(client, mongoDatabaseName, basePath); err != nil {
		log.Fatalf("Error loading threads: %v", err)
	}
	if err := loadCommuninitiesFromBucket(client, mongoDatabaseName, basePath); err != nil {
		log.Fatalf("Error loading communities: %v", err)
	}
}

func loadDatasetsFromLocal(client *mongo.Client, mongoDatabaseName string) {
	basePath := "/dataset"

	// Load data into MongoDB
	if err := loadThreadsFromLocal(client, mongoDatabaseName, basePath); err != nil {
		log.Fatalf("Error loading threads: %v", err)
	}
	if err := loadCommunitiesFromLocal(client, mongoDatabaseName, basePath); err != nil {
		log.Fatalf("Error loading communities: %v", err)
	}
}

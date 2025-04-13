package main

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Thread struct {
	ID          string `bson:"_id" json:"_id"`
	Title       string `bson:"title" json:"title"`
	Ups         int    `bson:"ups" json:"ups"`
	Downs       int    `bson:"downs" json:"downs"`
	Content     string `bson:"content" json:"content"`
	CommunityID string `bson:"community_id" json:"community_id"`
	NumComments int32  `bson:"num_comments" json:"num_comments"`
}

type Community struct {
	ID         string `bson:"_id" json:"_id"`
	Name       string `bson:"name" json:"name"`
	NumThreads int32  `bson:"num_threads" json:"num_threads"`
}

func loadThreadsDataset(client *mongo.Client, database, basePath string) error {
	var threads []Thread
	return loadDataset(client, database, "threads", basePath+"/threads.json", &threads)
}

func loadCommunitiesDataset(client *mongo.Client, database, basePath string) error {
	var communities []Community
	return loadDataset(client, database, "communities", basePath+"/communities.json", &communities)
}

// loads data from GCS bucket into MongoDB
func loadDataset[T any](client *mongo.Client, database, collectionName, jsonPath string, result *[]T) error {
	db := client.Database(database)
	collection := db.Collection(collectionName)

	// Check if collection already has data
	count, err := collection.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		return fmt.Errorf("failed to count documents in %s collection: %w", collectionName, err)
	}
	if count > 0 {
		log.Printf("%s collection already initialized, skipping JSON import.\n", collectionName)
		return nil
	}

	var data []byte

	if strings.HasPrefix(jsonPath, "gs://") {
		// Parse GCS path
		parts := strings.SplitN(jsonPath[5:], "/", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid GCS path: %s", jsonPath)
		}
		bucketName, objectName := parts[0], parts[1]

		ctx := context.Background()
		storageClient, err := storage.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("failed to create GCS client: %w", err)
		}
		defer storageClient.Close()

		reader, err := storageClient.Bucket(bucketName).Object(objectName).NewReader(ctx)
		if err != nil {
			return fmt.Errorf("failed to open GCS object: %w", err)
		}
		defer reader.Close()

		data, err = io.ReadAll(reader)
		if err != nil {
			return fmt.Errorf("failed to read from GCS object: %w", err)
		}
	} else {
		// Fallback to local file
		data, err = os.ReadFile(jsonPath)
		if err != nil {
			return fmt.Errorf("failed to read %s JSON file: %w", collectionName, err)
		}
	}

	// Unmarshal JSON
	if err := json.Unmarshal(data, result); err != nil {
		return fmt.Errorf("failed to unmarshal %s JSON: %w", collectionName, err)
	}

	// Convert to []interface{}
	var documents []interface{}
	for _, item := range *result {
		documents = append(documents, item)
	}

	if len(documents) > 0 {
		_, err = collection.InsertMany(context.TODO(), documents)
		if err != nil {
			return fmt.Errorf("failed to insert %s into MongoDB: %w", collectionName, err)
		}
		log.Printf("%s collection initialized with JSON data.\n", collectionName)
	} else {
		log.Printf("No data found in %s JSON file.\n", collectionName)
	}

	return nil
}

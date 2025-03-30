package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Thread struct {
	ID          string   `bson:"_id" json:"_id"`
	Title       string   `bson:"title" json:"title"`
	Ups         int      `bson:"ups" json:"ups"`
	Downs       int      `bson:"downs" json:"downs"`
	Content     string   `bson:"content" json:"content"`
	CommunityID string   `bson:"community_id" json:"community_id"`
	Comments    []string `bson:"comments" json:"comments"`
}

type Community struct {
	ID      string   `bson:"_id" json:"_id"`
	Name    string   `bson:"name" json:"name"`
	Threads []string `bson:"threads" json:"threads"`
}

// loadData is a helper function to load data from a JSON file into MongoDB.
func loadData[T any](client *mongo.Client, database, collectionName, jsonPath string, result *[]T) error {
	db := client.Database(database)
	collection := db.Collection(collectionName)

	// Check if the collection already contains documents
	count, err := collection.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		return fmt.Errorf("failed to count documents in %s collection: %w", collectionName, err)
	}
	if count > 0 {
		log.Printf("%s collection already initialized, skipping JSON import.\n", collectionName)
		return nil
	}

	// Read JSON file
	file, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read %s JSON file: %w", collectionName, err)
	}

	// Unmarshal JSON into the provided slice
	if err := json.Unmarshal(file, result); err != nil {
		return fmt.Errorf("failed to unmarshal %s JSON: %w", collectionName, err)
	}

	// Convert to a slice of interface{} for MongoDB insertion
	var documents []interface{}
	for _, item := range *result {
		documents = append(documents, item)
	}

	// Insert data into MongoDB if it's not empty
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

// loadThreads loads threads from a JSON file into MongoDB.
func loadThreads(client *mongo.Client, database, path string) error {
	var threads []Thread
	return loadData(client, database, "threads", path, &threads)
}

// loadCommunities loads communities from a JSON file into MongoDB.
func loadCommunities(client *mongo.Client, database, path string) error {
	var communities []Community
	return loadData(client, database, "communities", path, &communities)
}

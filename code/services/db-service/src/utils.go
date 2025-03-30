package server

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getFindOptions(offset int32, limit int32, sortBy string) *options.FindOptions {
	findOptions := options.Find()
	if offset >= 0 {
		findOptions.SetSkip(int64(offset))
	}
	if limit > 0 {
		findOptions.SetLimit(int64(limit))
	}
	if sortBy != "" {
		findOptions.SetSort(bson.D{{Key: sortBy, Value: -1}}) // sort in descending order
	}
	return findOptions
}

func generateUniqueId() string {
	id := primitive.NewObjectID()
	return id.Hex()
}

package server

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DefaultOffset int32 = 0
	DefaultLimit  int32 = 10
	MaxLimit      int32 = 100
)

func getFindOptions(offsetPtr *int32, limitPtr *int32, sortBy string) *options.FindOptions {
	findOptions := options.Find()
	offset := DefaultOffset
	if offsetPtr != nil {
		offset = *offsetPtr
	}
	limit := DefaultLimit
	if limitPtr != nil && *limitPtr < MaxLimit {
		limit = *limitPtr
	}
	if offset >= 0 {
		findOptions.SetSkip(int64(offset))
	}
	if limit > 0 {
		findOptions.SetLimit(int64(limit))
	}
	if sortBy != "" {
		findOptions.SetSort(bson.D{{Key: sortBy, Value: -1}})
	}

	return findOptions
}

func generateUniqueId() string {
	id := primitive.NewObjectID()
	return id.Hex()
}

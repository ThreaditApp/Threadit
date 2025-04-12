package server

import (
	"context"
	"errors"
	dbpb "gen/db-service/pb"
	models "gen/models/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *DBServer) ListCommunities(ctx context.Context, req *dbpb.ListCommunitiesRequest) (*dbpb.ListCommunitiesResponse, error) {
	collection := s.Mongo.Collection("communities")
	filter := bson.M{}
	if req.GetName() != "" {
		filter["name"] = bson.M{"$regex": req.GetName(), "$options": "i"} // case-insensitive name match
	}
	cursor, err := collection.Find(ctx, filter, getFindOptions(req.Offset, req.Limit, ""))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to find communities")
	}
	defer cursor.Close(ctx)

	var results []*models.Community
	for cursor.Next(ctx) {
		community := bson.M{}
		if err := cursor.Decode(&community); err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to list communities")
		}
		results = append(results, &models.Community{
			Id:         community["_id"].(string),
			Name:       community["name"].(string),
			NumThreads: community["num_threads"].(int32),
		})
	}
	if err := cursor.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "Unexpected cursor error")
	}

	return &dbpb.ListCommunitiesResponse{
		Communities: results,
	}, nil
}

func (s *DBServer) CreateCommunity(ctx context.Context, req *dbpb.CreateCommunityRequest) (*dbpb.CreateCommunityResponse, error) {
	collection := s.Mongo.Collection("communities")
	community := bson.M{
		"_id":         generateUniqueId(),
		"name":        req.GetName(),
		"num_threads": 0,
	}

	count, err := collection.CountDocuments(ctx, bson.M{"name": req.GetName()})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to check for name uniqueness")
	}
	if count > 0 {
		return nil, status.Errorf(codes.AlreadyExists, "Community name already in use")
	}
	_, err = collection.InsertOne(ctx, community)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create community")
	}

	return &dbpb.CreateCommunityResponse{
		Id: community["_id"].(string),
	}, nil
}

func (s *DBServer) GetCommunity(ctx context.Context, req *dbpb.GetCommunityRequest) (*models.Community, error) {
	collection := s.Mongo.Collection("communities")
	filter := bson.M{"_id": req.GetId()}

	var community bson.M
	err := collection.FindOne(ctx, filter).Decode(&community)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.NotFound, "Community not found")
		}
		return nil, status.Errorf(codes.Internal, "Failed to get community")
	}

	return &models.Community{
		Id:         community["_id"].(string),
		Name:       community["name"].(string),
		NumThreads: community["num_threads"].(int32),
	}, nil
}

func (s *DBServer) UpdateCommunity(ctx context.Context, req *dbpb.UpdateCommunityRequest) (*emptypb.Empty, error) {
	collection := s.Mongo.Collection("communities")
	setValues, incValues := bson.M{}, bson.M{}

	if req.Name != nil {
		filter := bson.M{"name": req.GetName(), "_id": bson.M{"$ne": req.GetId()}}
		count, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to check for name uniqueness")
		}
		if count > 0 {
			return nil, status.Errorf(codes.AlreadyExists, "Community name already in use")
		}
		setValues["name"] = req.GetName()
	}
	if req.NumThreadsOffset != nil {
		if offset := req.GetNumThreadsOffset(); offset == 1 {
			incValues["num_threads"] = 1
		} else {
			incValues["num_threads"] = -1
		}
	}
	if len(setValues) == 0 && len(incValues) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "No update fields provided")
	}

	update := bson.M{"$set": setValues, "$inc": incValues}
	result, err := collection.UpdateOne(ctx, bson.M{"_id": req.GetId()}, update)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update community")
	}
	if result.MatchedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "Community not found")
	}

	return &emptypb.Empty{}, nil
}

func (s *DBServer) DeleteCommunity(ctx context.Context, req *dbpb.DeleteCommunityRequest) (*emptypb.Empty, error) {
	collection := s.Mongo.Collection("communities")

	result, err := collection.DeleteOne(ctx, bson.M{"_id": req.GetId()})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete community")
	}
	if result.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "Community not found")
	}

	return &emptypb.Empty{}, nil
}

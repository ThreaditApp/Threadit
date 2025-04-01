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
	"log"
)

func (s *DBServer) ListCommunities(ctx context.Context, req *dbpb.ListCommunitiesRequest) (*dbpb.ListCommunitiesResponse, error) {
	collection := s.Mongo.Collection("communities")
	filter := bson.M{}
	if req.GetName() != "" {
		filter["name"] = bson.M{"$regex": req.GetName(), "$options": "i"} // case-insensitive name match
	}

	findOptions := getFindOptions(req.Offset, req.Limit, "")
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Printf("list communities: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to find communities")
	}
	defer cursor.Close(ctx)

	var results []*models.Community
	for cursor.Next(ctx) {
		var community struct {
			ID         string `bson:"_id"`
			Name       string `bson:"name"`
			NumThreads int32  `bson:"num_threads"`
		}
		if err := cursor.Decode(&community); err != nil {
			log.Printf("list communities: %v", err)
			return nil, status.Errorf(codes.Internal, "failed to list communities")
		}
		results = append(results, &models.Community{
			Id:         community.ID,
			Name:       community.Name,
			NumThreads: community.NumThreads,
		})
	}

	if err := cursor.Err(); err != nil {
		log.Printf("list communities: %v", err)
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	return &dbpb.ListCommunitiesResponse{
		Communities: results,
	}, nil
}

func (s *DBServer) CreateCommunity(ctx context.Context, req *dbpb.CreateCommunityRequest) (*dbpb.CreateCommunityResponse, error) {
	communityID := generateUniqueId()
	community := &models.Community{
		Id:         communityID,
		Name:       req.GetName(),
		NumThreads: 0,
	}
	doc := bson.M{
		"_id":         community.Id,
		"name":        community.Name,
		"num_threads": community.NumThreads,
	}

	collection := s.Mongo.Collection("communities")
	_, err := collection.InsertOne(ctx, doc)
	if err != nil {
		log.Printf("create community: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create community")
	}
	return &dbpb.CreateCommunityResponse{
		Id: communityID,
	}, nil
}

func (s *DBServer) GetCommunity(ctx context.Context, req *dbpb.GetCommunityRequest) (*models.Community, error) {
	collection := s.Mongo.Collection("communities")
	filter := bson.M{
		"_id": req.GetId(),
	}
	var community bson.M
	err := collection.FindOne(ctx, filter).Decode(&community)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.NotFound, "community not found")
		}
		log.Printf("get community: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get community")
	}
	return &models.Community{
		Id:         community["_id"].(string),
		Name:       community["name"].(string),
		NumThreads: community["num_threads"].(int32),
	}, nil
}

func (s *DBServer) UpdateCommunity(ctx context.Context, req *dbpb.UpdateCommunityRequest) (*emptypb.Empty, error) {
	collection := s.Mongo.Collection("communities")
	update := bson.M{}
	setFields := bson.M{}
	incFields := bson.M{}

	if req.Name != nil {
		setFields["name"] = req.GetName()
	}
	if req.NumThreadsOffset != nil {
		offset := req.GetNumThreadsOffset()
		if offset == 1 {
			incFields["num_threads"] = 1
		} else if offset == -1 {
			incFields["num_threads"] = -1
		} else {
			return nil, status.Errorf(codes.InvalidArgument, "invalid num threads offset")
		}
	}
	if len(setFields) > 0 {
		update["$set"] = setFields
	}
	if len(incFields) > 0 {
		update["$inc"] = incFields
	}
	if len(update) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "no update fields provided")
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": req.GetId()}, update)
	if err != nil {
		log.Printf("update community: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update community")
	}
	if result.MatchedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "community not found")
	}
	return &emptypb.Empty{}, nil
}

func (s *DBServer) DeleteCommunity(ctx context.Context, req *dbpb.DeleteCommunityRequest) (*emptypb.Empty, error) {
	collection := s.Mongo.Collection("communities")
	result, err := collection.DeleteOne(ctx, bson.M{"_id": req.GetId()})
	if err != nil {
		log.Printf("delete community: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete community")
	}
	if result.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "community not found")
	}
	return &emptypb.Empty{}, nil
}

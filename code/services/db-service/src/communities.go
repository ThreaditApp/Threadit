package server

import (
	"context"
	dbpb "gen/db-service/pb"
	models "gen/models/pb"

	"go.mongodb.org/mongo-driver/bson"
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

	findOptions := getFindOptions(req.GetLimit(), req.GetOffset(), "")
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find communities: %v", err)
	}
	defer cursor.Close(ctx)

	var results []*models.Community
	for cursor.Next(ctx) {
		var community struct {
			ID      string   `bson:"_id"`
			Name    string   `bson:"name"`
			Threads []string `bson:"threads"`
		}
		if err := cursor.Decode(&community); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to decode community: %v", err)
		}
		results = append(results, &models.Community{
			Id:      community.ID,
			Name:    community.Name,
			Threads: community.Threads,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "cursor error: %v", err)
	}
	return &dbpb.ListCommunitiesResponse{
		Communities: results,
	}, nil
}

func (s *DBServer) CreateCommunity(ctx context.Context, req *dbpb.CreateCommunityRequest) (*dbpb.CreateCommunityResponse, error) {
	communityID := generateUniqueId()
	community := &models.Community{
		Id:      communityID,
		Name:    req.GetName(),
		Threads: []string{},
	}
	doc := bson.M{
		"_id":     community.Id,
		"name":    community.Name,
		"threads": community.Threads,
	}

	collection := s.Mongo.Collection("communities")
	_, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create community: %v", err)
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
		return nil, status.Errorf(codes.NotFound, "community not found: %v", err)
	}
	return &models.Community{
		Id:   community["_id"].(string),
		Name: community["name"].(string),
	}, nil
}

func (s *DBServer) UpdateCommunity(ctx context.Context, req *dbpb.UpdateCommunityRequest) (*emptypb.Empty, error) {
	collection := s.Mongo.Collection("communities")
	filter := bson.M{"_id": req.GetId()}
	update := bson.M{"$set": bson.M{"name": req.GetName()}}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update community: %v", err)
	}
	if result.MatchedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "community not found")
	}
	return &emptypb.Empty{}, nil
}

func (s *DBServer) DeleteCommunity(ctx context.Context, req *dbpb.DeleteCommunityRequest) (*emptypb.Empty, error) {
	collection := s.Mongo.Collection("communities")

	// attempt deletion
	result, err := collection.DeleteOne(ctx, bson.M{"_id": req.GetId()})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete community: %v", err)
	}
	if result.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "community not found")
	}
	return &emptypb.Empty{}, nil
}

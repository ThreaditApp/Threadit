package server

import (
	"context"
	pb "db-service/src/pb"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DBServer struct {
	pb.UnimplementedDBServiceServer
	Mongo mongo.Client
}

func (s *DBServer) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	collection := s.Mongo.Database("mongo-database").Collection("users")
	user := bson.M{
		"username": in.GetUsername(),
		"email":    in.GetEmail(),
		"bio":      in.GetBio(),
	}

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *DBServer) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	collection := s.Mongo.Database("mongo-database").Collection("users")
	filter := bson.M{
		"username": in.GetUsername(),
		"email":    in.GetEmail(),
		"bio":      in.GetBio(),
	}
	var user bson.M
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{
		Username: user["username"].(string),
		Email:    user["email"].(string),
		Bio:      user["bio"].(string),
	}, nil
}

func (s *DBServer) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*emptypb.Empty, error) {
	collection := s.Mongo.Database("mongo-database").Collection("users")
	filter := bson.M{
		"id": in.GetId(),
	}

	// TODO: update only non-empty fields
	update := bson.M{
		"$set": bson.M{
			"username": in.GetUsername(),
			"email":    in.GetEmail(),
			"bio":      in.GetBio(),
		},
	}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
func (s *DBServer) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	collection := s.Mongo.Database("mongo-database").Collection("users")
	filter := bson.M{
		"id": in.GetId(),
	}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Feed service
func (s *DBServer) GetUserFeed(ctx context.Context, req *pb.GetUserFeedRequest) (*pb.GetUserFeedResponse, error) {
	collection := s.Mongo.Database("mongo-database").Collection("feeds")

	filter := bson.M{}
	sortOptions := bson.D{}
	if req.Sort != "" {
		sortOptions = bson.D{{Key: req.Sort, Value: -1}}
	}

	page := int64(1)
	pageSize := int64(25)
	if req.Page > 0 {
		page = int64(req.Page)
	}
	if req.PageSize > 0 {
		pageSize = int64(req.PageSize)
	}

	skip := (page - 1) * pageSize

	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(pageSize).
		SetSort(sortOptions)

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer cursor.Close(ctx)

	var feeds []bson.M
	if err := cursor.All(ctx, &feeds); err != nil {
		return nil, fmt.Errorf("error decoding database results: %w", err)
	}

	totalItems, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error counting documents: %w", err)
	}

	totalPages := (totalItems + pageSize - 1) / pageSize

	posts := make([]*pb.Feed, len(feeds))
	for i, feed := range feeds {
		posts[i] = &pb.Feed{
			Id:          feed["_id"].(string),
			Type:        feed["type"].(string),
			CommunityId: feed["community_id"].(string),
			Title:       feed["title"].(string),
			Content:     feed["content"].(string),
			CreatedAt:   timestamppb.New(feed["created_at"].(primitive.DateTime).Time()),
			UpdatedAt:   timestamppb.New(feed["updated_at"].(primitive.DateTime).Time()),
		}
	}

	return &pb.GetUserFeedResponse{
		Posts: posts,
		Pagination: &pb.Pagination{
			CurrentPage: int32(page),
			PerPage:     int32(pageSize),
			TotalItems:  int32(totalItems),
			TotalPages:  int32(totalPages),
		},
	}, nil
}

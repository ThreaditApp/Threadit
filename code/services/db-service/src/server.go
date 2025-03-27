package server

import (
	"context"
	pb "db-service/src/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DBServer struct {
	pb.UnimplementedDBServiceServer
	Mongo *mongo.Client
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
		"id:": in.GetId(),
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

// TODO: implement other methods for operations for communities, threads, comments and followers

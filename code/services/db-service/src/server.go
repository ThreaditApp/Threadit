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

// Comment service
func (s *DBServer) ListComments(ctx context.Context, in *pb.ListCommentsRequest) (*pb.ListCommentsResponse, error) {
	collection := s.Mongo.Database("mongo-database").Collection("comments")

	filter := bson.M{}
	if in.GetPostId() != "" {
		filter["post_id"] = in.GetPostId()
	}

	page := int64(1)
	pageSize := int64(25)
	if in.GetPage() > 0 {
		page = int64(in.GetPage())
	}
	if in.GetPageSize() > 0 {
		pageSize = int64(in.GetPageSize())
	}

	skip := (page - 1) * pageSize

	sortOptions := bson.D{bson.E{Key: "created_at", Value: -1}}

	totalItems, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(pageSize).
		SetSort(sortOptions)

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []*pb.Comment
	for cursor.Next(ctx) {
		var comment bson.M
		if err := cursor.Decode(&comment); err != nil {
			return nil, err
		}
		comments = append(comments, ConvertToProtoComment(comment))
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	totalPages := (totalItems + pageSize - 1) / pageSize

	return &pb.ListCommentsResponse{
		Comments: comments,
		Pagination: &pb.Pagination{
			CurrentPage: int32(page),
			PerPage:     int32(pageSize),
			TotalItems:  int32(totalItems),
			TotalPages:  int32(totalPages),
		},
	}, nil
}

func (s *DBServer) CreateComment(ctx context.Context, in *pb.CreateCommentRequest) (*pb.CreateCommentResponse, error) {
	collection := s.Mongo.Database("mongo-database").Collection("comments")

	if in.GetPostId() == "" || in.GetContent() == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	now := timestamppb.Now()

	comment := bson.M{
		"post_id":    in.GetPostId(),
		"user_id":    in.GetUserId(),
		"content":    in.GetContent(),
		"parent_id":  in.GetParentId(),
		"created_at": now.AsTime(),
	}

	res, err := collection.InsertOne(ctx, comment)
	if err != nil {
		return nil, err
	}

	objectID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to parse inserted ID")
	}

	return &pb.CreateCommentResponse{
		Comment: &pb.Comment{
			Id:        objectID.Hex(),
			PostId:    in.GetPostId(),
			UserId:    in.GetUserId(),
			Content:   in.GetContent(),
			ParentId:  in.GetParentId(),
			CreatedAt: now,
		},
	}, nil
}

func (s *DBServer) GetComment(ctx context.Context, in *pb.GetCommentRequest) (*pb.GetCommentResponse, error) {
	collection := s.Mongo.Database("mongo-database").Collection("comments")

	if in.GetId() == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	filter := bson.M{
		"_id": in.GetId(),
	}

	var comment bson.M
	err := collection.FindOne(ctx, filter).Decode(&comment)
	if err != nil {
		return nil, err
	}

	protoComment := ConvertToProtoComment(comment)

	return &pb.GetCommentResponse{
		Comment: protoComment,
	}, nil
}

func (s *DBServer) UpdateComment(ctx context.Context, in *pb.UpdateCommentRequest) (*pb.UpdateCommentResponse, error) {
	collection := s.Mongo.Database("mongo-database").Collection("comments")

	if in.GetId() == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	filter := bson.M{
		"_id": in.GetId(),
	}

	update := bson.M{
		"$set": bson.M{
			"content":    in.GetContent(),
			"updated_at": timestamppb.Now().AsTime(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedComment bson.M
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedComment)
	if err != nil {
		return nil, err
	}

	protoComment := ConvertToProtoComment(updatedComment)

	return &pb.UpdateCommentResponse{
		Comment: protoComment,
	}, nil
}

func (s *DBServer) DeleteComment(ctx context.Context, in *pb.DeleteCommentRequest) (*emptypb.Empty, error) {
	collection := s.Mongo.Database("mongo-database").Collection("comments")

	if in.GetId() == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	filter := bson.M{
		"_id": in.GetId(),
	}

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func ConvertToProtoComment(comment bson.M) *pb.Comment {
	createdAt := timestamppb.New(comment["created_at"].(primitive.DateTime).Time())

	return &pb.Comment{
		Id:        comment["_id"].(string),
		PostId:    comment["post_id"].(string),
		UserId:    comment["user_id"].(string),
		Content:   comment["content"].(string),
		ParentId:  comment["parent_id"].(string),
		CreatedAt: createdAt,
	}
}

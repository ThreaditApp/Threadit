package server

import (
	"context"
	pb "db-service/src/pb"
	"fmt"
	"strings"

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

func (s *DBServer) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
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

func (s *DBServer) ListThreads(ctx context.Context, in *pb.ListThreadsRequest) (*pb.ListThreadsResponse, error) {
	collection := s.Mongo.Database("mongo-database").Collection("threads")

	filter := bson.M{}
	if in.GetCommunityId() != "" {
		filter["community_id"] = in.GetCommunityId()
	}
	if in.GetAuthorId() != "" {
		filter["author_id"] = in.GetAuthorId()
	}
	if in.GetSearch() != "" {
		searchTerm := ".*" + in.GetSearch() + ".*"
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": searchTerm, "$options": "i"}},
			{"content": bson.M{"$regex": searchTerm, "$options": "i"}},
		}
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

	sortField := "created_at"
	sortOrder := int32(-1)
	if in.GetSortBy() != "" {
		sortBy := strings.ToLower(in.GetSortBy())
		if sortBy == "updated_at" { // allowed fields
			sortField = sortBy
		}
	}
	if in.GetSortOrder() != "" {
		inSortOrder := strings.ToLower(in.GetSortOrder())
		if inSortOrder == "asc" {
			sortOrder = 1
		}
	}

	sortOptions := bson.D{{Key: sortField, Value: sortOrder}}

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

	var threads []*pb.Thread
	for cursor.Next(ctx) {
		var thread bson.M
		if err := cursor.Decode(&thread); err != nil {
			return nil, err
		}
		threads = append(threads, ConvertToProtoThread(thread))
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	totalPages := (totalItems + pageSize - 1) / pageSize // ceiling division

	return &pb.ListThreadsResponse{
		Threads: threads,
		Pagination: &pb.Pagination{
			CurrentPage: int32(page),
			PerPage:     int32(pageSize),
			TotalItems:  int32(totalItems),
			TotalPages:  int32(totalPages),
		},
	}, nil
}

func (s *DBServer) CreateThread(ctx context.Context, in *pb.CreateThreadRequest) (*pb.Thread, error) {
	collection := s.Mongo.Database("mongo-database").Collection("threads")

	if in.GetCommunityId() == "" || in.GetAuthorId() == "" || in.GetTitle() == "" || in.GetContent() == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	// TODO: validate community_id

	user, err := s.GetUser(ctx, &pb.GetUserRequest{Id: in.GetAuthorId()})
	if err != nil || user == nil {
		return nil, fmt.Errorf("author not found")
	}

	now := timestamppb.Now()

	thread := bson.M{
		"community_id": in.GetCommunityId(),
		"author_id":    in.GetAuthorId(),
		"title":        in.GetTitle(),
		"content":      in.GetContent(),
		"created_at":   now.AsTime(),
		"updated_at":   now.AsTime(),
	}

	res, err := collection.InsertOne(ctx, thread)
	if err != nil {
		return nil, err
	}

	objectID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to parse inserted ID")
	}

	return &pb.Thread{
		Id:          objectID.Hex(),
		CommunityId: in.GetCommunityId(),
		AuthorId:    in.GetAuthorId(),
		Title:       in.GetTitle(),
		Content:     in.GetContent(),
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (s *DBServer) GetThread(ctx context.Context, in *pb.GetThreadRequest) (*pb.Thread, error) {
	collection := s.Mongo.Database("mongo-database").Collection("threads")

	if in.Id == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	filter := bson.M{
		"_id": in.GetId(),
	}

	var thread bson.M
	err := collection.FindOne(ctx, filter).Decode(&thread)
	if err != nil {
		return nil, err
	}

	return ConvertToProtoThread(thread), nil
}

func (s *DBServer) UpdateThread(ctx context.Context, in *pb.UpdateThreadRequest) (*pb.Thread, error) {
	collection := s.Mongo.Database("mongo-database").Collection("threads")

	if in.GetId() == "" || in.GetAuthorId() == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	user, err := s.GetUser(ctx, &pb.GetUserRequest{Id: in.GetAuthorId()})
	if err != nil || user == nil {
		return nil, fmt.Errorf("author not found")
	}

	filter := bson.M{
		"_id":       in.GetId(),
		"author_id": in.GetAuthorId(),
	}

	update := bson.M{
		"$set": bson.M{
			"updated_at": timestamppb.Now().AsTime(),
		},
	}

	updateFields := bson.M{}
	if in.GetTitle() != "" {
		updateFields["title"] = in.GetTitle()
	}
	if in.GetContent() != "" {
		updateFields["content"] = in.GetContent()
	}

	if len(updateFields) > 0 {
		for k, v := range updateFields {
			update["$set"].(bson.M)[k] = v
		}
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedThread bson.M
	err = collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedThread)
	if err != nil {
		return nil, err
	}

	return ConvertToProtoThread(updatedThread), nil
}

func (s *DBServer) DeleteThread(ctx context.Context, in *pb.DeleteThreadRequest) (*emptypb.Empty, error) {
	collection := s.Mongo.Database("mongo-database").Collection("threads")

	if in.GetId() == "" || in.GetAuthorId() == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	filter := bson.M{
		"_id":       in.GetId(),
		"author_id": in.GetAuthorId(),
	}

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func ConvertToProtoThread(thread bson.M) *pb.Thread {
	createdAt := timestamppb.New(thread["created_at"].(primitive.DateTime).Time())
	updatedAt := timestamppb.New(thread["updated_at"].(primitive.DateTime).Time())

	return &pb.Thread{
		Id:          thread["_id"].(string),
		CommunityId: thread["community_id"].(string),
		AuthorId:    thread["author_id"].(string),
		Title:       thread["title"].(string),
		Content:     thread["content"].(string),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

// TODO: implement other methods for operations for communities, comments and followers

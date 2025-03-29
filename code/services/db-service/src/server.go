package server

import (
	"context"
	"fmt"
	"gen/db-service/pb"
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

func (s *DBServer) ListCommunities(ctx context.Context, in *pb.ListCommunitiesRequest) (*pb.ListCommunitiesResponse, error) {
	collection := s.Mongo.Database("mongo-database").Collection("communities")

	filter := bson.M{}
	if in.GetOwnerId() != "" {
		ownerId, err := primitive.ObjectIDFromHex(in.GetOwnerId())
		if err != nil {
			return nil, fmt.Errorf("invalid owner ID: %v", err)
		}
		filter["owner_id"] = ownerId
	}
	if in.GetName() != "" {
		searchTerm := ".*" + in.GetName() + ".*"
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": searchTerm, "$options": "i"}},
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

	var communities []*pb.Community
	for cursor.Next(ctx) {
		var community bson.M
		if err := cursor.Decode(&community); err != nil {
			return nil, err
		}
		communities = append(communities, ConvertToProtoCommunity(community))
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	totalPages := (totalItems + pageSize - 1) / pageSize // ceiling division

	return &pb.ListCommunitiesResponse{
		Communities: communities,
		Pagination: &pb.Pagination{
			CurrentPage: int32(page),
			PerPage:     int32(pageSize),
			TotalItems:  int32(totalItems),
			TotalPages:  int32(totalPages),
		},
	}, nil
}

func (s *DBServer) CreateCommunity(ctx context.Context, in *pb.CreateCommunityRequest) (*pb.Community, error) {
	collection := s.Mongo.Database("mongo-database").Collection("communities")

	if in.GetOwnerId() == "" || in.GetName() == "" || in.GetDescription() == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	ownerId, err := primitive.ObjectIDFromHex(in.GetOwnerId())
	if err != nil {
		return nil, fmt.Errorf("invalid owner ID: %v", err)
	}

	now := timestamppb.Now()

	community := bson.M{
		"owner_id":    ownerId,
		"name":        in.GetName(),
		"description": in.GetDescription(),
		"created_at":  now.AsTime(),
		"updated_at":  now.AsTime(),
	}

	res, err := collection.InsertOne(ctx, community)
	if err != nil {
		return nil, err
	}

	objectID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to parse inserted ID")
	}

	return &pb.Community{
		Id:          objectID.Hex(),
		OwnerId:     in.GetOwnerId(),
		Name:        in.GetName(),
		Description: in.GetDescription(),
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (s *DBServer) GetCommunity(ctx context.Context, in *pb.GetCommunityRequest) (*pb.Community, error) {
	collection := s.Mongo.Database("mongo-database").Collection("communities")

	if in.GetId() == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid community ID: %v", err)
	}
	filter := bson.M{
		"_id": id,
	}

	var community bson.M
	err = collection.FindOne(ctx, filter).Decode(&community)
	if err != nil {
		return nil, err
	}

	return ConvertToProtoCommunity(community), nil
}

func (s *DBServer) UpdateCommunity(ctx context.Context, in *pb.UpdateCommunityRequest) (*pb.Community, error) {
	collection := s.Mongo.Database("mongo-database").Collection("communities")

	if in.GetId() == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid community ID: %v", err)
	}
	filter := bson.M{
		"_id": id,
	}

	update := bson.M{
		"$set": bson.M{
			"updated_at": timestamppb.Now().AsTime(),
		},
	}

	updateFields := bson.M{}
	if in.GetName() != "" {
		updateFields["name"] = in.GetName()
	}
	if in.GetDescription() != "" {
		updateFields["content"] = in.GetDescription()
	}

	if len(updateFields) > 0 {
		for k, v := range updateFields {
			update["$set"].(bson.M)[k] = v
		}
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedCommunity bson.M
	err = collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedCommunity)
	if err != nil {
		return nil, err
	}

	return ConvertToProtoCommunity(updatedCommunity), nil
}

func (s *DBServer) DeleteCommunity(ctx context.Context, in *pb.DeleteCommunityRequest) (*emptypb.Empty, error) {
	collection := s.Mongo.Database("mongo-database").Collection("communities")

	if in.GetId() == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid community ID: %v", err)
	}
	filter := bson.M{
		"_id": id,
	}

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func ConvertToProtoCommunity(community bson.M) *pb.Community {
	createdAt := timestamppb.New(community["created_at"].(primitive.DateTime).Time())
	updatedAt := timestamppb.New(community["updated_at"].(primitive.DateTime).Time())

	return &pb.Community{
		Id:          community["_id"].(primitive.ObjectID).Hex(),
		OwnerId:     community["owner_id"].(primitive.ObjectID).Hex(),
		Name:        community["name"].(string),
		Description: community["description"].(string),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func (s *DBServer) ListThreads(ctx context.Context, in *pb.ListThreadsRequest) (*pb.ListThreadsResponse, error) {
	collection := s.Mongo.Database("mongo-database").Collection("threads")

	filter := bson.M{}
	if in.GetCommunityId() != "" {
		communityId, err := primitive.ObjectIDFromHex(in.GetCommunityId())
		if err != nil {
			return nil, fmt.Errorf("invalid community ID: %v", err)
		}
		filter["community_id"] = communityId
	}
	if in.GetAuthorId() != "" {
		authorId, err := primitive.ObjectIDFromHex(in.GetAuthorId())
		if err != nil {
			return nil, fmt.Errorf("invalid author ID: %v", err)
		}
		filter["author_id"] = authorId
	}
	if in.GetTitle() != "" {
		searchTerm := ".*" + in.GetTitle() + ".*"
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": searchTerm, "$options": "i"}},
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

	communityId, err := primitive.ObjectIDFromHex(in.GetCommunityId())
	if err != nil {
		return nil, fmt.Errorf("invalid community ID: %v", err)
	}

	authorId, err := primitive.ObjectIDFromHex(in.GetAuthorId())
	if err != nil {
		return nil, fmt.Errorf("invalid author ID: %v", err)
	}

	now := timestamppb.Now()

	thread := bson.M{
		"community_id": communityId,
		"author_id":    authorId,
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

	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid thread ID: %v", err)
	}

	filter := bson.M{
		"_id": id,
	}

	var thread bson.M
	err = collection.FindOne(ctx, filter).Decode(&thread)
	if err != nil {
		return nil, err
	}

	return ConvertToProtoThread(thread), nil
}

func (s *DBServer) UpdateThread(ctx context.Context, in *pb.UpdateThreadRequest) (*pb.Thread, error) {
	collection := s.Mongo.Database("mongo-database").Collection("threads")

	if in.GetId() == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid thread ID: %v", err)
	}

	filter := bson.M{
		"_id": id,
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

	if in.GetId() == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid thread ID: %v", err)
	}

	filter := bson.M{
		"_id": id,
	}

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func ConvertToProtoThread(thread bson.M) *pb.Thread {
	createdAt := timestamppb.New(thread["created_at"].(primitive.DateTime).Time())
	updatedAt := timestamppb.New(thread["updated_at"].(primitive.DateTime).Time())

	return &pb.Thread{
		Id:          thread["_id"].(primitive.ObjectID).Hex(),
		CommunityId: thread["community_id"].(primitive.ObjectID).Hex(),
		AuthorId:    thread["author_id"].(primitive.ObjectID).Hex(),
		Title:       thread["title"].(string),
		Content:     thread["content"].(string),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

// TODO: implement other methods

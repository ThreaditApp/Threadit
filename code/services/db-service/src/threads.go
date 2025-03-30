package server

import (
	"context"
	"fmt"
	dbpb "gen/db-service/pb"
	models "gen/models/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

func (s *DBServer) ListThreads(ctx context.Context, req *dbpb.ListThreadsRequest) (*dbpb.ListThreadsResponse, error) {
	collection := s.Mongo.Collection("threads")

	sortBy := strings.Replace(req.GetSortBy(), "votes", "ups", 1)
	findOptions := getFindOptions(req.GetLimit(), req.GetOffset(), sortBy)
	filter := bson.M{}
	if req.GetCommunityId() != "" {
		communityId, err := primitive.ObjectIDFromHex(req.GetCommunityId())
		if err != nil {
			return nil, fmt.Errorf("invalid community ID: %v", err)
		}
		filter["community_id"] = communityId
	}
	if req.GetTitle() != "" {
		filter["title"] = bson.M{"$regex": req.GetTitle(), "$options": "i"} // case-insensitive title match
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find threads: %v", err)
	}
	defer cursor.Close(ctx)
	var results []*models.Thread
	for cursor.Next(ctx) {
		var thread struct {
			ID          string `bson:"_id"`
			CommunityId string `bson:"community_id"`
			Title       string `bson:"title"`
			Content     string `bson:"content"`
			Ups         int32  `bson:"ups"`
			Downs       int32  `bson:"downs"`
		}
		if err := cursor.Decode(&thread); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to decode thread: %v", err)
		}
		results = append(results, &models.Thread{
			Id:          thread.ID,
			CommunityId: thread.CommunityId,
			Title:       thread.Title,
			Content:     thread.Content,
			Ups:         thread.Ups,
			Downs:       thread.Downs,
		})
	}
	if err := cursor.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "cursor error: %v", err)
	}
	return &dbpb.ListThreadsResponse{
		Threads: results,
	}, nil
}

func (s *DBServer) CreateThread(ctx context.Context, req *dbpb.CreateThreadRequest) (*dbpb.CreateThreadResponse, error) {
	collection := s.Mongo.Collection("threads")

	// create thread
	threadID := generateUniqueId()
	thread := &models.Thread{
		Id:          threadID,
		CommunityId: req.GetCommunityId(),
		Title:       req.GetTitle(),
		Content:     req.GetContent(),
	}
	doc := bson.M{
		"_id":          thread.Id,
		"community_id": thread.CommunityId,
		"title":        thread.Title,
		"content":      thread.Content,
	}
	_, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create thread: %v", err)
	}

	// update the community with the new thread ID
	communityCollection := s.Mongo.Collection("communities")
	communityID, err := primitive.ObjectIDFromHex(thread.CommunityId)
	if err != nil {
		return nil, fmt.Errorf("invalid community ID: %v", err)
	}
	update := bson.M{
		"$addToSet": bson.M{
			"threads": thread.Id,
		},
	}
	_, err = communityCollection.UpdateOne(ctx, bson.M{"_id": communityID}, update)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update community with thread: %v", err)
	}
	return &dbpb.CreateThreadResponse{
		Id: threadID,
	}, nil
}

func (s *DBServer) GetThread(ctx context.Context, req *dbpb.GetThreadRequest) (*models.Thread, error) {
	collection := s.Mongo.Collection("threads")
	id, err := primitive.ObjectIDFromHex(req.GetId())
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
	return &models.Thread{
		Id:          thread["_id"].(primitive.ObjectID).Hex(),
		CommunityId: thread["community_id"].(primitive.ObjectID).Hex(),
		Title:       thread["title"].(string),
		Content:     thread["content"].(string),
		Ups:         thread["ups"].(int32),
		Downs:       thread["downs"].(int32),
	}, nil
}

func (s *DBServer) UpdateThread(ctx context.Context, req *dbpb.UpdateThreadRequest) (*emptypb.Empty, error) {
	collection := s.Mongo.Collection("threads")
	update := bson.M{}
	setFields := bson.M{}
	incFields := bson.M{}

	if req.Title != nil {
		setFields["title"] = req.GetTitle()
	}
	if req.Content != nil {
		setFields["content"] = req.GetContent()
	}
	if req.VoteOffset != nil {
		offset := req.GetVoteOffset()
		if offset > 0 {
			incFields["ups"] = offset
		} else if offset < 0 {
			incFields["downs"] = -offset // ensure it's positive
		} else {
			return nil, status.Errorf(codes.InvalidArgument, "vote offset cannot be zero")
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
		return nil, status.Errorf(codes.Internal, "failed to update thread: %v", err)
	}
	if result.MatchedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "thread not found")
	}
	return &emptypb.Empty{}, nil
}

func (s *DBServer) DeleteThread(ctx context.Context, req *dbpb.DeleteThreadRequest) (*emptypb.Empty, error) {
	collection := s.Mongo.Collection("threads")

	// get thread
	threadRes, err := s.GetThread(ctx, &dbpb.GetThreadRequest{Id: req.GetId()})

	// attempt deletion
	result, err := collection.DeleteOne(ctx, bson.M{"_id": req.GetId()})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete thread: %v", err)
	}
	if result.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "thread not found")
	}

	// remove thread id from the community
	communityCollection := s.Mongo.Collection("communities")
	communityID, err := primitive.ObjectIDFromHex(threadRes.GetCommunityId())
	if err != nil {
		return nil, fmt.Errorf("invalid community ID: %v", err)
	}
	update := bson.M{
		"$pull": bson.M{
			"threads": req.GetId(),
		},
	}
	_, err = communityCollection.UpdateOne(ctx, bson.M{"_id": communityID}, update)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update community with thread: %v", err)
	}

	return &emptypb.Empty{}, nil
}

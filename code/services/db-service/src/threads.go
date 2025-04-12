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

func (s *DBServer) ListThreads(ctx context.Context, req *dbpb.ListThreadsRequest) (*dbpb.ListThreadsResponse, error) {
	collection := s.Mongo.Collection("threads")
	filter := bson.M{}
	if id := req.GetCommunityId(); id != "" {
		filter["community_id"] = id
	}
	if title := req.GetTitle(); title != "" {
		filter["title"] = bson.M{"$regex": title, "$options": "i"} // case-insensitive title match
	}

	cursor, err := collection.Find(ctx, filter, getFindOptions(req.Offset, req.Limit, req.GetSortBy()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to find threads")
	}
	defer cursor.Close(ctx)

	var results []*models.Thread
	for cursor.Next(ctx) {
		thread := bson.M{}
		if err := cursor.Decode(&thread); err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to list threads")
		}
		results = append(results, &models.Thread{
			Id:          thread["_id"].(string),
			CommunityId: thread["community_id"].(string),
			Title:       thread["title"].(string),
			Content:     thread["content"].(string),
			Ups:         thread["ups"].(int32),
			Downs:       thread["downs"].(int32),
			NumComments: thread["num_comments"].(int32),
		})
	}
	if err := cursor.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "Unexpected cursor error")
	}

	return &dbpb.ListThreadsResponse{
		Threads: results,
	}, nil
}

func (s *DBServer) CreateThread(ctx context.Context, req *dbpb.CreateThreadRequest) (*dbpb.CreateThreadResponse, error) {
	collection := s.Mongo.Collection("threads")
	// create thread
	thread := bson.M{
		"_id":          generateUniqueId(),
		"community_id": req.GetCommunityId(),
		"title":        req.GetTitle(),
		"content":      req.GetContent(),
		"ups":          0,
		"downs":        0,
		"num_comments": 0,
	}
	_, err := collection.InsertOne(ctx, thread)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create thread")
	}
	return &dbpb.CreateThreadResponse{
		Id: thread["_id"].(string),
	}, nil
}

func (s *DBServer) GetThread(ctx context.Context, req *dbpb.GetThreadRequest) (*models.Thread, error) {
	collection := s.Mongo.Collection("threads")
	filter := bson.M{"_id": req.Id}

	var thread bson.M
	err := collection.FindOne(ctx, filter).Decode(&thread)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.NotFound, "Thread not found")
		}
		return nil, status.Errorf(codes.Internal, "Failed to get thread")
	}
	return &models.Thread{
		Id:          thread["_id"].(string),
		CommunityId: thread["community_id"].(string),
		Title:       thread["title"].(string),
		Content:     thread["content"].(string),
		Ups:         thread["ups"].(int32),
		Downs:       thread["downs"].(int32),
		NumComments: thread["num_comments"].(int32),
	}, nil
}

func (s *DBServer) UpdateThread(ctx context.Context, req *dbpb.UpdateThreadRequest) (*emptypb.Empty, error) {
	collection := s.Mongo.Collection("threads")
	setValues, incValues := bson.M{}, bson.M{}

	if req.Title != nil {
		setValues["title"] = req.GetTitle()
	}
	if req.Content != nil {
		setValues["content"] = req.GetContent()
	}
	if req.NumCommentsOffset != nil {
		if offset := req.GetNumCommentsOffset(); offset == 1 {
			incValues["num_comments"] = 1
		} else {
			incValues["num_comments"] = -1
		}
	}
	if req.VoteOffset != nil {
		if offset := req.GetVoteOffset(); offset == 1 {
			incValues["ups"] = 1
		} else {
			incValues["downs"] = 1
		}
	}
	if len(setValues) == 0 && len(incValues) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "No update fields provided")
	}

	update := bson.M{"$set": setValues, "$inc": incValues}
	result, err := collection.UpdateOne(ctx, bson.M{"_id": req.GetId()}, update)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update thread")
	}
	if result.MatchedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "Thread not found")
	}
	return &emptypb.Empty{}, nil
}

func (s *DBServer) DeleteThread(ctx context.Context, req *dbpb.DeleteThreadRequest) (*emptypb.Empty, error) {
	collection := s.Mongo.Collection("threads")

	result, err := collection.DeleteOne(ctx, bson.M{"_id": req.GetId()})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete thread")
	}
	if result.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "Thread not found")
	}

	return &emptypb.Empty{}, nil
}

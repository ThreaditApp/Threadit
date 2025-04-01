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

func (s *DBServer) ListThreads(ctx context.Context, req *dbpb.ListThreadsRequest) (*dbpb.ListThreadsResponse, error) {
	collection := s.Mongo.Collection("threads")
	findOptions := getFindOptions(req.Offset, req.Limit, *req.SortBy)
	filter := bson.M{}
	if req.GetCommunityId() != "" {
		filter["community_id"] = req.GetCommunityId()
	}
	if req.GetTitle() != "" {
		filter["title"] = bson.M{"$regex": req.GetTitle(), "$options": "i"} // case-insensitive title match
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Printf("list threads: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to find threads")
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
			NumComments int32  `bson:"num_comments"`
		}
		if err := cursor.Decode(&thread); err != nil {
			log.Printf("list threads: %v", err)
			return nil, status.Errorf(codes.Internal, "failed to list threads")
		}
		results = append(results, &models.Thread{
			Id:          thread.ID,
			CommunityId: thread.CommunityId,
			Title:       thread.Title,
			Content:     thread.Content,
			Ups:         thread.Ups,
			Downs:       thread.Downs,
			NumComments: thread.NumComments,
		})
	}
	if err := cursor.Err(); err != nil {
		log.Printf("list threads: %v", err)
		return nil, status.Errorf(codes.Internal, "internal error")
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
		Ups:         0,
		Downs:       0,
		NumComments: 0,
	}
	doc := bson.M{
		"_id":          thread.Id,
		"community_id": thread.CommunityId,
		"title":        thread.Title,
		"content":      thread.Content,
		"ups":          thread.Ups,
		"downs":        thread.Downs,
		"num_comments": thread.NumComments,
	}
	_, err := collection.InsertOne(ctx, doc)
	if err != nil {
		log.Printf("create thread: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create thread")
	}
	return &dbpb.CreateThreadResponse{
		Id: threadID,
	}, nil
}

func (s *DBServer) GetThread(ctx context.Context, req *dbpb.GetThreadRequest) (*models.Thread, error) {
	collection := s.Mongo.Collection("threads")
	filter := bson.M{
		"_id": req.Id,
	}
	var thread bson.M
	err := collection.FindOne(ctx, filter).Decode(&thread)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.NotFound, "thread not found")
		}
		log.Printf("get thread: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get thread")
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
	update := bson.M{}
	setFields := bson.M{}
	incFields := bson.M{}

	if req.Title != nil {
		setFields["title"] = req.GetTitle()
	}
	if req.Content != nil {
		setFields["content"] = req.GetContent()
	}
	if req.NumCommentsOffset != nil {
		offset := req.GetNumCommentsOffset()
		if offset == 1 {
			incFields["num_comments"] = 1
		} else if offset == -1 {
			incFields["num_comments"] = -1
		} else {
			return nil, status.Errorf(codes.InvalidArgument, "invalid num comments offset")
		}
	}
	if req.VoteOffset != nil {
		offset := req.GetVoteOffset()
		if offset == 1 {
			incFields["ups"] = 1
		} else if offset == -1 {
			incFields["downs"] = 1
		} else {
			return nil, status.Errorf(codes.InvalidArgument, "invalid vote offset")
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
		log.Printf("update thread: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update thread")
	}
	if result.MatchedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "thread not found")
	}
	return &emptypb.Empty{}, nil
}

func (s *DBServer) DeleteThread(ctx context.Context, req *dbpb.DeleteThreadRequest) (*emptypb.Empty, error) {
	collection := s.Mongo.Collection("threads")
	result, err := collection.DeleteOne(ctx, bson.M{"_id": req.GetId()})
	if err != nil {
		log.Printf("delete thread: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete thread")
	}
	if result.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "thread not found")
	}
	return &emptypb.Empty{}, nil
}

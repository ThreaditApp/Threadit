package server

import (
	"context"
	"fmt"
	dbpb "gen/db-service/pb"
	models "gen/models/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DBServer struct {
	dbpb.UnimplementedDBServiceServer
	Mongo *mongo.Database
}

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

func (s *DBServer) CreateCommunity(ctx context.Context, req *dbpb.CreateCommunityRequest) (*models.Community, error) {
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
	return community, nil
}

func (s *DBServer) GetCommunity(ctx context.Context, req *dbpb.GetCommunityRequest) (*models.Community, error) {
	collection := s.Mongo.Collection("communities")
	id, err := primitive.ObjectIDFromHex(req.Id)
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
	return &models.Community{
		Id:   community["_id"].(primitive.ObjectID).Hex(),
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

func (s *DBServer) CreateThread(ctx context.Context, req *dbpb.CreateThreadRequest) (*models.Thread, error) {
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
	return thread, nil
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

func (s *DBServer) ListComments(ctx context.Context, req *dbpb.ListCommentsRequest) (*dbpb.ListCommentsResponse, error) {
	collection := s.Mongo.Collection("comments")
	sortBy := strings.Replace(req.GetSortBy(), "votes", "ups", 1)
	findOptions := getFindOptions(req.GetLimit(), req.GetOffset(), sortBy)
	filter := bson.M{}
	if req.GetThreadId() != "" {
		threadId, err := primitive.ObjectIDFromHex(req.GetThreadId())
		if err != nil {
			return nil, fmt.Errorf("invalid thread ID: %v", err)
		}
		filter["thread_id"] = threadId
	}
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find comments: %v", err)
	}
	defer cursor.Close(ctx)
	var results []*models.Comment
	for cursor.Next(ctx) {
		var comment struct {
			ID         string `bson:"_id"`
			Content    string `bson:"content"`
			Ups        int32  `bson:"ups"`
			Downs      int32  `bson:"downs"`
			ParentId   string `bson:"parent_id"`
			ParentType string `bson:"parent_type"`
		}
		if err := cursor.Decode(&comment); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to decode comment: %v", err)
		}
		enumInt, ok := models.CommentParentType_value[comment.ParentType]
		if !ok {
			return nil, fmt.Errorf("invalid CommentParentType: %s", comment.ParentType)
		}
		results = append(results, &models.Comment{
			Id:         comment.ID,
			Content:    comment.Content,
			Ups:        comment.Ups,
			Downs:      comment.Downs,
			ParentId:   comment.ParentId,
			ParentType: models.CommentParentType(enumInt),
		})
	}
	if err := cursor.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "cursor error: %v", err)
	}
	return &dbpb.ListCommentsResponse{
		Comments: results,
	}, nil
}

func (s *DBServer) CreateComment(ctx context.Context, req *dbpb.CreateCommentRequest) (*models.Comment, error) {
	collection := s.Mongo.Collection("comments")
	commentID := generateUniqueId()
	comment := &models.Comment{
		Id:         commentID,
		Content:    req.GetContent(),
		Ups:        0,
		Downs:      0,
		ParentId:   req.GetParentId(),
		ParentType: req.GetParentType(),
	}
	doc := bson.M{
		"_id":         comment.Id,
		"content":     comment.Content,
		"ups":         comment.Ups,
		"downs":       comment.Downs,
		"parent_id":   comment.ParentId,
		"parent_type": comment.ParentType.String(),
	}

	if _, err := collection.InsertOne(ctx, doc); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create comment: %v", err)
	}
	return comment, nil
}

func (s *DBServer) GetComment(ctx context.Context, req *dbpb.GetCommentRequest) (*models.Comment, error) {
	collection := s.Mongo.Collection("comments")
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid comment ID: %v", err)
	}
	filter := bson.M{
		"_id": id,
	}
	var comment bson.M
	err = collection.FindOne(ctx, filter).Decode(&comment)
	if err != nil {
		return nil, err
	}
	enumInt, ok := models.CommentParentType_value[comment["parent_type"].(string)]
	if !ok {
		return nil, fmt.Errorf("invalid CommentParentType: %s", comment["parent_type"])
	}
	return &models.Comment{
		Id:         comment["_id"].(primitive.ObjectID).Hex(),
		Content:    comment["content"].(string),
		Ups:        comment["ups"].(int32),
		Downs:      comment["downs"].(int32),
		ParentId:   comment["parent_id"].(string),
		ParentType: models.CommentParentType(enumInt),
	}, nil
}

func (s *DBServer) UpdateComment(ctx context.Context, req *dbpb.UpdateCommentRequest) (*emptypb.Empty, error) {
	collection := s.Mongo.Collection("comments")
	update := bson.M{}
	setFields := bson.M{}
	incFields := bson.M{}

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
		return nil, status.Errorf(codes.Internal, "failed to update comment: %v", err)
	}
	if result.MatchedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "comment not found")
	}
	return &emptypb.Empty{}, nil
}

func (s *DBServer) DeleteComment(ctx context.Context, req *dbpb.DeleteCommentRequest) (*emptypb.Empty, error) {
	collection := s.Mongo.Collection("comments")

	// get comment
	commentRes, err := s.GetComment(ctx, &dbpb.GetCommentRequest{Id: req.GetId()})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get comment: %v", err)
	}

	// attempt deletion
	result, err := collection.DeleteOne(ctx, bson.M{"_id": req.GetId()})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete comment: %v", err)
	}
	if result.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "comment not found")
	}

	// remove comment id from the parent
	var parentUpdate bson.M
	if commentRes.GetParentType() == models.CommentParentType_THREAD {
		parentUpdate = bson.M{
			"$pull": bson.M{
				"comments": req.GetId(),
			},
		}
	} else if commentRes.GetParentType() == models.CommentParentType_COMMENT {
		parentUpdate = bson.M{
			"$pull": bson.M{
				"comments": req.GetId(),
			},
		}
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent type")
	}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": commentRes.GetParentId()}, parentUpdate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update parent with comment: %v", err)
	}
	
	return &emptypb.Empty{}, nil
}


func getFindOptions(offset int32, limit int32, sortBy string) *options.FindOptions {
	findOptions := options.Find()
	if offset > 0 {
		findOptions.SetSkip(int64(offset))
	}
	if limit > 0 {
		findOptions.SetLimit(int64(limit))
	}
	if sortBy != "" {
		findOptions.SetSort(bson.D{{Key: sortBy, Value: -1}}) // sort in descending order
	}
	return findOptions
}

func generateUniqueId() string {
	id := primitive.NewObjectID()
	return id.Hex()
}

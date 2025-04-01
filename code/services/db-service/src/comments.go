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

func (s *DBServer) ListComments(ctx context.Context, req *dbpb.ListCommentsRequest) (*dbpb.ListCommentsResponse, error) {
	collection := s.Mongo.Collection("comments")
	findOptions := getFindOptions(req.Offset, req.Limit, *req.SortBy)
	filter := bson.M{}
	if req.GetThreadId() != "" {
		filter["thread_id"] = req.GetThreadId()
	}
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Printf("list comments: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to find comments")
	}
	defer cursor.Close(ctx)
	var results []*models.Comment
	for cursor.Next(ctx) {
		var comment struct {
			ID          string `bson:"_id"`
			Content     string `bson:"content"`
			Ups         int32  `bson:"ups"`
			Downs       int32  `bson:"downs"`
			ParentId    string `bson:"parent_id"`
			ParentType  string `bson:"parent_type"`
			NumComments int32  `bson:"num_comments"`
		}
		if err := cursor.Decode(&comment); err != nil {
			log.Printf("list comments: %v", err)
			return nil, status.Errorf(codes.Internal, "failed to decode comment")
		}
		enumInt, ok := models.CommentParentType_value[comment.ParentType]
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "invalid parent type: %v", comment.ParentType)
		}
		results = append(results, &models.Comment{
			Id:          comment.ID,
			Content:     comment.Content,
			Ups:         comment.Ups,
			Downs:       comment.Downs,
			ParentId:    comment.ParentId,
			ParentType:  models.CommentParentType(enumInt),
			NumComments: comment.NumComments,
		})
	}
	if err := cursor.Err(); err != nil {
		log.Printf("list comments: %v", err)
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	return &dbpb.ListCommentsResponse{
		Comments: results,
	}, nil
}

func (s *DBServer) CreateComment(ctx context.Context, req *dbpb.CreateCommentRequest) (*dbpb.CreateCommentResponse, error) {
	collection := s.Mongo.Collection("comments")
	commentID := generateUniqueId()
	comment := &models.Comment{
		Id:          commentID,
		Content:     req.GetContent(),
		Ups:         0,
		Downs:       0,
		ParentId:    req.GetParentId(),
		ParentType:  req.GetParentType(),
		NumComments: 0,
	}
	doc := bson.M{
		"_id":          comment.Id,
		"content":      comment.Content,
		"ups":          comment.Ups,
		"downs":        comment.Downs,
		"parent_id":    comment.ParentId,
		"parent_type":  comment.ParentType.String(),
		"num_comments": comment.NumComments,
	}

	if _, err := collection.InsertOne(ctx, doc); err != nil {
		log.Printf("create comment: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create comment")
	}
	return &dbpb.CreateCommentResponse{
		Id: commentID,
	}, nil
}

func (s *DBServer) GetComment(ctx context.Context, req *dbpb.GetCommentRequest) (*models.Comment, error) {
	collection := s.Mongo.Collection("comments")
	filter := bson.M{
		"_id": req.GetId(),
	}
	var comment bson.M
	err := collection.FindOne(ctx, filter).Decode(&comment)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.NotFound, "comment not found")
		}
		log.Printf("get comment: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get comment")
	}
	enumInt, ok := models.CommentParentType_value[comment["parent_type"].(string)]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent type: %v", comment["parent_type"])
	}
	return &models.Comment{
		Id:          comment["_id"].(string),
		Content:     comment["content"].(string),
		Ups:         comment["ups"].(int32),
		Downs:       comment["downs"].(int32),
		ParentId:    comment["parent_id"].(string),
		ParentType:  models.CommentParentType(enumInt),
		NumComments: comment["num_comments"].(int32),
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
		log.Printf("update comment: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update comment")
	}
	if result.MatchedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "comment not found")
	}
	return &emptypb.Empty{}, nil
}

func (s *DBServer) DeleteComment(ctx context.Context, req *dbpb.DeleteCommentRequest) (*emptypb.Empty, error) {
	collection := s.Mongo.Collection("comments")
	result, err := collection.DeleteOne(ctx, bson.M{"_id": req.GetId()})
	if err != nil {
		log.Printf("delete comment: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete comment")
	}
	if result.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "comment not found")
	}
	return &emptypb.Empty{}, nil
}

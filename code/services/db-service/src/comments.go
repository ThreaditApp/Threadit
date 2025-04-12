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

func (s *DBServer) ListComments(ctx context.Context, req *dbpb.ListCommentsRequest) (*dbpb.ListCommentsResponse, error) {
	collection := s.Mongo.Collection("comments")
	filter := bson.M{}
	if req.GetThreadId() != "" {
		filter["thread_id"] = req.GetThreadId()
	}
	cursor, err := collection.Find(ctx, filter, getFindOptions(req.Offset, req.Limit, req.GetSortBy()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to find comments")
	}
	defer cursor.Close(ctx)

	var results []*models.Comment
	for cursor.Next(ctx) {
		comment := bson.M{}
		if err := cursor.Decode(&comment); err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to decode comment")
		}
		enumInt, _ := models.CommentParentType_value[comment["parent_type"].(string)]
		results = append(results, &models.Comment{
			Id:          comment["_id"].(string),
			Content:     comment["content"].(string),
			Ups:         comment["ups"].(int32),
			Downs:       comment["downs"].(int32),
			ParentId:    comment["parent_id"].(string),
			ParentType:  models.CommentParentType(enumInt),
			NumComments: comment["num_comments"].(int32),
		})
	}
	if err := cursor.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "Unexpected cursor error")
	}
	return &dbpb.ListCommentsResponse{
		Comments: results,
	}, nil
}

func (s *DBServer) CreateComment(ctx context.Context, req *dbpb.CreateCommentRequest) (*dbpb.CreateCommentResponse, error) {
	collection := s.Mongo.Collection("comments")
	comment := bson.M{
		"_id":          generateUniqueId(),
		"content":      req.GetContent(),
		"ups":          0,
		"downs":        0,
		"parent_id":    req.GetParentId(),
		"parent_type":  req.GetParentType().String(),
		"num_comments": 0,
	}

	if _, err := collection.InsertOne(ctx, comment); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create comment")
	}
	return &dbpb.CreateCommentResponse{
		Id: comment["_id"].(string),
	}, nil
}

func (s *DBServer) GetComment(ctx context.Context, req *dbpb.GetCommentRequest) (*models.Comment, error) {
	collection := s.Mongo.Collection("comments")
	filter := bson.M{"_id": req.GetId()}

	var comment bson.M
	err := collection.FindOne(ctx, filter).Decode(&comment)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.NotFound, "Comment not found")
		}
		return nil, status.Errorf(codes.Internal, "Failed to get comment")
	}

	enumInt, _ := models.CommentParentType_value[comment["parent_type"].(string)]
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
	setValues, incValues := bson.M{}, bson.M{}

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
		return nil, status.Errorf(codes.Internal, "Failed to update comment")
	}
	if result.MatchedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "Comment not found")
	}
	return &emptypb.Empty{}, nil
}

func (s *DBServer) DeleteComment(ctx context.Context, req *dbpb.DeleteCommentRequest) (*emptypb.Empty, error) {
	collection := s.Mongo.Collection("comments")

	result, err := collection.DeleteOne(ctx, bson.M{"_id": req.GetId()})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete comment")
	}
	if result.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "Comment not found")
	}

	return &emptypb.Empty{}, nil
}

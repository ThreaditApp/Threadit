package server

import (
	"context"
	commentpb "gen/comment-service/pb"
	dbpb "gen/db-service/pb"
	models "gen/models/pb"
	threadpb "gen/thread-service/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CommentServer struct {
	commentpb.UnimplementedCommentServiceServer
	DBClient     dbpb.DBServiceClient
	ThreadClient threadpb.ThreadServiceClient
}

func (s *CommentServer) ListComments(ctx context.Context, req *commentpb.ListCommentsRequest) (*commentpb.ListCommentsResponse, error) {
	// input validation
	if req.ThreadId != nil && *req.ThreadId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "thread id cannot be empty")
	}
	if req.Offset != nil && *req.Offset < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "offset must be a non-negative integer")
	}
	if req.Limit != nil && *req.Limit <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "limit must be a positive integer")
	}

	// call the database service to fetch comments
	res, err := s.DBClient.ListComments(ctx, &dbpb.ListCommentsRequest{
		ThreadId: req.ThreadId,
		Offset:   req.Offset,
		Limit:    req.Limit,
		SortBy:   req.SortBy,
	})
	if err != nil {
		return nil, err
	}
	return &commentpb.ListCommentsResponse{
		Comments: res.Comments,
	}, nil
}

func (s *CommentServer) CreateComment(ctx context.Context, req *commentpb.CreateCommentRequest) (*commentpb.CreateCommentResponse, error) {
	// input validation
	if req.ParentId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "parent id is required")
	}
	if req.Content == "" {
		return nil, status.Errorf(codes.InvalidArgument, "content is required")
	}
	if len(req.Content) > 500 {
		return nil, status.Errorf(codes.InvalidArgument, "content exceeds maximum length of 500 characters")
	}

	// TODO: update parent num_comments

	// create comment
	res, err := s.DBClient.CreateComment(ctx, &dbpb.CreateCommentRequest{
		Content:    req.Content,
		ParentId:   req.ParentId,
		ParentType: req.ParentType,
	})
	if err != nil {
		return nil, err
	}
	return &commentpb.CreateCommentResponse{
		Id: res.Id,
	}, nil
}

func (s *CommentServer) GetComment(ctx context.Context, req *commentpb.GetCommentRequest) (*models.Comment, error) {
	// Input validation
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	// fetch comment
	res, err := s.DBClient.GetComment(ctx, &dbpb.GetCommentRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *CommentServer) UpdateComment(ctx context.Context, req *commentpb.UpdateCommentRequest) (*emptypb.Empty, error) {
	// input validation
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}
	if req.Content != nil && len(*req.Content) > 500 {
		return nil, status.Errorf(codes.InvalidArgument, "content exceeds maximum length of 500 characters")
	}

	// update comment
	_, err := s.DBClient.UpdateComment(ctx, &dbpb.UpdateCommentRequest{
		Id:      req.Id,
		Content: req.Content,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *CommentServer) DeleteComment(ctx context.Context, req *commentpb.DeleteCommentRequest) (*emptypb.Empty, error) {
	// input validation
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	// call the database service to delete a comment
	_, err := s.DBClient.DeleteComment(ctx, &dbpb.DeleteCommentRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	// TODO: update parent num_comments

	return &emptypb.Empty{}, nil
}

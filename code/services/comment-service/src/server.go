package server

import (
	"context"
	commentpb "gen/comment-service/pb"
	dbpb "gen/db-service/pb"
	models "gen/models/pb"
	threadpb "gen/thread-service/pb"
	"math"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CommentServer struct {
	commentpb.UnimplementedCommentServiceServer
	DBClient     dbpb.DBServiceClient
	ThreadClient threadpb.ThreadServiceClient
}

const (
	MaxCommentLength = 500
)

func (s *CommentServer) CheckHealth(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *CommentServer) ListComments(ctx context.Context, req *commentpb.ListCommentsRequest) (*commentpb.ListCommentsResponse, error) {
	// validate inputs
	if req.ThreadId != nil && req.GetThreadId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Thread id cannot be empty")
	}
	if req.Offset != nil && req.GetOffset() < 0 {
		return nil, status.Error(codes.InvalidArgument, "Offset must be a positive integer")
	}
	if req.Limit != nil && req.GetLimit() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Limit must be a positive integer")
	}
	if req.SortBy != nil && req.GetSortBy() == "" {
		return nil, status.Error(codes.InvalidArgument, "Sort cannot be empty")
	}

	// fetch comments
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
	// validate inputs
	if req.GetParentId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Parent id is required")
	}
	if req.GetContent() == "" {
		return nil, status.Error(codes.InvalidArgument, "Content is required")
	}
	if len(req.GetContent()) > MaxCommentLength {
		return nil, status.Errorf(codes.InvalidArgument, "Content exceeds maximum length of %d characters", MaxCommentLength)
	}

	// create comment
	res, err := s.DBClient.CreateComment(ctx, &dbpb.CreateCommentRequest{
		Content:    req.Content,
		ParentId:   req.ParentId,
		ParentType: req.ParentType,
	})
	if err != nil {
		return nil, err
	}

	// update parent num_comments
	numCommentsOffset := int32(1)
	if req.ParentType == models.CommentParentType_THREAD {
		_, err = s.ThreadClient.UpdateThread(ctx, &threadpb.UpdateThreadRequest{
			Id:                req.ParentId,
			NumCommentsOffset: &numCommentsOffset,
		})
		if err != nil {
			return nil, err
		}
	} else {
		_, err = s.UpdateComment(ctx, &commentpb.UpdateCommentRequest{
			Id:                req.ParentId,
			NumCommentsOffset: &numCommentsOffset,
		})
	}

	return &commentpb.CreateCommentResponse{
		Id: res.Id,
	}, nil
}

func (s *CommentServer) GetComment(ctx context.Context, req *commentpb.GetCommentRequest) (*models.Comment, error) {
	// validate input
	if req.GetId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Comment id is required")
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
	// validate inputs
	if req.GetId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Comment id is required")
	}
	if req.Content != nil && len(req.GetContent()) > MaxCommentLength {
		return nil, status.Errorf(codes.InvalidArgument, "Content exceeds maximum length of %d characters", MaxCommentLength)
	}
	if req.VoteOffset != nil && math.Abs(float64(req.GetVoteOffset())) != 1 {
		return nil, status.Errorf(codes.InvalidArgument, "Vote offset must be either -1 or 1")
	}
	if req.NumCommentsOffset != nil && math.Abs(float64(req.GetNumCommentsOffset())) != 1 {
		return nil, status.Errorf(codes.InvalidArgument, "Number comments offset must be either -1 or 1")
	}

	// update comment
	_, err := s.DBClient.UpdateComment(ctx, &dbpb.UpdateCommentRequest{
		Id:                req.Id,
		Content:           req.Content,
		VoteOffset:        req.VoteOffset,
		NumCommentsOffset: req.NumCommentsOffset,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *CommentServer) DeleteComment(ctx context.Context, req *commentpb.DeleteCommentRequest) (*emptypb.Empty, error) {
	// validate input
	if req.GetId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Comment id is required")
	}

	// fetch comment
	res, err := s.DBClient.GetComment(ctx, &dbpb.GetCommentRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	// delete comment
	_, err = s.DBClient.DeleteComment(ctx, &dbpb.DeleteCommentRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	// update parent num_comments
	numCommentsOffset := int32(-1)
	if res.ParentType == models.CommentParentType_THREAD {
		_, err = s.ThreadClient.UpdateThread(ctx, &threadpb.UpdateThreadRequest{
			Id:                res.ParentId,
			NumCommentsOffset: &numCommentsOffset,
		})
		if err != nil {
			return nil, err
		}
	} else {
		_, err = s.UpdateComment(ctx, &commentpb.UpdateCommentRequest{
			Id:                res.ParentId,
			NumCommentsOffset: &numCommentsOffset,
		})
	}

	return &emptypb.Empty{}, nil
}

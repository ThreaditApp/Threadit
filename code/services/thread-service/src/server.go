package server

import (
	"context"
	commentpb "gen/comment-service/pb"
	communitypb "gen/community-service/pb"
	dbpb "gen/db-service/pb"
	models "gen/models/pb"
	threadpb "gen/thread-service/pb"
	"math"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ThreadServer struct {
	threadpb.UnimplementedThreadServiceServer
	CommunityClient communitypb.CommunityServiceClient
	DBClient        dbpb.DBServiceClient
	CommentClient   commentpb.CommentServiceClient
}

const (
	MinTitleLength   = 3
	MaxTitleLength   = 50
	MaxContentLength = 500
	MinContentLength = 3
)

func (s *ThreadServer) CheckHealth(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *ThreadServer) ListThreads(ctx context.Context, req *threadpb.ListThreadsRequest) (*threadpb.ListThreadsResponse, error) {
	// validate inputs
	if req.CommunityId != nil && req.GetCommunityId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Community id cannot be empty")
	}
	if req.Title != nil && req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "Title cannot be empty")
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

	// fetch threads
	res, err := s.DBClient.ListThreads(ctx, &dbpb.ListThreadsRequest{
		CommunityId: req.CommunityId,
		Title:       req.Title,
		Offset:      req.Offset,
		Limit:       req.Limit,
		SortBy:      req.SortBy,
	})
	if err != nil {
		return nil, err
	}
	return &threadpb.ListThreadsResponse{
		Threads: res.Threads,
	}, nil
}

func (s *ThreadServer) CreateThread(ctx context.Context, req *threadpb.CreateThreadRequest) (*threadpb.CreateThreadResponse, error) {
	// validate inputs
	if req.GetCommunityId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Community id is required")
	}
	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "Title is required")
	}
	if req.GetContent() == "" {
		return nil, status.Error(codes.InvalidArgument, "Content is required")
	}
	if len(req.GetTitle()) < MinTitleLength || len(req.GetTitle()) > MaxTitleLength {
		return nil, status.Errorf(codes.InvalidArgument, "Title must be between %d and %d characters long", MinTitleLength, MaxTitleLength)
	}
	if len(req.GetContent()) < MinContentLength || len(req.GetContent()) > MaxContentLength {
		return nil, status.Errorf(codes.InvalidArgument, "Content must be between %d and %d characters long", MinContentLength, MaxContentLength)
	}

	// create thread
	res, err := s.DBClient.CreateThread(ctx, &dbpb.CreateThreadRequest{
		CommunityId: req.CommunityId,
		Title:       req.Title,
		Content:     req.Content,
	})
	if err != nil {
		return nil, err
	}

	// update community num_threads
	numThreadsOffset := int32(1)
	_, err = s.CommunityClient.UpdateCommunity(ctx, &communitypb.UpdateCommunityRequest{
		Id:               req.CommunityId,
		NumThreadsOffset: &numThreadsOffset,
	})
	if err != nil {
		return nil, err
	}

	return &threadpb.CreateThreadResponse{
		Id: res.Id,
	}, nil
}

func (s *ThreadServer) GetThread(ctx context.Context, req *threadpb.GetThreadRequest) (*models.Thread, error) {
	// validate inputs
	if req.GetId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Thread id is required")
	}

	// fetch thread
	res, err := s.DBClient.GetThread(ctx, &dbpb.GetThreadRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *ThreadServer) UpdateThread(ctx context.Context, req *threadpb.UpdateThreadRequest) (*emptypb.Empty, error) {
	// validate inputs
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Thread id is required")
	}
	titleLen := len(req.GetTitle())
	if req.Title != nil && (titleLen < MinTitleLength || titleLen > MaxTitleLength) {
		return nil, status.Errorf(codes.InvalidArgument, "Title must be between %d and %d characters long", MinTitleLength, MaxTitleLength)
	}
	contentLen := len(req.GetContent())
	if req.Content != nil && (contentLen < MinContentLength || contentLen > MaxContentLength) {
		return nil, status.Errorf(codes.InvalidArgument, "Content must be between %d and %d characters long", MinContentLength, MaxContentLength)
	}
	if req.VoteOffset != nil && math.Abs(float64(req.GetVoteOffset())) != 1 {
		return nil, status.Error(codes.InvalidArgument, "Vote offset must be either -1 or 1")
	}
	if req.NumCommentsOffset != nil && math.Abs(float64(req.GetNumCommentsOffset())) != 1 {
		return nil, status.Error(codes.InvalidArgument, "Number comments offset must be either -1 or 1")
	}

	// update thread
	_, err := s.DBClient.UpdateThread(ctx, &dbpb.UpdateThreadRequest{
		Id:                req.Id,
		Title:             req.Title,
		Content:           req.Content,
		VoteOffset:        req.VoteOffset,
		NumCommentsOffset: req.NumCommentsOffset,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *ThreadServer) DeleteThread(ctx context.Context, req *threadpb.DeleteThreadRequest) (*emptypb.Empty, error) {
	// validate input
	if req.GetId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Thread id is required")
	}

	// get thread
	thread, err := s.DBClient.GetThread(ctx, &dbpb.GetThreadRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	// find all comments from thread
	comments, err := s.CommentClient.ListComments(ctx, &commentpb.ListCommentsRequest{
		ThreadId: &req.Id,
	})
	if err != nil {
		return nil, err
	}

	// delete comments
	for _, comment := range comments.Comments {
		_, err = s.CommentClient.DeleteComment(ctx, &commentpb.DeleteCommentRequest{
			Id: comment.Id,
		})
		if err != nil {
			return nil, err
		}
	}

	// delete thread
	_, err = s.DBClient.DeleteThread(ctx, &dbpb.DeleteThreadRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	// update community num_threads
	numThreadsOffset := int32(-1)
	_, err = s.CommunityClient.UpdateCommunity(ctx, &communitypb.UpdateCommunityRequest{
		Id:               thread.CommunityId,
		NumThreadsOffset: &numThreadsOffset,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

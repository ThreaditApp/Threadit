package server

import (
	"context"
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
}

const (
	MinThreadTitleLength   = 3
	MaxThreadTitleLength   = 50
	MaxThreadContentLength = 500
	MinThreadContentLength = 3
)

func (s *ThreadServer) ListThreads(ctx context.Context, req *threadpb.ListThreadsRequest) (*threadpb.ListThreadsResponse, error) {
	// validate inputs
	if req.CommunityId != nil && req.GetCommunityId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "community id cannot be empty")
	}
	if req.Title != nil && req.GetTitle() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "title cannot be empty")
	}
	if req.Offset != nil && req.GetOffset() < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "offset must be a non-negative integer")
	}
	if req.Limit != nil && req.GetLimit() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "limit must be a positive integer")
	}
	if req.SortBy != nil && req.GetSortBy() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "sort cannot be empty")
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
		return nil, status.Errorf(codes.InvalidArgument, "community id is required")
	}
	if req.GetTitle() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "title is required")
	}
	if req.GetContent() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "content is required")
	}
	if len(req.GetTitle()) < MinThreadTitleLength || len(req.GetTitle()) > MaxThreadTitleLength {
		return nil, status.Errorf(codes.InvalidArgument, "title must be between %d and %d characters long", MinThreadTitleLength, MaxThreadTitleLength)
	}
	if len(req.GetContent()) < MinThreadContentLength || len(req.GetContent()) > MaxThreadContentLength {
		return nil, status.Errorf(codes.InvalidArgument, "content must be between %d and %d characters long", MinThreadContentLength, MaxThreadContentLength)
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
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
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
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}
	if req.Title != nil && req.GetTitle() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "title cannot be empty")
	}
	if req.Content != nil && (len(req.GetContent()) < MinThreadContentLength || len(req.GetContent()) > MaxThreadContentLength) {
		return nil, status.Errorf(codes.InvalidArgument, "content must be between %d and %d characters long", MinThreadContentLength, MaxThreadContentLength)
	}
	if req.VoteOffset != nil && math.Abs(float64(req.GetVoteOffset())) != 1 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid vote offset")
	}
	if req.NumCommentsOffset != nil && math.Abs(float64(req.GetNumCommentsOffset())) != 1 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid num comments offset")
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
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	// get thread
	thread, err := s.DBClient.GetThread(ctx, &dbpb.GetThreadRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
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

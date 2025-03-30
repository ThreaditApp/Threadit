package server

import (
	"context"
	communitypb "gen/community-service/pb"
	dbpb "gen/db-service/pb"
	models "gen/models/pb"
	threadpb "gen/thread-service/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ThreadServer struct {
	threadpb.UnimplementedThreadServiceServer
	CommunityClient communitypb.CommunityServiceClient
	DBClient        dbpb.DBServiceClient
}

func (s *ThreadServer) ListThreads(ctx context.Context, req *threadpb.ListThreadsRequest) (*threadpb.ListThreadsResponse, error) {
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
	if req.CommunityId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "community_id is required")
	}
	if req.Title == "" {
		return nil, status.Errorf(codes.InvalidArgument, "title is required")
	}
	if len(req.Title) < 3 || len(req.Title) > 50 {
		return nil, status.Errorf(codes.InvalidArgument, "title must be between 3 and 50 characters long")
	}
	if req.Content == "" {
		return nil, status.Errorf(codes.InvalidArgument, "content is required")
	}
	if len(req.Content) < 3 || len(req.Content) > 500 {
		return nil, status.Errorf(codes.InvalidArgument, "content must be between 3 and 500 characters long")
	}

	// check if community exists
	_, err := s.CommunityClient.GetCommunity(ctx, &communitypb.GetCommunityRequest{
		Id: req.CommunityId,
	})
	if err != nil {
		return nil, err
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
	return &threadpb.CreateThreadResponse{
		Id: res.Id,
	}, nil
}

func (s *ThreadServer) GetThread(ctx context.Context, req *threadpb.GetThreadRequest) (*models.Thread, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}
	res, err := s.DBClient.GetThread(ctx, &dbpb.GetThreadRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *ThreadServer) UpdateThread(ctx context.Context, req *threadpb.UpdateThreadRequest) (*emptypb.Empty, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}
	if req.Title == nil && req.Content == nil {
		return nil, status.Errorf(codes.InvalidArgument, "at least one of title or content is required")
	}
	if req.Title != nil && *req.Title != "" && (len(*req.Title) < 3 || len(*req.Title) > 50) {
		return nil, status.Errorf(codes.InvalidArgument, "title must be between 3 and 50 characters long")
	}
	if req.Content != nil && *req.Content != "" && (len(*req.Content) < 3 || len(*req.Content) > 500) {
		return nil, status.Errorf(codes.InvalidArgument, "content must be between 3 and 500 characters long")
	}
	_, err := s.DBClient.UpdateThread(ctx, &dbpb.UpdateThreadRequest{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *ThreadServer) DeleteThread(ctx context.Context, req *threadpb.DeleteThreadRequest) (*emptypb.Empty, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}
	_, err := s.DBClient.DeleteThread(ctx, &dbpb.DeleteThreadRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

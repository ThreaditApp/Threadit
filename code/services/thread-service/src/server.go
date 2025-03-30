package server

import (
	"context"
	"fmt"
	communitypb "gen/community-service/pb"
	dbpb "gen/db-service/pb"
	threadpb "gen/thread-service/pb"
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
		return nil, fmt.Errorf("error calling database service: %w", err)
	}
	return &threadpb.ListThreadsResponse{
		Threads: res.Threads,
	}, nil
}

func (s *ThreadServer) CreateThread(ctx context.Context, req *threadpb.CreateThreadRequest) (*threadpb.CreateThreadResponse, error) {
	// check if community exists
	_, err := s.CommunityClient.GetCommunity(ctx, &communitypb.GetCommunityRequest{
		Id: req.CommunityId,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling community service: %w", err)
	}

	// create thread
	res, err := s.DBClient.CreateThread(ctx, &dbpb.CreateThreadRequest{
		CommunityId: req.CommunityId,
		Title:       req.Title,
		Content:     req.Content,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}
	return &threadpb.CreateThreadResponse{
		Id: res.Id,
	}, nil
}

func (s *ThreadServer) GetThread(ctx context.Context, req *threadpb.GetThreadRequest) (*threadpb.GetThreadResponse, error) {
	res, err := s.DBClient.GetThread(ctx, &dbpb.GetThreadRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}
	return &threadpb.GetThreadResponse{
		Thread: res.Thread,
	}, nil
}

func (s *ThreadServer) UpdateThread(ctx context.Context, req *threadpb.UpdateThreadRequest) (*emptypb.Empty, error) {
	_, err := s.DBClient.UpdateThread(ctx, &dbpb.UpdateThreadRequest{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *ThreadServer) DeleteThread(ctx context.Context, req *threadpb.DeleteThreadRequest) (*emptypb.Empty, error) {
	_, err := s.DBClient.DeleteThread(ctx, &dbpb.DeleteThreadRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}
	return &emptypb.Empty{}, nil
}

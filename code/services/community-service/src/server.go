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

type CommunityServer struct {
	communitypb.UnimplementedCommunityServiceServer
	DBClient     dbpb.DBServiceClient
	ThreadClient threadpb.ThreadServiceClient
}

const (
	MinLength = 3
	MaxLength = 50
)

func (s *CommunityServer) CheckHealth(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *CommunityServer) ListCommunities(ctx context.Context, req *communitypb.ListCommunitiesRequest) (*communitypb.ListCommunitiesResponse, error) {
	// validate inputs
	if req.Offset != nil && req.GetOffset() < 0 {
		return nil, status.Error(codes.InvalidArgument, "Offset must be a non-negative integer")
	}
	if req.Limit != nil && req.GetLimit() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Limit must be a positive integer")
	}

	// fetch communities
	res, err := s.DBClient.ListCommunities(ctx, &dbpb.ListCommunitiesRequest{
		Name:   req.Name,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		return nil, err
	}
	return &communitypb.ListCommunitiesResponse{
		Communities: res.Communities,
	}, nil
}

func (s *CommunityServer) CreateCommunity(ctx context.Context, req *communitypb.CreateCommunityRequest) (*communitypb.CreateCommunityResponse, error) {
	// validate inputs
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "Community name is required")
	}
	if len(req.GetName()) < MinLength || len(req.GetName()) > MaxLength {
		return nil, status.Errorf(codes.InvalidArgument, "Name must be between %d and %d characters long", MinLength, MaxLength)
	}

	// create community
	res, err := s.DBClient.CreateCommunity(ctx, &dbpb.CreateCommunityRequest{
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}
	return &communitypb.CreateCommunityResponse{
		Id: res.Id,
	}, nil
}

func (s *CommunityServer) GetCommunity(ctx context.Context, req *communitypb.GetCommunityRequest) (*models.Community, error) {
	// validate input
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Community id is required")
	}

	// fetch community
	res, err := s.DBClient.GetCommunity(ctx, &dbpb.GetCommunityRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *CommunityServer) UpdateCommunity(ctx context.Context, req *communitypb.UpdateCommunityRequest) (*emptypb.Empty, error) {
	// validate inputs
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Community id is required")
	}
	nameLen := len(req.GetName())
	if req.Name != nil && (nameLen < MinLength || nameLen > MaxLength) {
		return nil, status.Errorf(codes.InvalidArgument, "Name must be between %d and %d characters long", MinLength, MaxLength)
	}
	if req.NumThreadsOffset != nil && math.Abs(float64(req.GetNumThreadsOffset())) != 1 {
		return nil, status.Error(codes.InvalidArgument, "Number of threads offset must be either -1 or 1")
	}

	// update community
	_, err := s.DBClient.UpdateCommunity(ctx, &dbpb.UpdateCommunityRequest{
		Id:               req.Id,
		Name:             req.Name,
		NumThreadsOffset: req.NumThreadsOffset,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CommunityServer) DeleteCommunity(ctx context.Context, req *communitypb.DeleteCommunityRequest) (*emptypb.Empty, error) {
	// validate input
	if req.GetId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Community id is required")
	}

	// find all threads from community
	res, err := s.ThreadClient.ListThreads(ctx, &threadpb.ListThreadsRequest{
		CommunityId: &req.Id,
	})
	if err != nil {
		return nil, err
	}

	// delete threads
	for _, thread := range res.Threads {
		_, err = s.ThreadClient.DeleteThread(ctx, &threadpb.DeleteThreadRequest{
			Id: thread.Id,
		})
		if err != nil {
			return nil, err
		}
	}

	// delete community
	_, err = s.DBClient.DeleteCommunity(ctx, &dbpb.DeleteCommunityRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

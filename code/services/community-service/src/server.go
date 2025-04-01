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
	MinCommunityNameLength = 3
	MaxCommunityNameLength = 50
)

func (s *CommunityServer) ListCommunities(ctx context.Context, req *communitypb.ListCommunitiesRequest) (*communitypb.ListCommunitiesResponse, error) {
	// validate inputs
	if req.Name != nil && len(req.GetName()) < MinCommunityNameLength {
		return nil, status.Errorf(codes.InvalidArgument, "invalid community name %s", req.Name)
	}
	if req.Offset != nil && req.GetOffset() < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "offset must be a non-negative integer")
	}
	if req.Limit != nil && req.GetLimit() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "limit must be a positive integer")
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
	if len(req.GetName()) < MinCommunityNameLength {
		return nil, status.Errorf(codes.InvalidArgument, "invalid community name")
	}
	if len(req.Name) < MinCommunityNameLength || len(req.Name) > MaxCommunityNameLength {
		return nil, status.Errorf(codes.InvalidArgument, "name must be between %d and %d characters long", MinCommunityNameLength, MaxCommunityNameLength)
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
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
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
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}
	if req.GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}
	if len(req.GetName()) < MinCommunityNameLength || len(req.GetName()) > MaxCommunityNameLength {
		return nil, status.Errorf(codes.InvalidArgument, "name must be between %d and %d characters long", MinCommunityNameLength, MaxCommunityNameLength)
	}
	if math.Abs(float64(req.GetNumThreadsOffset())) != 1 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid num threads offset")
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
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	// delete community
	_, err := s.DBClient.DeleteCommunity(ctx, &dbpb.DeleteCommunityRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
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
	return &emptypb.Empty{}, nil
}

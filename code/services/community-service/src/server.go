package server

import (
	"context"
	communitypb "gen/community-service/pb"
	dbpb "gen/db-service/pb"
	models "gen/models/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CommunityServer struct {
	communitypb.UnimplementedCommunityServiceServer
	DBClient dbpb.DBServiceClient
}

func (s *CommunityServer) ListCommunities(ctx context.Context, req *communitypb.ListCommunitiesRequest) (*communitypb.ListCommunitiesResponse, error) {
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
	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}
	if len(req.Name) < 3 || len(req.Name) > 50 {
		return nil, status.Errorf(codes.InvalidArgument, "name must be between 3 and 50 characters long")
	}
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
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}
	res, err := s.DBClient.GetCommunity(ctx, &dbpb.GetCommunityRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *CommunityServer) UpdateCommunity(ctx context.Context, req *communitypb.UpdateCommunityRequest) (*emptypb.Empty, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}
	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}
	if len(req.Name) < 3 || len(req.Name) > 50 {
		return nil, status.Errorf(codes.InvalidArgument, "name must be between 3 and 50 characters long")
	}
	_, err := s.DBClient.UpdateCommunity(ctx, &dbpb.UpdateCommunityRequest{
		Id:   req.Id,
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CommunityServer) DeleteCommunity(ctx context.Context, req *communitypb.DeleteCommunityRequest) (*emptypb.Empty, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}
	_, err := s.DBClient.DeleteCommunity(ctx, &dbpb.DeleteCommunityRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

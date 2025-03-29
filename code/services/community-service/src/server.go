package server

import (
	"context"
	"fmt"
	communitypb "gen/community-service/pb"
	dbpb "gen/db-service/pb"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CommunityServer struct {
	communitypb.UnimplementedCommunityServiceServer
	DBClient dbpb.DBServiceClient
}

func (s *CommunityServer) ListCommunities(ctx context.Context, req *communitypb.ListCommunitiesRequest) (*communitypb.ListCommunitiesResponse, error) {
	res, err := s.DBClient.ListCommunities(ctx, &dbpb.ListCommunitiesRequest{
		Page:      req.Page,
		PageSize:  req.PageSize,
		OwnerId:   req.OwnerId,
		Name:      req.Name,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	communities := make([]*communitypb.Community, len(res.Communities))
	for i, community := range res.Communities {
		communities[i] = &communitypb.Community{
			Id:          community.Id,
			OwnerId:     community.OwnerId,
			Name:        community.Name,
			Description: community.Description,
			CreatedAt:   community.CreatedAt,
			UpdatedAt:   community.UpdatedAt,
		}
	}

	return &communitypb.ListCommunitiesResponse{
		Communities: communities,
		Pagination: &communitypb.Pagination{
			CurrentPage: res.Pagination.CurrentPage,
			PerPage:     res.Pagination.PerPage,
			TotalItems:  res.Pagination.TotalItems,
			TotalPages:  res.Pagination.TotalPages,
		},
	}, nil
}

func (s *CommunityServer) CreateCommunity(ctx context.Context, req *communitypb.CreateCommunityRequest) (*communitypb.Community, error) {

	res, err := s.DBClient.CreateCommunity(ctx, &dbpb.CreateCommunityRequest{
		Name: req.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	return &communitypb.Community{
		Id:   res.Id,
		Name: res.Name,
	}, nil
}

func (s *CommunityServer) GetCommunity(ctx context.Context, req *communitypb.GetCommunityRequest) (*communitypb.Community, error) {
	res, err := s.DBClient.GetCommunity(ctx, &dbpb.GetCommunityRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	return &communitypb.Community{
		Id:   res.Id,
		Name: res.Name,
	}, nil
}

func (s *CommunityServer) UpdateCommunity(ctx context.Context, req *communitypb.UpdateCommunityRequest) (*communitypb.Community, error) {
	community, err := s.DBClient.GetCommunity(ctx, &dbpb.GetCommunityRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	res, err := s.DBClient.UpdateCommunity(ctx, &dbpb.UpdateCommunityRequest{
		Id:   req.Id,
		Name: req.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	return &communitypb.Community{
		Id:   res.Id,
		Name: res.Name,
	}, nil
}

func (s *CommunityServer) DeleteCommunity(ctx context.Context, req *communitypb.DeleteCommunityRequest) (*emptypb.Empty, error) {

	community, err := s.DBClient.GetCommunity(ctx, &dbpb.GetCommunityRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	_, err = s.DBClient.DeleteCommunity(ctx, &dbpb.DeleteCommunityRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	return &emptypb.Empty{}, nil
}

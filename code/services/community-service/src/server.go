package server

import (
	"community-service/src/pb"
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CommunityServer struct {
	pb.UnimplementedCommunityServiceServer
	DBClient pb.DBServiceClient
}

func (s *CommunityServer) ListCommunities(ctx context.Context, req *pb.ListCommunitiesRequest) (*pb.ListCommunitiesResponse, error) {
	res, err := s.DBClient.ListCommunities(ctx, &pb.ListCommunitiesRequest{
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

	communities := make([]*pb.Community, len(res.Communities))
	for i, community := range res.Communities {
		communities[i] = &pb.Community{
			Id:          community.Id,
			OwnerId:     community.OwnerId,
			Name:        community.Name,
			Description: community.Description,
			CreatedAt:   community.CreatedAt,
			UpdatedAt:   community.UpdatedAt,
		}
	}

	return &pb.ListCommunitiesResponse{
		Communities: communities,
		Pagination: &pb.Pagination{
			CurrentPage: res.Pagination.CurrentPage,
			PerPage:     res.Pagination.PerPage,
			TotalItems:  res.Pagination.TotalItems,
			TotalPages:  res.Pagination.TotalPages,
		},
	}, nil
}

func (s *CommunityServer) CreateCommunity(ctx context.Context, req *pb.CreateCommunityRequest) (*pb.Community, error) {
	userId, err := getCurrentUserId(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.DBClient.CreateCommunity(ctx, &pb.CreateCommunityRequest{
		OwnerId:     userId,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	return &pb.Community{
		Id:          res.Id,
		OwnerId:     res.OwnerId,
		Name:        res.Name,
		Description: res.Description,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}, nil
}

func (s *CommunityServer) GetCommunity(ctx context.Context, req *pb.GetCommunityRequest) (*pb.Community, error) {
	res, err := s.DBClient.GetCommunity(ctx, &pb.GetCommunityRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	return &pb.Community{
		Id:          res.Id,
		OwnerId:     res.OwnerId,
		Name:        res.Name,
		Description: res.Description,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}, nil
}

func (s *CommunityServer) UpdateCommunity(ctx context.Context, req *pb.UpdateCommunityRequest) (*pb.Community, error) {
	userId, err := getCurrentUserId(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.DBClient.UpdateCommunity(ctx, &pb.UpdateCommunityRequest{
		Id:          req.Id,
		OwnerId:     userId,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	return &pb.Community{
		Id:          res.Id,
		OwnerId:     res.OwnerId,
		Name:        res.Name,
		Description: res.Description,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}, nil
}

func (s *CommunityServer) DeleteCommunity(ctx context.Context, req *pb.DeleteCommunityRequest) (*emptypb.Empty, error) {
	userId, err := getCurrentUserId(ctx)
	if err != nil {
		return nil, err
	}

	_, err = s.DBClient.DeleteCommunity(ctx, &pb.DeleteCommunityRequest{
		Id:      req.Id,
		OwnerId: userId,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func getCurrentUserId(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("no metadata found in context")
	}
	userIds := md.Get("x-user-id")
	if len(userIds) == 0 {
		return "", fmt.Errorf("user id not found in metadata")
	}
	return userIds[0], nil
}

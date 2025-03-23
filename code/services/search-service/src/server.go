package server

import (
	"fmt"
	"context"
	pb "search-service/src/pb"
)

type SearchServer struct {
	pb.UnimplementedSearchServiceServer
	UserClient userpb.UserServiceClient
	CommunityClient communitypb.CommunityServiceClient
	ThreadClient threadpb.ThreadServiceClient
}

func (s *SearchServer) GlobalSearch(ctx context.Context, req *pb.GlobalSearchRequest) (*pb.GlobalSearchResponse, error) {
	userRes, err := s.UserClient.GetUsersByName(ctx, &userpb.UserQueryRequest{
		Query: req.Query,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling user service: %w", err)
	}
	communityRes, err := s.CommunityClient.GetCommunitiesByName(ctx, &communitypb.CommunityQueryRequest{
		Query: req.Query,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling community service: %w", err)
	}
	threadRes, err := s.ThreadClient.GetThreadsByTitle(ctx, &threadpb.ThreadQueryRequest{
		Query: req.Query,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling thread service: %w", err)
	}
	return &pb.GlobalSearchResponse{
		UserIds: userRes.UserIds,
		CommunityIds: communityRes.CommunityIds,
		ThreadIds: threadRes.ThreadIds,
	}, nil
}

func (s *SearchServer) UserSearch(ctx context.Context, req *pb.UserSearchRequest) (*pb.UserSearchResponse, error) {
	res, err := s.UserClient.GetUsersByName(ctx, &userpb.UserQueryRequest{
		Query: req.Query,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling user service: %w", err)
	}
	return &pb.UserSearchResponse{
		UserIds: res.UserIds,
	}, nil
}

func (s *SearchServer) CommunitySearch(ctx context.Context, req *pb.CommunitySearchRequest) (*pb.CommunitySearchResponse, error) {
	res, err := s.UserClient.GetCommunitiesByName(ctx, &communitypb.CommunityQueryRequest{
		Query: req.Query,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling community service: %w", err)
	}
	return &pb.CommunitySearchResponse{
		CommunityIds: res.CommunityIds,
	}, nil
}

func (s *SearchServer) ThreadSearch(ctx context.Context, req *pb.ThreadSearchRequest) (*pb.ThreadSearchResponse, error) {
	res, err := s.UserClient.GetThreadsByTitle(ctx, &threadpb.ThreadQueryRequest{
		Query: req.Query,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling thread service: %w", err)
	}
	return &pb.ThreadSearchResponse{
		ThreadIds: res.ThreadIds,
	}, nil
}
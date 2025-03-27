package server

import (
	"context"
	"fmt"
	"search-service/src/pb"
)

type SearchServer struct {
	pb.UnimplementedSearchServiceServer
	UserClient      pb.UserServiceClient
	CommunityClient pb.CommunityServiceClient
	ThreadClient    pb.ThreadServiceClient
}

func (s *SearchServer) GlobalSearch(ctx context.Context, req *pb.GlobalSearchRequest) (*pb.GlobalSearchResponse, error) {
	userResults, err := s.searchUsers(ctx, req.Query)
	if err != nil {
		return nil, err
	}
	communityResults, err := s.searchCommunities(ctx, req.Query)
	if err != nil {
		return nil, err
	}
	threadResults, err := s.searchThreads(ctx, req.Query)
	if err != nil {
		return nil, err
	}
	return &pb.GlobalSearchResponse{
		UserIds:      userResults,
		CommunityIds: communityResults,
		ThreadIds:    threadResults,
	}, nil
}

func (s *SearchServer) UserSearch(ctx context.Context, req *pb.UserSearchRequest) (*pb.UserSearchResponse, error) {
	results, err := s.searchUsers(ctx, req.Query)
	if err != nil {
		return nil, err
	}
	return &pb.UserSearchResponse{
		UserIds: results,
	}, nil
}

func (s *SearchServer) CommunitySearch(ctx context.Context, req *pb.CommunitySearchRequest) (*pb.CommunitySearchResponse, error) {
	results, err := s.searchCommunities(ctx, req.Query)
	if err != nil {
		return nil, err
	}
	return &pb.CommunitySearchResponse{
		CommunityIds: results,
	}, nil
}

func (s *SearchServer) ThreadSearch(ctx context.Context, req *pb.ThreadSearchRequest) (*pb.ThreadSearchResponse, error) {
	results, err := s.searchThreads(ctx, req.Query)
	if err != nil {
		return nil, err
	}
	return &pb.ThreadSearchResponse{
		ThreadIds: results,
	}, nil
}

func (s *SearchServer) searchUsers(ctx context.Context, query string) ([]string, error) {
	res, err := s.UserClient.GetUsersByName(ctx, &pb.UserQueryRequest{
		Query: query,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling user service: %w", err)
	}
	return res.UserIds, nil
}

func (s *SearchServer) searchCommunities(ctx context.Context, query string) ([]string, error) {
	res, err := s.CommunityClient.GetCommunitiesByName(ctx, &pb.CommunityQueryRequest{
		Query: query,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling community service: %w", err)
	}
	return res.CommunityIds, nil
}

func (s *SearchServer) searchThreads(ctx context.Context, query string) ([]string, error) {
	res, err := s.ThreadClient.GetThreadsByTitle(ctx, &pb.ThreadQueryRequest{
		Query: query,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling thread service: %w", err)
	}
	return res.ThreadIds, nil
}

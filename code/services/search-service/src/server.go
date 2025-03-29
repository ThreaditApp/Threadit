package server

import (
	"context"
	"fmt"
	communitypb "gen/community-service/pb"
	searchpb "gen/search-service/pb"
	threadpb "gen/thread-service/pb"
	userpb "gen/user-service/pb"
)

type SearchServer struct {
	searchpb.UnimplementedSearchServiceServer
	UserClient      userpb.UserServiceClient
	CommunityClient communitypb.CommunityServiceClient
	ThreadClient    threadpb.ThreadServiceClient
}

func (s *SearchServer) GlobalSearch(ctx context.Context, req *searchpb.GlobalSearchRequest) (*searchpb.GlobalSearchResponse, error) {
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
	return &searchpb.GlobalSearchResponse{
		UserIds:      userResults,
		CommunityIds: communityResults,
		ThreadIds:    threadResults,
	}, nil
}

func (s *SearchServer) UserSearch(ctx context.Context, req *searchpb.UserSearchRequest) (*searchpb.UserSearchResponse, error) {
	results, err := s.searchUsers(ctx, req.Query)
	if err != nil {
		return nil, err
	}
	return &searchpb.UserSearchResponse{
		UserIds: results,
	}, nil
}

func (s *SearchServer) CommunitySearch(ctx context.Context, req *searchpb.CommunitySearchRequest) (*searchpb.CommunitySearchResponse, error) {
	results, err := s.searchCommunities(ctx, req.Query)
	if err != nil {
		return nil, err
	}
	return &searchpb.CommunitySearchResponse{
		CommunityIds: results,
	}, nil
}

func (s *SearchServer) ThreadSearch(ctx context.Context, req *searchpb.ThreadSearchRequest) (*searchpb.ThreadSearchResponse, error) {
	results, err := s.searchThreads(ctx, req.Query)
	if err != nil {
		return nil, err
	}
	return &searchpb.ThreadSearchResponse{
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

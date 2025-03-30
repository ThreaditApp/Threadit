package server

import (
	"context"
	communitypb "gen/community-service/pb"
	models "gen/models/pb"
	searchpb "gen/search-service/pb"
	threadpb "gen/thread-service/pb"
)

type SearchServer struct {
	searchpb.UnimplementedSearchServiceServer
	CommunityClient communitypb.CommunityServiceClient
	ThreadClient    threadpb.ThreadServiceClient
}

func (s *SearchServer) GlobalSearch(ctx context.Context, req *searchpb.SearchRequest) (*searchpb.GlobalSearchResponse, error) {
	communityResults, err := s.searchCommunities(ctx, &req.Query)
	if err != nil {
		return nil, err
	}
	threadResults, err := s.searchThreads(ctx, &req.Query)
	if err != nil {
		return nil, err
	}
	return &searchpb.GlobalSearchResponse{
		CommunityResults: communityResults,
		ThreadResults:    threadResults,
	}, nil
}

func (s *SearchServer) CommunitySearch(ctx context.Context, req *searchpb.SearchRequest) (*searchpb.CommunitySearchResponse, error) {
	results, err := s.searchCommunities(ctx, &req.Query)
	if err != nil {
		return nil, err
	}
	return &searchpb.CommunitySearchResponse{
		Results: results,
	}, nil
}

func (s *SearchServer) ThreadSearch(ctx context.Context, req *searchpb.SearchRequest) (*searchpb.ThreadSearchResponse, error) {
	results, err := s.searchThreads(ctx, &req.Query)
	if err != nil {
		return nil, err
	}
	return &searchpb.ThreadSearchResponse{
		Results: results,
	}, nil
}

func (s *SearchServer) searchCommunities(ctx context.Context, query *string) ([]*models.Community, error) {
	res, err := s.CommunityClient.ListCommunities(ctx, &communitypb.ListCommunitiesRequest{
		Name: query,
	})
	if err != nil {
		return nil, err
	}
	return res.Communities, nil
}

func (s *SearchServer) searchThreads(ctx context.Context, query *string) ([]*models.Thread, error) {
	res, err := s.ThreadClient.ListThreads(ctx, &threadpb.ListThreadsRequest{
		Title: query,
	})
	if err != nil {
		return nil, err
	}
	return res.Threads, nil
}

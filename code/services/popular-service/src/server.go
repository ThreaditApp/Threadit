package server

import (
	"context"
	"fmt"
	commentpb "gen/comment-service/pb"
	popularpb "gen/popular-service/pb"
	threadpb "gen/thread-service/pb"
)

type PopularServer struct {
	popularpb.UnimplementedPopularServiceServer
	ThreadClient  threadpb.ThreadServiceClient
	CommentClient commentpb.CommentServiceClient
}

func (s *PopularServer) GetPopularThreads(ctx context.Context, req *popularpb.GetPopularThreadsRequest) (*popularpb.GetPopularThreadsResponse, error) {
	sortBy := "votes"
	res, err := s.ThreadClient.ListThreads(ctx, &threadpb.ListThreadsRequest{
		Offset: req.Offset,
		Limit:  req.Limit,
		SortBy: &sortBy,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling thread service: %w", err)
	}
	return &popularpb.GetPopularThreadsResponse{
		Threads: res.Threads,
	}, nil
}

func (s *PopularServer) GetPopularComments(ctx context.Context, req *popularpb.GetPopularCommentsRequest) (*popularpb.GetPopularCommentsResponse, error) {
	sortBy := "votes"
	res, err := s.CommentClient.ListComments(ctx, &commentpb.ListCommentsRequest{
		Offset:   req.Offset,
		Limit:    req.Limit,
		SortBy:   &sortBy,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling comment service: %w", err)
	}
	return &popularpb.GetPopularCommentsResponse{
		Comments: res.Comments,
	}, nil
}
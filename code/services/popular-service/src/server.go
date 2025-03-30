package server

import (
	"context"
	commentpb "gen/comment-service/pb"
	popularpb "gen/popular-service/pb"
	threadpb "gen/thread-service/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PopularServer struct {
	popularpb.UnimplementedPopularServiceServer
	ThreadClient  threadpb.ThreadServiceClient
	CommentClient commentpb.CommentServiceClient
}

func (s *PopularServer) GetPopularThreads(ctx context.Context, req *popularpb.GetPopularThreadsRequest) (*popularpb.GetPopularThreadsResponse, error) {
	// input validation
	if req.Offset != nil && *req.Offset < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "offset must be a non-negative integer")
	}
	if req.Limit != nil && *req.Limit <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "limit must be a positive integer")
	}

	sortBy := "votes"
	res, err := s.ThreadClient.ListThreads(ctx, &threadpb.ListThreadsRequest{
		Offset: req.Offset,
		Limit:  req.Limit,
		SortBy: &sortBy,
	})
	if err != nil {
		return nil, err
	}
	return &popularpb.GetPopularThreadsResponse{
		Threads: res.Threads,
	}, nil
}

func (s *PopularServer) GetPopularComments(ctx context.Context, req *popularpb.GetPopularCommentsRequest) (*popularpb.GetPopularCommentsResponse, error) {
	// input validation
	if req.Offset != nil && *req.Offset < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "offset must be a non-negative integer")
	}
	if req.Limit != nil && *req.Limit <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "limit must be a positive integer")
	}

	sortBy := "votes"
	res, err := s.CommentClient.ListComments(ctx, &commentpb.ListCommentsRequest{
		Offset: req.Offset,
		Limit:  req.Limit,
		SortBy: &sortBy,
	})
	if err != nil {
		return nil, err
	}
	return &popularpb.GetPopularCommentsResponse{
		Comments: res.Comments,
	}, nil
}

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
	// validate inputs
	if req.Offset != nil && req.GetOffset() < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "offset must be a non-negative integer")
	}
	if req.Limit != nil && req.GetLimit() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "limit must be a positive integer")
	}

	// fetch threads
	sortBy := "ups" // upvotes
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
	// validate inputs
	if req.Offset != nil && req.GetOffset() < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "offset must be a non-negative integer")
	}
	if req.Limit != nil && req.GetLimit() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "limit must be a positive integer")
	}

	// fetch comments
	sortBy := "ups" // upvotes
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

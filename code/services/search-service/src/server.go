package server

import (
	"context"
	pb "search-service/src/pb"
)

type SearchServer struct {
	pb.UnimplementedSearchServiceServer
}

func (s *SearchServer) GlobalSearch(ctx context.Context, req *pb.GlobalSearchRequest) (*pb.GlobalSearchResponse, error) {
	return &pb.GlobalSearchResponse{}, nil
}

func (s *SearchServer) UserSearch(ctx context.Context, req *pb.UserSearchRequest) (*pb.UserSearchResponse, error) {
	return &pb.UserSearchResponse{}, nil
}

func (s *SearchServer) CommunitySearch(ctx context.Context, req *pb.CommunitySearchRequest) (*pb.CommunitySearchResponse, error) {
	return &pb.CommunitySearchResponse{}, nil
}

func (s *SearchServer) ThreadSearch(ctx context.Context, req *pb.ThreadSearchRequest) (*pb.ThreadSearchResponse, error) {
	return &pb.ThreadSearchResponse{}, nil
}
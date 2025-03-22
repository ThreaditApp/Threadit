package server

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "social-service/src/pb"
)

type SocialServer struct {
	pb.UnimplementedSocialServiceServer
}

func (s *SocialServer) FollowUser(ctx context.Context, req *pb.FollowUserRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *SocialServer) UnfollowUser(ctx context.Context, req *pb.FollowUserRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *SocialServer) GetFollowers(ctx context.Context, req *pb.GetFollowersRequest) (*pb.GetFollowersResponse, error) {
	return &pb.GetFollowersResponse{}, nil
}

func (s *SocialServer) GetFollowing(ctx context.Context, req *pb.GetFollowingRequest) (*pb.GetFollowingResponse, error) {
	return &pb.GetFollowingResponse{}, nil
}

func (s *SocialServer) FollowCommunity(ctx context.Context, req *pb.FollowCommunityRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *SocialServer) UnfollowCommunity(ctx context.Context, req *pb.FollowCommunityRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *SocialServer) GetCommunityFollowers(ctx context.Context, req *pb.GetCommunityFollowersRequest) (*pb.GetCommunityFollowersResponse, error) {
	return &pb.GetCommunityFollowersResponse{}, nil
}
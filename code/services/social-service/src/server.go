package server

import (
	"context"
	"fmt"
	dbpb "gen/db-service/pb"
	socialpb "gen/social-service/pb"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SocialServer struct {
	socialpb.UnimplementedSocialServiceServer
	DBClient dbpb.DBServiceClient
}

func (s *SocialServer) FollowUser(ctx context.Context, req *socialpb.FollowUserRequest) (*emptypb.Empty, error) {
	user, _ := getCurrentUserId(ctx)
	_, err := s.DBClient.FollowUser(ctx, &dbpb.FollowUserRequest{
		UserId:       user,
		TargetUserId: req.UserId,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *SocialServer) UnfollowUser(ctx context.Context, req *socialpb.FollowUserRequest) (*emptypb.Empty, error) {
	_, err := s.DBClient.UnfollowUser(ctx, &dbpb.UnfollowUserRequest{
		UserId:       getCurrentUserId(ctx),
		TargetUserId: req.UserId,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *SocialServer) GetFollowers(ctx context.Context, req *socialpb.GetFollowersRequest) (*socialpb.GetFollowersResponse, error) {
	res, err := s.DBClient.GetFollowers(ctx, &dbpb.GetFollowersRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}
	return &socialpb.GetFollowersResponse{
		UserIds: res.UserIds,
	}, nil
}

func (s *SocialServer) GetFollowing(ctx context.Context, req *socialpb.GetFollowingRequest) (*socialpb.GetFollowingResponse, error) {
	res, err := s.DBClient.GetFollowing(ctx, &dbpb.GetFollowingRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}
	return &socialpb.GetFollowingResponse{
		UserIds: res.UserIds,
	}, nil

}

func (s *SocialServer) FollowCommunity(ctx context.Context, req *socialpb.FollowCommunityRequest) (*emptypb.Empty, error) {
	_, err := s.DBClient.FollowCommunity(ctx, &dbpb.FollowCommunityRequest{
		UserId:      getCurrentUserId(ctx),
		CommunityId: req.CommunityId,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *SocialServer) UnfollowCommunity(ctx context.Context, req *socialpb.FollowCommunityRequest) (*emptypb.Empty, error) {
	_, err := s.DBClient.UnfollowCommunity(ctx, &dbpb.UnfollowCommunityRequest{
		UserId:      getCurrentUserId(ctx),
		CommunityId: req.CommunityId,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *SocialServer) GetCommunityFollowers(ctx context.Context, req *socialpb.GetCommunityFollowersRequest) (*socialpb.GetCommunityFollowersResponse, error) {
	res, err := s.DBClient.GetCommunityFollowers(ctx, &dbpb.GetCommunityFollowersRequest{
		CommunityId: req.CommunityId,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}
	return &socialpb.GetCommunityFollowersResponse{
		UserIds: res.UserIds,
	}, nil
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

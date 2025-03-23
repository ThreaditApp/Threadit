package server

import (
	"fmt"
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "vote-service/src/pb"
)

type VoteServer struct {
	pb.UnimplementedVoteServiceServer
	ThreadClient threadpb.ThreadServiceClient
	CommentClient commentpb.CommentServiceClient
}

func (s *VoteServer) UpvoteThread(ctx context.Context, req *pb.VoteThreadRequest) (*emptypb.Empty, error) {
	_, err := s.ThreadClient.UpdateVote(ctx, &threadpb.VoteThreadRequest{
		ThreadId: req.ThreadId,
		UserId: getCurrentUserId(ctx),
		Value: 1,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *VoteServer) DownvoteThread(ctx context.Context, req *pb.VoteThreadRequest) (*emptypb.Empty, error) {
	_, err := s.ThreadClient.UpdateVote(ctx, &threadpb.VoteThreadRequest{
		ThreadId: req.ThreadId,
		UserId: getCurrentUserId(ctx),
		Value: -1,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *VoteServer) RemoveThreadVote(ctx context.Context, req *pb.VoteThreadRequest) (*emptypb.Empty, error) {
	_, err := s.ThreadClient.UpdateVote(ctx, &threadpb.VoteThreadRequest{
		ThreadId: req.ThreadId,
		UserId: getCurrentUserId(ctx),
		Value: 0,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *VoteServer) UpvoteComment(ctx context.Context, req *pb.VoteCommentRequest) (*emptypb.Empty, error) {
	_, err := s.CommentClient.UpdateVote(ctx, &commentpb.VoteCommentRequest{
		CommentId: req.CommentId,
		UserId: getCurrentUserId(ctx),
		Value: 1,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *VoteServer) DownvoteComment(ctx context.Context, req *pb.VoteCommentRequest) (*emptypb.Empty, error) {
	_, err := s.CommentClient.UpdateVote(ctx, &commentpb.VoteCommentRequest{
		CommentId: req.CommentId,
		UserId: getCurrentUserId(ctx),
		Value: -1,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *VoteServer) RemoveCommentVote(ctx context.Context, req *pb.VoteCommentRequest) (*emptypb.Empty, error) {
	_, err := s.CommentClient.UpdateVote(ctx, &commentpb.VoteCommentRequest{
		CommentId: req.CommentId,
		UserId: getCurrentUserId(ctx),
		Value: 0,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}
	return &emptypb.Empty{}, nil
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

package server

import (
	"context"
	"fmt"
	commentpb "gen/comment-service/pb"
	threadpb "gen/thread-service/pb"
	votepb "gen/vote-service/pb"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type VoteServer struct {
	votepb.UnimplementedVoteServiceServer
	ThreadClient  threadpb.ThreadServiceClient
	CommentClient commentpb.CommentServiceClient
}

func (s *VoteServer) UpvoteThread(ctx context.Context, req *votepb.VoteThreadRequest) (*emptypb.Empty, error) {
	return s.updateThreadVote(ctx, req, 1)
}

func (s *VoteServer) DownvoteThread(ctx context.Context, req *votepb.VoteThreadRequest) (*emptypb.Empty, error) {
	return s.updateThreadVote(ctx, req, -1)
}

func (s *VoteServer) RemoveThreadVote(ctx context.Context, req *votepb.VoteThreadRequest) (*emptypb.Empty, error) {
	return s.updateThreadVote(ctx, req, 0)
}

func (s *VoteServer) UpvoteComment(ctx context.Context, req *votepb.VoteCommentRequest) (*emptypb.Empty, error) {
	return s.updateCommentVote(ctx, req, 1)
}

func (s *VoteServer) DownvoteComment(ctx context.Context, req *votepb.VoteCommentRequest) (*emptypb.Empty, error) {
	return s.updateCommentVote(ctx, req, -1)
}

func (s *VoteServer) RemoveCommentVote(ctx context.Context, req *votepb.VoteCommentRequest) (*emptypb.Empty, error) {
	return s.updateCommentVote(ctx, req, 0)
}

func (s *VoteServer) updateThreadVote(ctx context.Context, req *votepb.VoteThreadRequest, value int) (*emptypb.Empty, error) {
	_, err := s.ThreadClient.UpdateVote(ctx, &threadpb.VoteThreadRequest{
		ThreadId: req.ThreadId,
		UserId:   getCurrentUserId(ctx),
		Value:    value,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *VoteServer) updateCommentVote(ctx context.Context, req *pb.VoteCommentRequest, value int) (*emptypb.Empty, error) {
	_, err := s.CommentClient.UpdateVote(ctx, &commentpb.VoteCommentRequest{
		CommentId: req.CommentId,
		UserId:    getCurrentUserId(ctx),
		Value:     value,
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

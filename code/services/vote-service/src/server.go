package server

import (
	"context"
	"fmt"
	commentpb "gen/comment-service/pb"
	threadpb "gen/thread-service/pb"
	votepb "gen/vote-service/pb"
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

func (s *VoteServer) updateThreadVote(ctx context.Context, req *votepb.VoteThreadRequest, value int32) (*emptypb.Empty, error) {
	_, err := s.ThreadClient.UpdateThread(ctx, &threadpb.UpdateThreadRequest{
		Id:         req.ThreadId,
		VoteOffset: &value,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *VoteServer) updateCommentVote(ctx context.Context, req *votepb.VoteCommentRequest, value int32) (*emptypb.Empty, error) {
	_, err := s.CommentClient.UpdateComment(ctx, &commentpb.UpdateCommentRequest{
		Id:         req.CommentId,
		VoteOffset: &value,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}
	return &emptypb.Empty{}, nil
}

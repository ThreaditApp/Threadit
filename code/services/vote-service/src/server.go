package server

import (
	"context"
	commentpb "gen/comment-service/pb"
	threadpb "gen/thread-service/pb"
	votepb "gen/vote-service/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type VoteServer struct {
	votepb.UnimplementedVoteServiceServer
	ThreadClient  threadpb.ThreadServiceClient
	CommentClient commentpb.CommentServiceClient
}

func (s *VoteServer) CheckHealth(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *VoteServer) UpvoteThread(ctx context.Context, req *votepb.VoteThreadRequest) (*emptypb.Empty, error) {
	return s.updateThreadVote(ctx, req, 1)
}

func (s *VoteServer) DownvoteThread(ctx context.Context, req *votepb.VoteThreadRequest) (*emptypb.Empty, error) {
	return s.updateThreadVote(ctx, req, -1)
}

func (s *VoteServer) UpvoteComment(ctx context.Context, req *votepb.VoteCommentRequest) (*emptypb.Empty, error) {
	return s.updateCommentVote(ctx, req, 1)
}

func (s *VoteServer) DownvoteComment(ctx context.Context, req *votepb.VoteCommentRequest) (*emptypb.Empty, error) {
	return s.updateCommentVote(ctx, req, -1)
}

func (s *VoteServer) updateThreadVote(ctx context.Context, req *votepb.VoteThreadRequest, value int32) (*emptypb.Empty, error) {
	if req.GetThreadId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Thread id is required")
	}
	_, err := s.ThreadClient.UpdateThread(ctx, &threadpb.UpdateThreadRequest{
		Id:         req.ThreadId,
		VoteOffset: &value,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *VoteServer) updateCommentVote(ctx context.Context, req *votepb.VoteCommentRequest, value int32) (*emptypb.Empty, error) {
	if req.GetCommentId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Comment id is required")
	}
	_, err := s.CommentClient.UpdateComment(ctx, &commentpb.UpdateCommentRequest{
		Id:         req.CommentId,
		VoteOffset: &value,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

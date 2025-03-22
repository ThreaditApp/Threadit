package server

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "vote-service/src/pb"
)

type VoteServer struct {
	pb.UnimplementedVoteServiceServer
}

func (s *VoteServer) UpvoteThread(ctx context.Context, req *pb.VoteThreadRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *VoteServer) DownvoteThread(ctx context.Context, req *pb.VoteThreadRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *VoteServer) RemoveThreadVote(ctx context.Context, req *pb.VoteThreadRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *VoteServer) UpvoteComment(ctx context.Context, req *pb.VoteCommentRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *VoteServer) DownvoteComment(ctx context.Context, req *pb.VoteCommentRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *VoteServer) RemoveCommentVote(ctx context.Context, req *pb.VoteCommentRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

package test

import (
	"context"
	"testing"

	commentpb "gen/comment-service/pb"
	threadpb "gen/thread-service/pb"
	votepb "gen/vote-service/pb"
	src "vote-service/src"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/grpc"
)

type MockThreadClient struct {
	threadpb.ThreadServiceClient
	UpdateThreadFunc func(ctx context.Context, req *threadpb.UpdateThreadRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

func (m *MockThreadClient) UpdateThread(ctx context.Context, req *threadpb.UpdateThreadRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return m.UpdateThreadFunc(ctx, req, opts...)
}

type MockCommentClient struct {
	commentpb.CommentServiceClient
	UpdateCommentFunc func(ctx context.Context, req *commentpb.UpdateCommentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

func (m *MockCommentClient) UpdateComment(ctx context.Context, req *commentpb.UpdateCommentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return m.UpdateCommentFunc(ctx, req, opts...)
}

func TestUpvoteThread_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     *votepb.VoteThreadRequest
		wantErr error
	}{
		{
			name:    "missing thread id",
			req:     &votepb.VoteThreadRequest{},
			wantErr: status.Error(codes.InvalidArgument, "Thread id is required"),
		},
		{
			name: "valid request",
			req: &votepb.VoteThreadRequest{
				ThreadId: "123",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &src.VoteServer{
				ThreadClient: &MockThreadClient{
					UpdateThreadFunc: func(ctx context.Context, req *threadpb.UpdateThreadRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
						return &emptypb.Empty{}, nil
					},
				},
				CommentClient: &MockCommentClient{},
			}

			_, err := server.UpvoteThread(context.Background(), tt.req)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpvoteComment_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     *votepb.VoteCommentRequest
		wantErr error
	}{
		{
			name:    "missing comment id",
			req:     &votepb.VoteCommentRequest{},
			wantErr: status.Error(codes.InvalidArgument, "Comment id is required"),
		},
		{
			name: "valid request",
			req: &votepb.VoteCommentRequest{
				CommentId: "123",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &src.VoteServer{
				ThreadClient: &MockThreadClient{
					UpdateThreadFunc: func(ctx context.Context, req *threadpb.UpdateThreadRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
						return &emptypb.Empty{}, nil
					},
				},
				CommentClient: &MockCommentClient{
					UpdateCommentFunc: func(ctx context.Context, req *commentpb.UpdateCommentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
						return &emptypb.Empty{}, nil
					},
				},
			}

			_, err := server.UpvoteComment(context.Background(), tt.req)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

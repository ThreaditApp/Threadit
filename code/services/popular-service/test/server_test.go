package test

import (
	"context"
	"testing"

	commentpb "gen/comment-service/pb"
	models "gen/models/pb"
	popularpb "gen/popular-service/pb"
	threadpb "gen/thread-service/pb"
	src "popular-service/src"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc"
)

type MockThreadClient struct {
	threadpb.ThreadServiceClient
	ListThreadsFunc func(ctx context.Context, req *threadpb.ListThreadsRequest, opts ...grpc.CallOption) (*threadpb.ListThreadsResponse, error)
}

func (m *MockThreadClient) ListThreads(ctx context.Context, req *threadpb.ListThreadsRequest, opts ...grpc.CallOption) (*threadpb.ListThreadsResponse, error) {
	return m.ListThreadsFunc(ctx, req, opts...)
}

type MockCommentClient struct {
	commentpb.CommentServiceClient
	ListCommentsFunc func(ctx context.Context, req *commentpb.ListCommentsRequest, opts ...grpc.CallOption) (*commentpb.ListCommentsResponse, error)
}

func (m *MockCommentClient) ListComments(ctx context.Context, req *commentpb.ListCommentsRequest, opts ...grpc.CallOption) (*commentpb.ListCommentsResponse, error) {
	return m.ListCommentsFunc(ctx, req, opts...)
}

func TestGetPopularThreads_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     *popularpb.GetPopularThreadsRequest
		wantErr error
	}{
		{
			name: "negative offset",
			req: &popularpb.GetPopularThreadsRequest{
				Offset: int32Ptr(-1),
			},
			wantErr: status.Error(codes.InvalidArgument, "Offset must be a positive integer"),
		},
		{
			name: "zero limit",
			req: &popularpb.GetPopularThreadsRequest{
				Limit: int32Ptr(0),
			},
			wantErr: status.Error(codes.InvalidArgument, "Limit must be a positive integer"),
		},
		{
			name: "negative limit",
			req: &popularpb.GetPopularThreadsRequest{
				Limit: int32Ptr(-1),
			},
			wantErr: status.Error(codes.InvalidArgument, "Limit must be a positive integer"),
		},
		{
			name: "valid request",
			req: &popularpb.GetPopularThreadsRequest{
				Offset: int32Ptr(0),
				Limit:  int32Ptr(10),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &src.PopularServer{
				ThreadClient: &MockThreadClient{
					ListThreadsFunc: func(ctx context.Context, req *threadpb.ListThreadsRequest, opts ...grpc.CallOption) (*threadpb.ListThreadsResponse, error) {
						return &threadpb.ListThreadsResponse{
							Threads: []*models.Thread{},
						}, nil
					},
				},
				CommentClient: &MockCommentClient{},
			}

			_, err := server.GetPopularThreads(context.Background(), tt.req)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetPopularComments_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     *popularpb.GetPopularCommentsRequest
		wantErr error
	}{
		{
			name: "negative offset",
			req: &popularpb.GetPopularCommentsRequest{
				Offset: int32Ptr(-1),
			},
			wantErr: status.Error(codes.InvalidArgument, "Offset must be a positive integer"),
		},
		{
			name: "zero limit",
			req: &popularpb.GetPopularCommentsRequest{
				Limit: int32Ptr(0),
			},
			wantErr: status.Error(codes.InvalidArgument, "Limit must be a positive integer"),
		},
		{
			name: "negative limit",
			req: &popularpb.GetPopularCommentsRequest{
				Limit: int32Ptr(-1),
			},
			wantErr: status.Error(codes.InvalidArgument, "Limit must be a positive integer"),
		},
		{
			name: "valid request",
			req: &popularpb.GetPopularCommentsRequest{
				Offset: int32Ptr(0),
				Limit:  int32Ptr(10),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &src.PopularServer{
				ThreadClient: &MockThreadClient{},
				CommentClient: &MockCommentClient{
					ListCommentsFunc: func(ctx context.Context, req *commentpb.ListCommentsRequest, opts ...grpc.CallOption) (*commentpb.ListCommentsResponse, error) {
						return &commentpb.ListCommentsResponse{
							Comments: []*models.Comment{},
						}, nil
					},
				},
			}

			_, err := server.GetPopularComments(context.Background(), tt.req)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func int32Ptr(i int32) *int32 {
	return &i
}

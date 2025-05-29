package test

import (
	"context"
	"testing"

	communitypb "gen/community-service/pb"
	models "gen/models/pb"
	searchpb "gen/search-service/pb"
	threadpb "gen/thread-service/pb"
	src "search-service/src"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc"
)

type MockCommunityClient struct {
	communitypb.CommunityServiceClient
	ListCommunitiesFunc func(ctx context.Context, req *communitypb.ListCommunitiesRequest, opts ...grpc.CallOption) (*communitypb.ListCommunitiesResponse, error)
}

func (m *MockCommunityClient) ListCommunities(ctx context.Context, req *communitypb.ListCommunitiesRequest, opts ...grpc.CallOption) (*communitypb.ListCommunitiesResponse, error) {
	return m.ListCommunitiesFunc(ctx, req, opts...)
}

type MockThreadClient struct {
	threadpb.ThreadServiceClient
	ListThreadsFunc func(ctx context.Context, req *threadpb.ListThreadsRequest, opts ...grpc.CallOption) (*threadpb.ListThreadsResponse, error)
}

func (m *MockThreadClient) ListThreads(ctx context.Context, req *threadpb.ListThreadsRequest, opts ...grpc.CallOption) (*threadpb.ListThreadsResponse, error) {
	return m.ListThreadsFunc(ctx, req, opts...)
}

func TestGlobalSearch_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     *searchpb.SearchRequest
		wantErr error
	}{
		{
			name:    "empty query",
			req:     &searchpb.SearchRequest{Query: ""},
			wantErr: status.Error(codes.InvalidArgument, "Query is empty"),
		},
		{
			name: "negative offset",
			req: &searchpb.SearchRequest{
				Query:  "test",
				Offset: int32Ptr(-1),
			},
			wantErr: status.Error(codes.InvalidArgument, "Offset cannot be negative"),
		},
		{
			name: "negative limit",
			req: &searchpb.SearchRequest{
				Query: "test",
				Limit: int32Ptr(-1),
			},
			wantErr: status.Error(codes.InvalidArgument, "Limit cannot be negative"),
		},
		{
			name: "valid request",
			req: &searchpb.SearchRequest{
				Query:  "test",
				Offset: int32Ptr(0),
				Limit:  int32Ptr(10),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &src.SearchServer{
				CommunityClient: &MockCommunityClient{
					ListCommunitiesFunc: func(ctx context.Context, req *communitypb.ListCommunitiesRequest, opts ...grpc.CallOption) (*communitypb.ListCommunitiesResponse, error) {
						return &communitypb.ListCommunitiesResponse{
							Communities: []*models.Community{},
						}, nil
					},
				},
				ThreadClient: &MockThreadClient{
					ListThreadsFunc: func(ctx context.Context, req *threadpb.ListThreadsRequest, opts ...grpc.CallOption) (*threadpb.ListThreadsResponse, error) {
						return &threadpb.ListThreadsResponse{
							Threads: []*models.Thread{},
						}, nil
					},
				},
			}

			_, err := server.GlobalSearch(context.Background(), tt.req)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCommunitySearch_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     *searchpb.SearchRequest
		wantErr error
	}{
		{
			name:    "empty query",
			req:     &searchpb.SearchRequest{Query: ""},
			wantErr: status.Error(codes.InvalidArgument, "Query is empty"),
		},
		{
			name: "negative offset",
			req: &searchpb.SearchRequest{
				Query:  "test",
				Offset: int32Ptr(-1),
			},
			wantErr: status.Error(codes.InvalidArgument, "Offset cannot be negative"),
		},
		{
			name: "negative limit",
			req: &searchpb.SearchRequest{
				Query: "test",
				Limit: int32Ptr(-1),
			},
			wantErr: status.Error(codes.InvalidArgument, "Limit cannot be negative"),
		},
		{
			name: "valid request",
			req: &searchpb.SearchRequest{
				Query:  "test",
				Offset: int32Ptr(0),
				Limit:  int32Ptr(10),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &src.SearchServer{
				CommunityClient: &MockCommunityClient{
					ListCommunitiesFunc: func(ctx context.Context, req *communitypb.ListCommunitiesRequest, opts ...grpc.CallOption) (*communitypb.ListCommunitiesResponse, error) {
						return &communitypb.ListCommunitiesResponse{
							Communities: []*models.Community{},
						}, nil
					},
				},
			}

			_, err := server.CommunitySearch(context.Background(), tt.req)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestThreadSearch_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     *searchpb.SearchRequest
		wantErr error
	}{
		{
			name:    "empty query",
			req:     &searchpb.SearchRequest{Query: ""},
			wantErr: status.Error(codes.InvalidArgument, "Query is empty"),
		},
		{
			name: "negative offset",
			req: &searchpb.SearchRequest{
				Query:  "test",
				Offset: int32Ptr(-1),
			},
			wantErr: status.Error(codes.InvalidArgument, "Offset cannot be negative"),
		},
		{
			name: "negative limit",
			req: &searchpb.SearchRequest{
				Query: "test",
				Limit: int32Ptr(-1),
			},
			wantErr: status.Error(codes.InvalidArgument, "Limit cannot be negative"),
		},
		{
			name: "valid request",
			req: &searchpb.SearchRequest{
				Query:  "test",
				Offset: int32Ptr(0),
				Limit:  int32Ptr(10),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &src.SearchServer{
				ThreadClient: &MockThreadClient{
					ListThreadsFunc: func(ctx context.Context, req *threadpb.ListThreadsRequest, opts ...grpc.CallOption) (*threadpb.ListThreadsResponse, error) {
						return &threadpb.ListThreadsResponse{
							Threads: []*models.Thread{},
						}, nil
					},
				},
			}

			_, err := server.ThreadSearch(context.Background(), tt.req)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func int32Ptr(i int32) *int32 { return &i }

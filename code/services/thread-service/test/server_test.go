package test

import (
	"context"
	"testing"

	communitypb "gen/community-service/pb"
	dbpb "gen/db-service/pb"
	models "gen/models/pb"
	threadpb "gen/thread-service/pb"
	src "thread-service/src"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/grpc"
)

type MockDBClient struct {
	dbpb.DBServiceClient
	ListThreadsFunc  func(ctx context.Context, req *dbpb.ListThreadsRequest, opts ...grpc.CallOption) (*dbpb.ListThreadsResponse, error)
	CreateThreadFunc func(ctx context.Context, req *dbpb.CreateThreadRequest, opts ...grpc.CallOption) (*dbpb.CreateThreadResponse, error)
	GetThreadFunc    func(ctx context.Context, req *dbpb.GetThreadRequest, opts ...grpc.CallOption) (*models.Thread, error)
	UpdateThreadFunc func(ctx context.Context, req *dbpb.UpdateThreadRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteThreadFunc func(ctx context.Context, req *dbpb.DeleteThreadRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

func (m *MockDBClient) ListThreads(ctx context.Context, req *dbpb.ListThreadsRequest, opts ...grpc.CallOption) (*dbpb.ListThreadsResponse, error) {
	return m.ListThreadsFunc(ctx, req, opts...)
}

func (m *MockDBClient) CreateThread(ctx context.Context, req *dbpb.CreateThreadRequest, opts ...grpc.CallOption) (*dbpb.CreateThreadResponse, error) {
	return m.CreateThreadFunc(ctx, req, opts...)
}

func (m *MockDBClient) GetThread(ctx context.Context, req *dbpb.GetThreadRequest, opts ...grpc.CallOption) (*models.Thread, error) {
	return m.GetThreadFunc(ctx, req, opts...)
}

func (m *MockDBClient) UpdateThread(ctx context.Context, req *dbpb.UpdateThreadRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return m.UpdateThreadFunc(ctx, req, opts...)
}

func (m *MockDBClient) DeleteThread(ctx context.Context, req *dbpb.DeleteThreadRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return m.DeleteThreadFunc(ctx, req, opts...)
}

type MockCommunityClient struct {
	communitypb.CommunityServiceClient
	GetCommunityFunc    func(ctx context.Context, req *communitypb.GetCommunityRequest, opts ...grpc.CallOption) (*models.Community, error)
	UpdateCommunityFunc func(ctx context.Context, req *communitypb.UpdateCommunityRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

func (m *MockCommunityClient) GetCommunity(ctx context.Context, req *communitypb.GetCommunityRequest, opts ...grpc.CallOption) (*models.Community, error) {
	return m.GetCommunityFunc(ctx, req, opts...)
}

func (m *MockCommunityClient) UpdateCommunity(ctx context.Context, req *communitypb.UpdateCommunityRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return m.UpdateCommunityFunc(ctx, req, opts...)
}

func TestCreateThread_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     *threadpb.CreateThreadRequest
		wantErr error
	}{
		{
			name:    "missing community id",
			req:     &threadpb.CreateThreadRequest{},
			wantErr: status.Error(codes.InvalidArgument, "Community id is required"),
		},
		{
			name: "missing title",
			req: &threadpb.CreateThreadRequest{
				CommunityId: "123",
			},
			wantErr: status.Error(codes.InvalidArgument, "Title is required"),
		},
		{
			name: "missing content",
			req: &threadpb.CreateThreadRequest{
				CommunityId: "123",
				Title:       "test thread",
			},
			wantErr: status.Error(codes.InvalidArgument, "Content is required"),
		},
		{
			name: "valid request",
			req: &threadpb.CreateThreadRequest{
				CommunityId: "123",
				Title:       "test thread",
				Content:     "test content",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &src.ThreadServer{
				DBClient: &MockDBClient{
					CreateThreadFunc: func(ctx context.Context, req *dbpb.CreateThreadRequest, opts ...grpc.CallOption) (*dbpb.CreateThreadResponse, error) {
						return &dbpb.CreateThreadResponse{Id: "123"}, nil
					},
				},
				CommunityClient: &MockCommunityClient{
					GetCommunityFunc: func(ctx context.Context, req *communitypb.GetCommunityRequest, opts ...grpc.CallOption) (*models.Community, error) {
						return &models.Community{
							Id:   "123",
							Name: "test-community",
						}, nil
					},
					UpdateCommunityFunc: func(ctx context.Context, req *communitypb.UpdateCommunityRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
						return &emptypb.Empty{}, nil
					},
				},
			}
			_, err := server.CreateThread(context.Background(), tt.req)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetThread_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     *threadpb.GetThreadRequest
		wantErr error
	}{
		{
			name:    "missing id",
			req:     &threadpb.GetThreadRequest{},
			wantErr: status.Error(codes.InvalidArgument, "Thread id is required"),
		},
		{
			name: "valid request",
			req: &threadpb.GetThreadRequest{
				Id: "123",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &src.ThreadServer{
				DBClient: &MockDBClient{
					GetThreadFunc: func(ctx context.Context, req *dbpb.GetThreadRequest, opts ...grpc.CallOption) (*models.Thread, error) {
						return &models.Thread{
							Id:          "123",
							CommunityId: "456",
							Title:       "test thread",
							Content:     "test content",
						}, nil
					},
				},
			}

			_, err := server.GetThread(context.Background(), tt.req)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListThreads_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     *threadpb.ListThreadsRequest
		wantErr error
	}{
		{
			name:    "empty community id",
			req:     &threadpb.ListThreadsRequest{CommunityId: strPtr("")},
			wantErr: status.Error(codes.InvalidArgument, "Community id cannot be empty"),
		},
		{
			name:    "empty title",
			req:     &threadpb.ListThreadsRequest{Title: strPtr("")},
			wantErr: status.Error(codes.InvalidArgument, "Title cannot be empty"),
		},
		{
			name:    "negative_offset",
			req:     &threadpb.ListThreadsRequest{CommunityId: strPtr("123"), Offset: int32Ptr(-1), Limit: int32Ptr(10), SortBy: strPtr("new")},
			wantErr: status.Error(codes.InvalidArgument, "Offset must be a positive integer"),
		},
		{
			name:    "zero_limit",
			req:     &threadpb.ListThreadsRequest{CommunityId: strPtr("123"), Offset: int32Ptr(0), Limit: int32Ptr(0), SortBy: strPtr("new")},
			wantErr: status.Error(codes.InvalidArgument, "Limit must be a positive integer"),
		},
		{
			name:    "empty_sort",
			req:     &threadpb.ListThreadsRequest{CommunityId: strPtr("123"), Offset: int32Ptr(0), Limit: int32Ptr(10), SortBy: strPtr("")},
			wantErr: status.Error(codes.InvalidArgument, "Sort cannot be empty"),
		},
		{
			name: "valid request",
			req: &threadpb.ListThreadsRequest{
				CommunityId: strPtr("123"),
				Title:       strPtr("test"),
				Offset:      int32Ptr(0),
				Limit:       int32Ptr(10),
				SortBy:      strPtr("created_at"),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &src.ThreadServer{
				DBClient: &MockDBClient{
					ListThreadsFunc: func(ctx context.Context, req *dbpb.ListThreadsRequest, opts ...grpc.CallOption) (*dbpb.ListThreadsResponse, error) {
						return &dbpb.ListThreadsResponse{
							Threads: []*models.Thread{},
						}, nil
					},
				},
				CommunityClient: &MockCommunityClient{
					GetCommunityFunc: func(ctx context.Context, req *communitypb.GetCommunityRequest, opts ...grpc.CallOption) (*models.Community, error) {
						return &models.Community{
							Id:   "123",
							Name: "test-community",
						}, nil
					},
				},
			}

			_, err := server.ListThreads(context.Background(), tt.req)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func strPtr(s string) *string { return &s }
func int32Ptr(i int32) *int32 { return &i }

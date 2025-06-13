package test

import (
	"context"
	"testing"

	src "community-service/src"
	communitypb "gen/community-service/pb"
	dbpb "gen/db-service/pb"
	models "gen/models/pb"
	threadpb "gen/thread-service/pb"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/grpc"
)

type MockDBClient struct {
	dbpb.DBServiceClient
	ListCommunitiesFunc  func(ctx context.Context, req *dbpb.ListCommunitiesRequest, opts ...grpc.CallOption) (*dbpb.ListCommunitiesResponse, error)
	CreateCommunityFunc  func(ctx context.Context, req *dbpb.CreateCommunityRequest, opts ...grpc.CallOption) (*dbpb.CreateCommunityResponse, error)
	GetCommunityFunc     func(ctx context.Context, req *dbpb.GetCommunityRequest, opts ...grpc.CallOption) (*models.Community, error)
	UpdateCommunityFunc  func(ctx context.Context, req *dbpb.UpdateCommunityRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteCommunityFunc  func(ctx context.Context, req *dbpb.DeleteCommunityRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

func (m *MockDBClient) ListCommunities(ctx context.Context, req *dbpb.ListCommunitiesRequest, opts ...grpc.CallOption) (*dbpb.ListCommunitiesResponse, error) {
	return m.ListCommunitiesFunc(ctx, req, opts...)
}

func (m *MockDBClient) CreateCommunity(ctx context.Context, req *dbpb.CreateCommunityRequest, opts ...grpc.CallOption) (*dbpb.CreateCommunityResponse, error) {
	return m.CreateCommunityFunc(ctx, req, opts...)
}

func (m *MockDBClient) GetCommunity(ctx context.Context, req *dbpb.GetCommunityRequest, opts ...grpc.CallOption) (*models.Community, error) {
	return m.GetCommunityFunc(ctx, req, opts...)
}

func (m *MockDBClient) UpdateCommunity(ctx context.Context, req *dbpb.UpdateCommunityRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return m.UpdateCommunityFunc(ctx, req, opts...)
}

func (m *MockDBClient) DeleteCommunity(ctx context.Context, req *dbpb.DeleteCommunityRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return m.DeleteCommunityFunc(ctx, req, opts...)
}

type MockThreadClient struct {
	threadpb.ThreadServiceClient
	ListThreadsFunc   func(ctx context.Context, req *threadpb.ListThreadsRequest, opts ...grpc.CallOption) (*threadpb.ListThreadsResponse, error)
	DeleteThreadFunc  func(ctx context.Context, req *threadpb.DeleteThreadRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

func (m *MockThreadClient) ListThreads(ctx context.Context, req *threadpb.ListThreadsRequest, opts ...grpc.CallOption) (*threadpb.ListThreadsResponse, error) {
	return m.ListThreadsFunc(ctx, req, opts...)
}

func (m *MockThreadClient) DeleteThread(ctx context.Context, req *threadpb.DeleteThreadRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return m.DeleteThreadFunc(ctx, req, opts...)
}

func TestCreateCommunity_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     *communitypb.CreateCommunityRequest
		wantErr error
	}{
		{
			name:    "missing name",
			req:     &communitypb.CreateCommunityRequest{},
			wantErr: status.Error(codes.InvalidArgument, "Community name is required"),
		},
		{
			name: "name too short",
			req: &communitypb.CreateCommunityRequest{
				Name: "ab",
			},
			wantErr: status.Error(codes.InvalidArgument, "Name must be between 3 and 50 characters long"),
		},
		{
			name: "name too long",
			req: &communitypb.CreateCommunityRequest{
				Name: "a",
			},
			wantErr: status.Error(codes.InvalidArgument, "Name must be between 3 and 50 characters long"),
		},
		{
			name: "valid request",
			req: &communitypb.CreateCommunityRequest{
				Name: "test-community",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &src.CommunityServer{
				DBClient: &MockDBClient{
					CreateCommunityFunc: func(ctx context.Context, req *dbpb.CreateCommunityRequest, opts ...grpc.CallOption) (*dbpb.CreateCommunityResponse, error) {
						return &dbpb.CreateCommunityResponse{
							Id: "123",
						}, nil
					},
				},
				ThreadClient: &MockThreadClient{},
			}

			_, err := server.CreateCommunity(context.Background(), tt.req)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetCommunity_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     *communitypb.GetCommunityRequest
		wantErr error
	}{
		{
			name:    "missing id",
			req:     &communitypb.GetCommunityRequest{},
			wantErr: status.Error(codes.InvalidArgument, "Community id is required"),
		},
		{
			name: "valid request",
			req: &communitypb.GetCommunityRequest{
				Id: "123",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &src.CommunityServer{
				DBClient: &MockDBClient{
					GetCommunityFunc: func(ctx context.Context, req *dbpb.GetCommunityRequest, opts ...grpc.CallOption) (*models.Community, error) {
						return &models.Community{
							Id:   "123",
							Name: "test-community",
						}, nil
					},
				},
				ThreadClient: &MockThreadClient{},
			}

			_, err := server.GetCommunity(context.Background(), tt.req)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

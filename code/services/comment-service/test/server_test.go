package test

import (
	"context"
	"testing"

	src "comment-service/src"
	commentpb "gen/comment-service/pb"
	dbpb "gen/db-service/pb"
	models "gen/models/pb"
	threadpb "gen/thread-service/pb"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MockDBClient struct {
	dbpb.DBServiceClient
	ListCommentsFunc  func(ctx context.Context, req *dbpb.ListCommentsRequest, opts ...grpc.CallOption) (*dbpb.ListCommentsResponse, error)
	CreateCommentFunc func(ctx context.Context, req *dbpb.CreateCommentRequest, opts ...grpc.CallOption) (*dbpb.CreateCommentResponse, error)
	GetCommentFunc    func(ctx context.Context, req *dbpb.GetCommentRequest, opts ...grpc.CallOption) (*models.Comment, error)
	UpdateCommentFunc func(ctx context.Context, req *dbpb.UpdateCommentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteCommentFunc func(ctx context.Context, req *dbpb.DeleteCommentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

func (m *MockDBClient) ListComments(ctx context.Context, req *dbpb.ListCommentsRequest, opts ...grpc.CallOption) (*dbpb.ListCommentsResponse, error) {
	return m.ListCommentsFunc(ctx, req, opts...)
}

func (m *MockDBClient) CreateComment(ctx context.Context, req *dbpb.CreateCommentRequest, opts ...grpc.CallOption) (*dbpb.CreateCommentResponse, error) {
	return m.CreateCommentFunc(ctx, req, opts...)
}

func (m *MockDBClient) GetComment(ctx context.Context, req *dbpb.GetCommentRequest, opts ...grpc.CallOption) (*models.Comment, error) {
	return m.GetCommentFunc(ctx, req, opts...)
}

func (m *MockDBClient) UpdateComment(ctx context.Context, req *dbpb.UpdateCommentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return m.UpdateCommentFunc(ctx, req, opts...)
}

func (m *MockDBClient) DeleteComment(ctx context.Context, req *dbpb.DeleteCommentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return m.DeleteCommentFunc(ctx, req, opts...)
}

type MockThreadClient struct {
	threadpb.ThreadServiceClient
	UpdateThreadFunc func(ctx context.Context, req *threadpb.UpdateThreadRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

func (m *MockThreadClient) UpdateThread(ctx context.Context, req *threadpb.UpdateThreadRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return m.UpdateThreadFunc(ctx, req, opts...)
}

func TestCreateComment_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     *commentpb.CreateCommentRequest
		wantErr error
	}{
		{
			name:    "missing parent id",
			req:     &commentpb.CreateCommentRequest{},
			wantErr: status.Error(codes.InvalidArgument, "Parent id is required"),
		},
		{
			name: "missing content",
			req: &commentpb.CreateCommentRequest{
				ParentId: "123",
			},
			wantErr: status.Error(codes.InvalidArgument, "Content is required"),
		},
		{
			name: "content too long",
			req: &commentpb.CreateCommentRequest{
				ParentId: "123",
				Content:  string(make([]byte, 501)),
			},
			wantErr: status.Error(codes.InvalidArgument, "Content exceeds maximum length of 500 characters"),
		},
		{
			name: "valid request",
			req: &commentpb.CreateCommentRequest{
				ParentId:   "123",
				Content:    "test comment",
				ParentType: models.CommentParentType_THREAD,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &src.CommentServer{
				DBClient: &MockDBClient{
					CreateCommentFunc: func(ctx context.Context, req *dbpb.CreateCommentRequest, opts ...grpc.CallOption) (*dbpb.CreateCommentResponse, error) {
						return &dbpb.CreateCommentResponse{
							Id: "123",
						}, nil
					},
				},
				ThreadClient: &MockThreadClient{
					UpdateThreadFunc: func(ctx context.Context, req *threadpb.UpdateThreadRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
						return &emptypb.Empty{}, nil
					},
				},
			}

			_, err := server.CreateComment(context.Background(), tt.req)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetComment_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     *commentpb.GetCommentRequest
		wantErr error
	}{
		{
			name:    "missing id",
			req:     &commentpb.GetCommentRequest{},
			wantErr: status.Error(codes.InvalidArgument, "Comment id is required"),
		},
		{
			name: "valid request",
			req: &commentpb.GetCommentRequest{
				Id: "123",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &src.CommentServer{
				DBClient: &MockDBClient{
					GetCommentFunc: func(ctx context.Context, req *dbpb.GetCommentRequest, opts ...grpc.CallOption) (*models.Comment, error) {
						return &models.Comment{
							Id:      "123",
							Content: "test comment",
						}, nil
					},
				},
				ThreadClient: &MockThreadClient{},
			}

			_, err := server.GetComment(context.Background(), tt.req)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

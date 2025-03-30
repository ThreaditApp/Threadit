package server

import (
	"context"
	commentpb "gen/comment-service/pb"
	dbpb "gen/db-service/pb"
	models "gen/models/pb"
	threadpb "gen/thread-service/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CommentServer struct {
	commentpb.UnimplementedCommentServiceServer
	DBClient     dbpb.DBServiceClient
	ThreadClient threadpb.ThreadServiceClient
}

func (s *CommentServer) ListComments(ctx context.Context, req *commentpb.ListCommentsRequest) (*commentpb.ListCommentsResponse, error) {
	res, err := s.DBClient.ListComments(ctx, &dbpb.ListCommentsRequest{
		ThreadId: req.ThreadId,
		Offset:   req.Offset,
		Limit:    req.Limit,
		SortBy:   req.SortBy,
	})
	if err != nil {
		return nil, err
	}
	return &commentpb.ListCommentsResponse{
		Comments: res.Comments,
	}, nil
}

func (s *CommentServer) CreateComment(ctx context.Context, req *commentpb.CreateCommentRequest) (*commentpb.CreateCommentResponse, error) {
	var err error
	if req.ParentType == models.CommentParentType_THREAD {
		// check if thread exists
		_, err = s.ThreadClient.GetThread(ctx, &threadpb.GetThreadRequest{
			Id: req.ParentId,
		})

	} else {
		// check if comment exists
		_, err = s.GetComment(ctx, &commentpb.GetCommentRequest{
			Id: req.ParentId,
		})
	}
	if err != nil {
		return nil, err
	}

	// create comment
	res, err := s.DBClient.CreateComment(ctx, &dbpb.CreateCommentRequest{
		Content:    req.Content,
		ParentId:   req.ParentId,
		ParentType: req.ParentType,
	})
	if err != nil {
		return nil, err
	}
	return &commentpb.CreateCommentResponse{
		Id: res.Id,
	}, nil
}

func (s *CommentServer) GetComment(ctx context.Context, req *commentpb.GetCommentRequest) (*models.Comment, error) {
	res, err := s.DBClient.GetComment(ctx, &dbpb.GetCommentRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *CommentServer) UpdateComment(ctx context.Context, req *commentpb.UpdateCommentRequest) (*emptypb.Empty, error) {
	_, err := s.DBClient.UpdateComment(ctx, &dbpb.UpdateCommentRequest{
		Id:      req.Id,
		Content: req.Content,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *CommentServer) DeleteComment(ctx context.Context, req *commentpb.DeleteCommentRequest) (*emptypb.Empty, error) {
	_, err := s.DBClient.DeleteComment(ctx, &dbpb.DeleteCommentRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

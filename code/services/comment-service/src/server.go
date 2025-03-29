package server

import (
	"context"
	"fmt"
	commentpb "gen/comment-service/pb"
	dbpb "gen/db-service/pb"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

type CommentServer struct {
	commentpb.UnimplementedCommentServiceServer
	DBClient dbpb.DBServiceClient
}

func (s *CommentServer) ListComments(ctx context.Context, req *commentpb.ListCommentsRequest) (*commentpb.ListCommentsResponse, error) {
	log.Printf("ListComments called with post_id: %s", req.PostId)

	res, err := s.DBClient.ListComments(ctx, &dbpb.ListCommentsRequest{
		PostId:   req.PostId,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	comments := make([]*commentpb.Comment, len(res.Comments))
	for i, comment := range res.Comments {
		comments[i] = &commentpb.Comment{
			Id:        comment.Id,
			PostId:    comment.PostId,
			UserId:    comment.UserId,
			Content:   comment.Content,
			ParentId:  comment.ParentId,
			CreatedAt: comment.CreatedAt,
		}
	}

	return &commentpb.ListCommentsResponse{
		Comments: comments,
		Pagination: &commentpb.Pagination{
			CurrentPage: res.Pagination.CurrentPage,
			PerPage:     res.Pagination.PerPage,
			TotalItems:  res.Pagination.TotalItems,
			TotalPages:  res.Pagination.TotalPages,
		},
	}, nil
}

func (s *CommentServer) CreateComment(ctx context.Context, req *commentpb.CreateCommentRequest) (*commentpb.Comment, error) {

	res, err := s.DBClient.CreateComment(ctx, &dbpb.CreateCommentRequest{
		ThreadId: req.ThreadId,
		Content:  req.Content,
		ParentId: req.ParentId, // TODO: fix
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	return &commentpb.Comment{
		Id:       res.Comment.Id,
		Content:  res.Comment.Content,
		ParentId: res.Comment.ParentId,
	}, nil
}

func (s *CommentServer) GetComment(ctx context.Context, req *commentpb.GetCommentRequest) (*commentpb.Comment, error) {
	log.Printf("GetComment called with id: %s", req.Id)

	res, err := s.DBClient.GetComment(ctx, &dbpb.GetCommentRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	return &commentpb.Comment{
		Id:       res.Id,
		Content:  res.Content,
		ParentId: res.ParentId,
	}, nil
}

func (s *CommentServer) UpdateComment(ctx context.Context, req *commentpb.UpdateCommentRequest) (*emptypb.Empty, error) {

	_, err := s.DBClient.UpdateComment(ctx, &dbpb.UpdateCommentRequest{
		Id:      req.Id,
		Content: req.Content,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	return nil, nil
}

func (s *CommentServer) DeleteComment(ctx context.Context, req *commentpb.DeleteCommentRequest) (*emptypb.Empty, error) {

	_, err := s.DBClient.DeleteComment(ctx, &dbpb.DeleteCommentRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	return &emptypb.Empty{}, nil
}

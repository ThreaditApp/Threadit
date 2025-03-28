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
	userId, err := getCurrentUserId(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.DBClient.CreateComment(ctx, &dbpb.CreateCommentRequest{
		PostId:   req.PostId,
		UserId:   userId,
		Content:  req.Content,
		ParentId: req.ParentId,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	return &commentpb.Comment{
		Id:        res.Comment.Id,
		PostId:    res.Comment.PostId,
		UserId:    res.Comment.UserId,
		Content:   res.Comment.Content,
		ParentId:  res.Comment.ParentId,
		CreatedAt: res.Comment.CreatedAt,
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
		Id:        res.Id,
		PostId:    res.PostId,
		UserId:    res.UserId,
		Content:   res.Content,
		ParentId:  res.ParentId,
		CreatedAt: res.CreatedAt,
	}, nil
}

func (s *CommentServer) UpdateComment(ctx context.Context, req *commentpb.UpdateCommentRequest) (*emptypb.Empty, error) {
	userId, err := getCurrentUserId(ctx)
	if err != nil {
		return nil, err
	}

	_, err = s.DBClient.UpdateComment(ctx, &dbpb.UpdateCommentRequest{
		Id:      req.Id,
		UserId:  userId,
		Content: req.Content,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	return nil, nil
}

func (s *CommentServer) DeleteComment(ctx context.Context, req *commentpb.DeleteCommentRequest) (*emptypb.Empty, error) {
	userId, err := getCurrentUserId(ctx)
	if err != nil {
		return nil, err
	}

	_, err = s.DBClient.DeleteComment(ctx, &dbpb.DeleteCommentRequest{
		Id:     req.Id,
		UserId: userId,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func getCurrentUserId(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("no metadata found in context")
	}
	userIds := md.Get("x-user-id")
	if len(userIds) == 0 {
		return "", fmt.Errorf("user id not found in metadata")
	}
	return userIds[0], nil
}

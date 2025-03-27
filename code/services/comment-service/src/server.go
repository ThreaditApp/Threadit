package server

import (
	"context"
	"fmt"
	"log"

	pb "comment-service/src/pb"
	dbpb "db-service/src/pb"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CommentServer struct {
	pb.UnimplementedCommentServiceServer
	DBClient pb.DBServiceClient
}

func (s *CommentServer) ListComments(ctx context.Context, req *pb.ListCommentsRequest) (*pb.ListCommentsResponse, error) {
	log.Printf("ListComments called with post_id: %s", req.PostId)

	res, err := s.DBClient.ListComments(ctx, &dbpb.ListCommentsRequest{
		PostId:   req.PostId,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	comments := make([]*pb.Comment, len(res.Comments))
	for i, comment := range res.Comments {
		comments[i] = &pb.Comment{
			Id:        comment.Id,
			PostId:    comment.PostId,
			UserId:    comment.UserId,
			Content:   comment.Content,
			ParentId:  comment.ParentId,
			CreatedAt: comment.CreatedAt,
		}
	}

	return &pb.ListCommentsResponse{
		Comments: comments,
		Pagination: &pb.Pagination{
			CurrentPage: res.Pagination.CurrentPage,
			PerPage:     res.Pagination.PerPage,
			TotalItems:  res.Pagination.TotalItems,
			TotalPages:  res.Pagination.TotalPages,
		},
	}, nil
}

func (s *CommentServer) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.Comment, error) {
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

	return &pb.Comment{
		Id:        res.Comment.Id,
		PostId:    res.Comment.PostId,
		UserId:    res.Comment.UserId,
		Content:   res.Comment.Content,
		ParentId:  res.Comment.ParentId,
		CreatedAt: res.Comment.CreatedAt,
	}, nil
}

func (s *CommentServer) GetComment(ctx context.Context, req *pb.GetCommentRequest) (*pb.Comment, error) {
	log.Printf("GetComment called with id: %s", req.Id)

	res, err := s.DBClient.GetComment(ctx, &dbpb.GetCommentRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	return &pb.Comment{
		Id:        res.Comment.Id,
		PostId:    res.Comment.PostId,
		UserId:    res.Comment.UserId,
		Content:   res.Comment.Content,
		ParentId:  res.Comment.ParentId,
		CreatedAt: res.Comment.CreatedAt,
	}, nil
}

func (s *CommentServer) UpdateComment(ctx context.Context, req *pb.UpdateCommentRequest) (*pb.Comment, error) {
	userId, err := getCurrentUserId(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.DBClient.UpdateComment(ctx, &dbpb.UpdateCommentRequest{
		Id:      req.Id,
		UserId:  userId,
		Content: req.Content,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	return &pb.Comment{
		Id:        res.Comment.Id,
		PostId:    res.Comment.PostId,
		UserId:    res.Comment.UserId,
		Content:   res.Comment.Content,
		ParentId:  res.Comment.ParentId,
		CreatedAt: res.Comment.CreatedAt,
	}, nil
}

func (s *CommentServer) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*emptypb.Empty, error) {
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

package server

import (
	"context"
	"fmt"
	"log"

	dbpb "db-service/src/pb"
	pb "feed-service/src/pb"
)

type FeedServer struct {
	pb.UnimplementedFeedServiceServer
	DBClient dbpb.DBServiceClient
}

func (s *FeedServer) GetUserFeed(ctx context.Context, req *pb.GetUserFeedRequest) (*pb.GetUserFeedResponse, error) {
	log.Printf("GetUserFeed called with page: %d, page_size: %d, sort: %s", req.Page, req.PageSize, req.Sort)

	res, err := s.DBClient.GetUserFeed(ctx, &dbpb.GetUserFeedRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Sort:     req.Sort,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	posts := make([]*pb.Feed, len(res.Posts))
	for i, post := range res.Posts {
		posts[i] = &pb.Feed{
			Id:          post.Id,
			Type:        post.Type,
			CommunityId: post.CommunityId,
			Title:       post.Title,
			Content:     post.Content,
			CreatedAt:   post.CreatedAt,
			UpdatedAt:   post.UpdatedAt,
		}
	}

	return &pb.GetUserFeedResponse{
		Posts: posts,
		Pagination: &pb.Pagination{
			CurrentPage: res.Pagination.CurrentPage,
			PerPage:     res.Pagination.PerPage,
			TotalItems:  res.Pagination.TotalItems,
			TotalPages:  res.Pagination.TotalPages,
		},
	}, nil
}

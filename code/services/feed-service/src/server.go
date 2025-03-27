package server

import (
	"context"
	"fmt"
	"log"

	pb "feed-service/src/pb"
	socialpb "social-service/src/pb"
	threadpb "thread-service/src/pb"
)

type FeedServer struct {
	pb.UnimplementedFeedServiceServer
	SocialClient socialpb.SocialServiceClient
	ThreadClient threadpb.ThreadServiceClient
}

func (s *FeedServer) GetUserFeed(ctx context.Context, req *pb.GetUserFeedRequest) (*pb.GetUserFeedResponse, error) {
	log.Printf("GetUserFeed called with page: %d, page_size: %d, sort: %s", req.Page, req.PageSize, req.Sort)

	// Get the list of communities and users the user is following
	followingRes, err := s.SocialClient.GetFollowing(ctx, &socialpb.GetFollowingRequest{})
	if err != nil {
		return nil, fmt.Errorf("error fetching following data: %w", err)
	}

	// Fetch threads from the followed communities and users
	threadRes, err := s.ThreadClient.ListThreads(ctx, &threadpb.ListThreadsRequest{
		CommunityId: followingRes.CommunityIds,
		AuthorId:    followingRes.UserIds,
		Page:        req.Page,
		PageSize:    req.PageSize,
		SortBy:      "created_at",
		SortOrder:   "desc",
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching threads: %w", err)
	}

	// Map threads to the response
	posts := make([]*pb.Thread, len(threadRes.Threads))
	for i, thread := range threadRes.Threads {
		posts[i] = &pb.Thread{
			Id:          thread.Id,
			Type:        "thread",
			CommunityId: thread.CommunityId,
			Title:       thread.Title,
			Content:     thread.Content,
			CreatedAt:   thread.CreatedAt,
			UpdatedAt:   thread.UpdatedAt,
		}
	}

	return &pb.GetUserFeedResponse{
		Posts: posts,
		Pagination: &pb.Pagination{
			CurrentPage: threadRes.Pagination.CurrentPage,
			PerPage:     threadRes.Pagination.PerPage,
			TotalItems:  threadRes.Pagination.TotalItems,
			TotalPages:  threadRes.Pagination.TotalPages,
		},
	}, nil

}

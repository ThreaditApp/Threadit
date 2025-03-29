package server

import (
	"context"
	"fmt"
	communitypb "gen/community-service/pb"
	dbpb "gen/db-service/pb"
	threadpb "gen/thread-service/pb"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ThreadServer struct {
	threadpb.UnimplementedThreadServiceServer
	CommunityClient communitypb.CommunityServiceClient
	DBClient        dbpb.DBServiceClient
}

func (s *ThreadServer) ListThreads(ctx context.Context, req *threadpb.ListThreadsRequest) (*threadpb.ListThreadsResponse, error) {
	res, err := s.DBClient.ListThreads(ctx, &dbpb.ListThreadsRequest{
		Page:        req.Page,
		PageSize:    req.PageSize,
		CommunityId: req.CommunityId,
		AuthorId:    req.AuthorId,
		Title:       req.Title,
		SortBy:      req.SortBy,
		SortOrder:   req.SortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	threads := make([]*threadpb.Thread, len(res.Threads))
	for i, thread := range res.Threads {
		threads[i] = &threadpb.Thread{
			Id:          thread.Id,
			CommunityId: thread.CommunityId,
			AuthorId:    thread.AuthorId,
			Title:       thread.Title,
			Content:     thread.Content,
			CreatedAt:   thread.CreatedAt,
			UpdatedAt:   thread.UpdatedAt,
		}
	}

	return &threadpb.ListThreadsResponse{
		Threads: threads,
		Pagination: &threadpb.Pagination{
			CurrentPage: res.Pagination.CurrentPage,
			PerPage:     res.Pagination.PerPage,
			TotalItems:  res.Pagination.TotalItems,
			TotalPages:  res.Pagination.TotalPages,
		},
	}, nil
}

func (s *ThreadServer) CreateThread(ctx context.Context, req *threadpb.CreateThreadRequest) (*threadpb.Thread, error) {
	userId, err := getCurrentUserId(ctx)
	if err != nil {
		return nil, err
	}

	_, err = s.CommunityClient.GetCommunity(ctx, &communitypb.GetCommunityRequest{
		Id: req.CommunityId,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling community service: %w", err)
	}

	res, err := s.DBClient.CreateThread(ctx, &dbpb.CreateThreadRequest{
		CommunityId: req.CommunityId,
		AuthorId:    userId,
		Title:       req.Title,
		Content:     req.Content,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	return &threadpb.Thread{
		Id:          res.Id,
		CommunityId: res.CommunityId,
		AuthorId:    res.AuthorId,
		Title:       res.Title,
		Content:     res.Content,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}, nil
}

func (s *ThreadServer) GetThread(ctx context.Context, req *threadpb.GetThreadRequest) (*threadpb.Thread, error) {
	res, err := s.DBClient.GetThread(ctx, &dbpb.GetThreadRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	return &threadpb.Thread{
		Id:          res.Id,
		CommunityId: res.CommunityId,
		AuthorId:    res.AuthorId,
		Title:       res.Title,
		Content:     res.Content,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}, nil
}

func (s *ThreadServer) UpdateThread(ctx context.Context, req *threadpb.UpdateThreadRequest) (*threadpb.Thread, error) {
	userId, err := getCurrentUserId(ctx)
	if err != nil {
		return nil, err
	}

	thread, err := s.DBClient.GetThread(ctx, &dbpb.GetThreadRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}
	if thread.AuthorId != userId {
		return nil, fmt.Errorf("user is not the author of the thread")
	}

	res, err := s.DBClient.UpdateThread(ctx, &dbpb.UpdateThreadRequest{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}

	return &threadpb.Thread{
		Id:          res.Id,
		CommunityId: res.CommunityId,
		AuthorId:    res.AuthorId,
		Title:       res.Title,
		Content:     res.Content,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}, nil
}

func (s *ThreadServer) DeleteThread(ctx context.Context, req *threadpb.DeleteThreadRequest) (*emptypb.Empty, error) {
	userId, err := getCurrentUserId(ctx)
	if err != nil {
		return nil, err
	}

	thread, err := s.DBClient.GetThread(ctx, &dbpb.GetThreadRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling database service: %w", err)
	}
	if thread.AuthorId != userId {
		return nil, fmt.Errorf("user is not the author of the thread")
	}

	_, err = s.DBClient.DeleteThread(ctx, &dbpb.DeleteThreadRequest{
		Id: req.Id,
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

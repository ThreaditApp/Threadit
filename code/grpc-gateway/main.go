package main

import (
	"context"
	"encoding/json"
	"fmt"
	commentpb "gen/comment-service/pb"
	communitypb "gen/community-service/pb"
	popularpb "gen/popular-service/pb"
	searchpb "gen/search-service/pb"
	threadpb "gen/thread-service/pb"
	votepb "gen/vote-service/pb"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func getGrpcServerAddress(hostEnvVar string, portEnvVar string) string {
	host := os.Getenv(hostEnvVar)
	if host == "" {
		log.Fatalf("missing %s env var", hostEnvVar)
	}
	port := os.Getenv(portEnvVar)
	if port == "" {
		log.Fatalf("missing %s env var", portEnvVar)
	}
	return fmt.Sprintf("%s:%s", host, port)
}

func connectGrpcClient(hostEnvVar string, portEnvVar string) *grpc.ClientConn {
	addr := getGrpcServerAddress(hostEnvVar, portEnvVar)
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024*1024*500), // 500MB
			grpc.MaxCallSendMsgSize(1024*1024*500), // 500MB
		),
	)
	if err != nil {
		log.Fatalf("failed to connect to %s: %v", addr, err)
	}
	return conn
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	health := map[string]bool{
		"community-service": false,
		"thread-service":    false,
		"comment-service":   false,
		"vote-service":      false,
		"search-service":    false,
		"popular-service":   false,
	}

	communityConn := connectGrpcClient("COMMUNITY_SERVICE_HOST", "COMMUNITY_SERVICE_PORT")
	defer communityConn.Close()
	communityClient := communitypb.NewCommunityServiceClient(communityConn)
	_, err := communityClient.CheckHealth(ctx, &emptypb.Empty{})
	health["community-service"] = err == nil

	threadConn := connectGrpcClient("THREAD_SERVICE_HOST", "THREAD_SERVICE_PORT")
	defer threadConn.Close()
	threadClient := threadpb.NewThreadServiceClient(threadConn)
	_, err = threadClient.CheckHealth(ctx, &emptypb.Empty{})
	health["thread-service"] = err == nil

	commentConn := connectGrpcClient("COMMENT_SERVICE_HOST", "COMMENT_SERVICE_PORT")
	defer commentConn.Close()
	commentClient := commentpb.NewCommentServiceClient(commentConn)
	_, err = commentClient.CheckHealth(ctx, &emptypb.Empty{})
	health["comment-service"] = err == nil

	voteConn := connectGrpcClient("VOTE_SERVICE_HOST", "VOTE_SERVICE_PORT")
	defer voteConn.Close()
	voteClient := votepb.NewVoteServiceClient(voteConn)
	_, err = voteClient.CheckHealth(ctx, &emptypb.Empty{})
	health["vote-service"] = err == nil

	searchConn := connectGrpcClient("SEARCH_SERVICE_HOST", "SEARCH_SERVICE_PORT")
	defer searchConn.Close()
	searchClient := searchpb.NewSearchServiceClient(searchConn)
	_, err = searchClient.CheckHealth(ctx, &emptypb.Empty{})
	health["search-service"] = err == nil

	popularConn := connectGrpcClient("POPULAR_SERVICE_HOST", "POPULAR_SERVICE_PORT")
	defer popularConn.Close()
	popularClient := popularpb.NewPopularServiceClient(popularConn)
	_, err = popularClient.CheckHealth(ctx, &emptypb.Empty{})
	health["popular-service"] = err == nil

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonResponse, err := json.Marshal(health)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func main() {
	gwmux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024*1024*500), // 500MB
			grpc.MaxCallSendMsgSize(1024*1024*500), // 500MB
		),
	}

	err := communitypb.RegisterCommunityServiceHandlerFromEndpoint(context.Background(), gwmux, getGrpcServerAddress("COMMUNITY_SERVICE_HOST", "COMMUNITY_SERVICE_PORT"), opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC gateway: %v", err)
	}

	err = threadpb.RegisterThreadServiceHandlerFromEndpoint(context.Background(), gwmux, getGrpcServerAddress("THREAD_SERVICE_HOST", "THREAD_SERVICE_PORT"), opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC gateway: %v", err)
	}

	err = commentpb.RegisterCommentServiceHandlerFromEndpoint(context.Background(), gwmux, getGrpcServerAddress("COMMENT_SERVICE_HOST", "COMMENT_SERVICE_PORT"), opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC gateway: %v", err)
	}

	err = votepb.RegisterVoteServiceHandlerFromEndpoint(context.Background(), gwmux, getGrpcServerAddress("VOTE_SERVICE_HOST", "VOTE_SERVICE_PORT"), opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC gateway: %v", err)
	}

	err = searchpb.RegisterSearchServiceHandlerFromEndpoint(context.Background(), gwmux, getGrpcServerAddress("SEARCH_SERVICE_HOST", "SEARCH_SERVICE_PORT"), opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC gateway: %v", err)
	}

	err = popularpb.RegisterPopularServiceHandlerFromEndpoint(context.Background(), gwmux, getGrpcServerAddress("POPULAR_SERVICE_HOST", "POPULAR_SERVICE_PORT"), opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC gateway: %v", err)
	}

	http.HandleFunc("/health", handleHealthCheck)

	http.Handle("/", gwmux)

	port := os.Getenv("GRPC_GATEWAY_PORT")
	if port == "" {
		log.Fatalf("missing GRPC_GATEWAY_PORT env var")
	}

	log.Printf("gRPC Gateway server listening on :%s", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

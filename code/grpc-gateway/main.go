package main

import (
	"context"
	"fmt"
	commentpb "gen/comment-service/pb"
	communitypb "gen/community-service/pb"
	popularpb "gen/popular-service/pb"
	searchpb "gen/search-service/pb"
	threadpb "gen/thread-service/pb"
	votepb "gen/vote-service/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
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

func matchHeader(key string) (string, bool) {
	if key == "X-User-Id" {
		return "x-user-id", true
	}
	return key, false
}

func main() {
	gwmux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(matchHeader),
	)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

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

	err = popularpb.RegisterPopularServiceHandlerFromEndpoint(context.Background(), gwmux, getGrpcServerAddress("FEED_SERVICE_HOST", "FEED_SERVICE_PORT"), opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC gateway: %v", err)
	}

	http.Handle("/", gwmux)

	port := os.Getenv("GRPC_GATEWAY_PORT")
	if port == "" {
		log.Fatalf("missing GRPC_GATEWAY_PORT env var")
	}

	fmt.Printf("gRPC Gateway server listening on :%s", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

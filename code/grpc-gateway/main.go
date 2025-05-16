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
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"threadit/grpc-gateway/middleware"
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

func main() {
	ctx := context.Background()
	mux := runtime.NewServeMux()

	// Initialize auth handler
	authHandler := middleware.NewAuthHandler(
		os.Getenv("KEYCLOAK_URL"),
		os.Getenv("KEYCLOAK_CLIENT_ID"),
		os.Getenv("KEYCLOAK_CLIENT_SECRET"),
		os.Getenv("KEYCLOAK_REALM"),
	)

	// Create a new ServeMux for both gRPC-Gateway and auth routes
	httpMux := http.NewServeMux()

	// Register auth routes
	authHandler.RegisterRoutes(httpMux)

	// gRPC dial options with message size configurations
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024*1024*500), // 500MB
			grpc.MaxCallSendMsgSize(1024*1024*500), // 500MB
		),
	}

	// Register gRPC-Gateway routes with auth middleware
	httpMux.Handle("/api/v1/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Auth middleware for API routes
		authMiddleware := middleware.NewAuthMiddleware(middleware.KeycloakConfig{
			Realm:        os.Getenv("KEYCLOAK_REALM"),
			ClientID:     os.Getenv("KEYCLOAK_CLIENT_ID"),
			ClientSecret: os.Getenv("KEYCLOAK_CLIENT_SECRET"),
			KeycloakURL:  os.Getenv("KEYCLOAK_URL"),
		})
		
		authMiddleware.Handler(mux).ServeHTTP(w, r)
	}))

	// Register service handlers
	if err := registerServices(ctx, mux, opts); err != nil {
		log.Fatalf("Failed to register services: %v", err)
	}

	port := os.Getenv("GRPC_GATEWAY_PORT")
	if port == "" {
		log.Fatalf("missing GRPC_GATEWAY_PORT env var")
	}

	log.Printf("gRPC Gateway server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, httpMux); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func registerServices(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) error {
	// Register Community Service
	if err := communitypb.RegisterCommunityServiceHandlerFromEndpoint(
		ctx, mux, getGrpcServerAddress("COMMUNITY_SERVICE_HOST", "COMMUNITY_SERVICE_PORT"), opts,
	); err != nil {
		return fmt.Errorf("failed to register community service: %v", err)
	}

	// Register Thread Service
	if err := threadpb.RegisterThreadServiceHandlerFromEndpoint(
		ctx, mux, getGrpcServerAddress("THREAD_SERVICE_HOST", "THREAD_SERVICE_PORT"), opts,
	); err != nil {
		return fmt.Errorf("failed to register thread service: %v", err)
	}

	// Register Comment Service
	if err := commentpb.RegisterCommentServiceHandlerFromEndpoint(
		ctx, mux, getGrpcServerAddress("COMMENT_SERVICE_HOST", "COMMENT_SERVICE_PORT"), opts,
	); err != nil {
		return fmt.Errorf("failed to register comment service: %v", err)
	}

	// Register Vote Service
	if err := votepb.RegisterVoteServiceHandlerFromEndpoint(
		ctx, mux, getGrpcServerAddress("VOTE_SERVICE_HOST", "VOTE_SERVICE_PORT"), opts,
	); err != nil {
		return fmt.Errorf("failed to register vote service: %v", err)
	}

	// Register Search Service
	if err := searchpb.RegisterSearchServiceHandlerFromEndpoint(
		ctx, mux, getGrpcServerAddress("SEARCH_SERVICE_HOST", "SEARCH_SERVICE_PORT"), opts,
	); err != nil {
		return fmt.Errorf("failed to register search service: %v", err)
	}

	// Register Popular Service
	if err := popularpb.RegisterPopularServiceHandlerFromEndpoint(
		ctx, mux, getGrpcServerAddress("POPULAR_SERVICE_HOST", "POPULAR_SERVICE_PORT"), opts,
	); err != nil {
		return fmt.Errorf("failed to register popular service: %v", err)
	}

	return nil
}

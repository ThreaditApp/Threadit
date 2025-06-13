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
	gorun "runtime"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
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
	// Set maximum number of CPUs to use
	gorun.GOMAXPROCS(gorun.NumCPU())

	gwmux := runtime.NewServeMux()
	ctx := context.Background()

	// gRPC dial options with message size configurations
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024*1024*500), // 500MB
			grpc.MaxCallSendMsgSize(1024*1024*500), // 500MB
		),
	}

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

	// Register gRPC-Gateway routes with auth middleware
	httpMux.Handle("/api", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Auth middleware for API routes
		authMiddleware := middleware.NewAuthMiddleware(middleware.KeycloakConfig{
			Realm:        os.Getenv("KEYCLOAK_REALM"),
			ClientID:     os.Getenv("KEYCLOAK_CLIENT_ID"),
			ClientSecret: os.Getenv("KEYCLOAK_CLIENT_SECRET"),
			KeycloakURL:  os.Getenv("KEYCLOAK_URL"),
		})

		authMiddleware.Handler(gwmux).ServeHTTP(w, r)
	}))

	// Register service handlers
	if err := registerServices(ctx, gwmux, opts); err != nil {
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

	http.HandleFunc("/health", handleHealthCheck)
	http.Handle("/", mux)

	return nil
}

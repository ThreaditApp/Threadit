package main

import (
	server "feed-service/src"
	"fmt"
	feedpb "gen/feed-service/pb"
	socialpb "gen/social-service/pb"
	threadpb "gen/thread-service/pb"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func connectGrpcClient(serviceName string, portEnvVar string) *grpc.ClientConn {
	port := os.Getenv(portEnvVar)
	if port == "" {
		log.Fatalf("missing %s env var", portEnvVar)
	}
	addr := fmt.Sprintf("localhost:%s", port)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to %s: %v", serviceName, err)
	}
	return conn
}

func main() {
	// Connect to other services
	threadConn := connectGrpcClient("thread service", "THREAD_SERVICE_PORT")
	defer threadConn.Close()

	socialConn := connectGrpcClient("social service", "SOCIAL_SERVICE_PORT")
	defer socialConn.Close()

	// Create feed service with database client
	feedService := &server.FeedServer{
		ThreadClient: threadpb.NewThreadServiceClient(threadConn),
		SocialClient: socialpb.NewSocialServiceClient(socialConn),
	}

	// get env port
	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		log.Fatalf("missing SERVICE_PORT env var")
	}

	// Start gRPC server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	feedpb.RegisterFeedServiceServer(grpcServer, feedService)

	log.Printf("gRPC server is listening on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

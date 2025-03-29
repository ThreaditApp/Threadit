package main

import (
	"fmt"
	communitypb "gen/community-service/pb"
	searchpb "gen/search-service/pb"
	threadpb "gen/thread-service/pb"
	userpb "gen/user-service/pb"
	"log"
	"net"
	"os"
	server "search-service/src"

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
	// connect to other services
	userConn := connectGrpcClient("user service", "USER_SERVICE_PORT")
	defer userConn.Close()
	communityConn := connectGrpcClient("community service", "COMMUNITY_SERVICE_PORT")
	defer communityConn.Close()
	threadConn := connectGrpcClient("thread service", "THREAD_SERVICE_PORT")
	defer threadConn.Close()

	// create search server with clients
	searchServer := &server.SearchServer{
		UserClient:      userpb.NewUserServiceClient(userConn),
		CommunityClient: communitypb.NewCommunityServiceClient(communityConn),
		ThreadClient:    threadpb.NewThreadServiceClient(threadConn),
	}

	// get env port
	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		log.Fatalf("missing SERVICE_PORT env var")
	}

	// start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	searchpb.RegisterSearchServiceServer(grpcServer, searchServer)

	log.Printf("gRPC server is listening on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

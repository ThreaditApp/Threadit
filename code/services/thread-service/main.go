package main

import (
	"context"
	"fmt"
	communitypb "gen/community-service/pb"
	dbpb "gen/db-service/pb"
	threadpb "gen/thread-service/pb"
	"log"
	"net"
	"os"
	server "thread-service/src"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func connectGrpcClient(hostEnvVar string, portEnvVar string) *grpc.ClientConn {
	host := os.Getenv(hostEnvVar)
	if host == "" {
		log.Fatalf("missing %s env var", hostEnvVar)
	}
	port := os.Getenv(portEnvVar)
	if port == "" {
		log.Fatalf("missing %s env var", portEnvVar)
	}
	addr := fmt.Sprintf("%s:%s", host, port)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to %s: %v", addr, err)
	}
	return conn
}

func main() {
	// connect to database service
	dbConn := connectGrpcClient("DB_SERVICE_HOST", "DB_SERVICE_PORT")
	defer dbConn.Close()

	// connect to community service
	communityConn := connectGrpcClient("COMMUNITY_SERVICE_HOST", "COMMUNITY_SERVICE_PORT")
	defer communityConn.Close()

	// create thread service with database service
	threadService := &server.ThreadServer{
		DBClient:        dbpb.NewDBServiceClient(dbConn),
		CommunityClient: communitypb.NewCommunityServiceClient(communityConn),
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
	threadpb.RegisterThreadServiceServer(grpcServer, threadService)

	// hardcode a thread creation
	threadService.CreateThread(context.Background(), &threadpb.CreateThreadRequest{
		CommunityId: "1",
		Title:       "Hello World",
		Content:     "This is a test thread",
	})

	log.Printf("gRPC server is listening on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

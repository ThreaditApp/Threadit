package main

import (
	"fmt"
	commentpb "gen/comment-service/pb"
	threadpb "gen/thread-service/pb"
	votepb "gen/vote-service/pb"
	"log"
	"net"
	"os"
	server "vote-service/src"

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
	threadConn := connectGrpcClient("thread service", "THREAD_SERVICE_PORT")
	defer threadConn.Close()
	commentConn := connectGrpcClient("comment service", "COMMENT_SERVICE_PORT")
	defer commentConn.Close()

	// create social service with database service
	voteService := &server.VoteServer{
		ThreadClient:  threadpb.NewThreadServiceClient(threadConn),
		CommentClient: commentpb.NewCommentServiceClient(commentConn),
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
	votepb.RegisterVoteServiceServer(grpcServer, voteService)

	log.Printf("gRPC server is listening on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

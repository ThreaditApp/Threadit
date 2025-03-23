package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	server "vote-service/src"
	pb "vote-service/src/pb"
)

func ConnectGrpcClient(serviceName string, port int) *grpc.ClientConn {
	addr := fmt.Sprintf("localhost:%d", port)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to %s: %v", serviceName, err)
	}
	return conn
}

func main() {

	// connect to other services
	threadConn := ConnectGrpcClient("thread service", 50056)
	defer threadConn.Close()
	commentConn := ConnectGrpcClient("comment service", 50057)
	defer commentConn.Close()

	// create social service with database service
	voteService := &server.VoteServer{
		ThreadClient:  threadpb.NewThreadServiceClient(threadConn),
		CommentClient: commentpb.NewCommentServiceClient(commentConn),
	}

	// start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterVoteServiceServer(grpcServer, voteService)

	log.Println("gRPC server is listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

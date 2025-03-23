package main

import (
	"fmt"
	"log"
	"net"
	server "search-service/src"
	pb "search-service/src/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func connectGrpcClient(serviceName string, port int) *grpc.ClientConn {
	addr := fmt.Sprintf("localhost:%d", port)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to %s: %v", serviceName, err)
	}
	return conn
}

func main() {
	// connect to other services
	userConn := connectGrpcClient("user service", 50052)
	defer userConn.Close()
	communityConn := connectGrpcClient("community service", 50053)
	defer communityConn.Close()
	threadConn := connectGrpcClient("thread service", 50054)
	defer threadConn.Close()

	// create search server with clients
	searchServer := &server.SearchServer{
		UserClient:      userpb.NewUserServiceClient(userConn),
		CommunityClient: communitypb.NewCommunityServiceClient(communityConn),
		ThreadClient:    threadpb.NewThreadServiceClient(threadConn),
	}

	// start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSearchServiceServer(grpcServer, searchServer)

	log.Println("search service is listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

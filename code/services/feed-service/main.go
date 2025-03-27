package main

import (
	dbpb "db-service/src/pb"
	server "feed-service/src"
	pb "feed-service/src/pb"
	"fmt"
	"log"
	"net"

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
	// Connect to database service
	dbConn := connectGrpcClient("database service", 50055)
	defer dbConn.Close()

	// Create feed service with database client
	feedService := &server.FeedServer{
		DBClient: dbpb.NewDBServiceClient(dbConn),
	}

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterFeedServiceServer(grpcServer, feedService)

	log.Println("gRPC server is listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

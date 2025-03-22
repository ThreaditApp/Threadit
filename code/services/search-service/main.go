package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	server "search-service/src"
	pb "search-service/src/pb"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSearchServiceServer(grpcServer, &server.SearchServer{})

	log.Println("gRPC server is listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

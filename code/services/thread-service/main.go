package main

import (
	"fmt"
	"log"
	"net"
	server "thread-service/src"
	pb "thread-service/src/pb"

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
	// connect to database service
	dbConn := connectGrpcClient("database service", 50055)
	defer dbConn.Close()

	// create thread service with database service
	threadService := &server.ThreadServer{
		DBClient: dbpb.NewDBServiceClient(dbConn),
	}

	// start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterThreadServiceServer(grpcServer, threadService)

	log.Println("gRPC server is listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

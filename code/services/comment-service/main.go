package main

import (
    "fmt"
    "log"
    "net"
	pb "comment-service/src/pb"
    server "comment-service/src"

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

    // Create comment service with database client
    commentService := &server.CommentServer{
        DBClient: pb.NewDBServiceClient(dbConn),
    }

    // Start gRPC server
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterCommentServiceServer(grpcServer, commentService)

    log.Println("gRPC server is listening on :50051")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
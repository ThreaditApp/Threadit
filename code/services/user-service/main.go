package main

import (
	"fmt"
	dbpb "gen/db-service/pb"
	userpb "gen/user-service/pb"
	"log"
	"net"
	"os"
	server "user-service/src"

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
	dbConn := connectGrpcClient("database service", "DATABASE_SERVICE_PORT")
	defer dbConn.Close()

	DBService := &server.UserServer{
		DBClient: dbpb.NewDBServiceClient(dbConn),
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
	userpb.RegisterUserServiceServer(grpcServer, DBService)

	log.Printf("gRPC server is listening on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

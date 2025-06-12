package main

import (
	"fmt"
	communitypb "gen/community-service/pb"
	searchpb "gen/search-service/pb"
	threadpb "gen/thread-service/pb"
	"log"
	"net"
	"os"
	"runtime"
	server "search-service/src"

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
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024*1024*500), // 500MB
			grpc.MaxCallSendMsgSize(1024*1024*500), // 500MB
		),
	)
	if err != nil {
		log.Fatalf("failed to connect to %s: %v", addr, err)
	}
	return conn
}

func main() {
	// Set maximum number of CPUs to use
	runtime.GOMAXPROCS(runtime.NumCPU())

	// connect to other services
	communityConn := connectGrpcClient("COMMUNITY_SERVICE_HOST", "COMMUNITY_SERVICE_PORT")
	defer communityConn.Close()
	threadConn := connectGrpcClient("THREAD_SERVICE_HOST", "THREAD_SERVICE_PORT")
	defer threadConn.Close()

	// create search server with clients
	searchServer := &server.SearchServer{
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

	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(1024*1024*500), // 500MB
		grpc.MaxSendMsgSize(1024*1024*500), // 500MB
	)
	searchpb.RegisterSearchServiceServer(grpcServer, searchServer)

	log.Printf("gRPC server is listening on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

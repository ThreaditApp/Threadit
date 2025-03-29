package main

import (
	"context"
	"fmt"
	communitypb "gen/community-service/pb"
	threadpb "gen/thread-service/pb"
	"strings"

	// userpb "gen/user-service/pb"
	// commentpb "gen/comment-service/pb"
	// votepb "gen/vote-service/pb"
	// socialpb "gen/social-service/pb"
	// searchpb "gen/search-service/pb"
	// feedpb "gen/feed-service/pb"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func getGrpcServerAddress(hostEnvVar string, portEnvVar string) string {
	host := os.Getenv(hostEnvVar)
	if host == "" {
		log.Fatalf("missing %s env var", hostEnvVar)
	}
	port := os.Getenv(portEnvVar)
	if port == "" {
		log.Fatalf("missing %s env var", portEnvVar)
	}
	return fmt.Sprintf("%s:%s", host, port)
}

func matchHeader(key string) (string, bool) {
	if key == "X-User-Id" {
		return "x-user-id", true
	}
	return key, false
}

func handleError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) { // TODO: Improve error conversion
	if st, ok := status.FromError(err); ok {
		if strings.Contains(st.Message(), "no documents in result") {
			http.Error(w, "", http.StatusNotFound)
			return
		}
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func main() {
	gwmux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(matchHeader),
		runtime.WithErrorHandler(handleError),
	)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := communitypb.RegisterCommunityServiceHandlerFromEndpoint(context.Background(), gwmux, getGrpcServerAddress("COMMUNITY_SERVICE_HOST", "COMMUNITY_SERVICE_PORT"), opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC gateway: %v", err)
	}

	err = threadpb.RegisterThreadServiceHandlerFromEndpoint(context.Background(), gwmux, getGrpcServerAddress("THREAD_SERVICE_HOST", "THREAD_SERVICE_PORT"), opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC gateway: %v", err)
	}

	// err = userpb.RegisterUserServiceHandlerFromEndpoint(context.Background(), gwmux, getGrpcServerAddress("USER_SERVICE_HOST", "USER_SERVICE_PORT"), opts)
	// if err != nil {
	// 	log.Fatalf("Failed to register gRPC gateway: %v", err)
	// }

	// err = commentpb.RegisterCommentServiceHandlerFromEndpoint(context.Background(), gwmux, getGrpcServerAddress("COMMENT_SERVICE_HOST", "COMMENT_SERVICE_PORT"), opts)
	// if err != nil {
	// 	log.Fatalf("Failed to register gRPC gateway: %v", err)
	// }

	// err = votepb.RegisterVoteServiceHandlerFromEndpoint(context.Background(), gwmux, getGrpcServerAddress("VOTE_SERVICE_HOST", "VOTE_SERVICE_PORT"), opts)
	// if err != nil {
	// 	log.Fatalf("Failed to register gRPC gateway: %v", err)
	// }

	// err = socialpb.RegisterSocialServiceHandlerFromEndpoint(context.Background(), gwmux, getGrpcServerAddress("SOCIAL_SERVICE_HOST", "SOCIAL_SERVICE_PORT"), opts)
	// if err != nil {
	// 	log.Fatalf("Failed to register gRPC gateway: %v", err)
	// }

	// err = searchpb.RegisterSearchServiceHandlerFromEndpoint(context.Background(), gwmux, getGrpcServerAddress("SEARCH_SERVICE_HOST", "SEARCH_SERVICE_PORT"), opts)
	// if err != nil {
	// 	log.Fatalf("Failed to register gRPC gateway: %v", err)
	// }

	// err = feedpb.RegisterFeedServiceHandlerFromEndpoint(context.Background(), gwmux, getGrpcServerAddress("FEED_SERVICE_HOST", "FEED_SERVICE_PORT"), opts)
	// if err != nil {
	// 	log.Fatalf("Failed to register gRPC gateway: %v", err)
	// }

	http.Handle("/", gwmux)

	port := os.Getenv("GRPC_GATEWAY_PORT")
	if port == "" {
		log.Fatalf("missing GRPC_GATEWAY_PORT env var")
	}

	fmt.Printf("gRPC Gateway server listening on :%s", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

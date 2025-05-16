package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
	"your-module/code/services/auth"
)

type AuthMiddleware struct {
	keycloak *auth.KeycloakClient
}

func NewAuthMiddleware(config auth.KeycloakConfig) (*AuthMiddleware, error) {
	kc, err := auth.NewKeycloakClient(config)
	if err != nil {
		return nil, err
	}
	return &AuthMiddleware{keycloak: kc}, nil
}

func (am *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for public endpoints
		if isPublicEndpoint(r.URL.Path, r.Method) {
			next.ServeHTTP(w, r)
			return
		}

		// Extract token from Authorization header
		token, err := auth.ExtractBearerToken(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate token
		claims, err := am.keycloak.ValidateToken(r.Context(), token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Check required roles for protected endpoints
		if !hasRequiredRole(r.URL.Path, claims) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), "user_claims", claims)
		
		// Forward token to gRPC services
		md := metadata.Pairs("authorization", "Bearer "+token)
		ctx = metadata.NewOutgoingContext(ctx, md)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func isPublicEndpoint(path, method string) bool {
	// Auth endpoints are always public
	authPaths := []string{
		"/auth/login",
		"/auth/register",
		"/auth/logout",
	}
	for _, ap := range authPaths {
		if path == ap {
			return true
		}
	}

	// Only GET requests can be public for these paths
	if method != http.MethodGet {
		return false
	}

	publicGetPaths := []string{
		"/communities",
		"/threads",
		"/comments",
		"/search",
		"/search/thread",
		"/search/community",
		"/popular/threads",
		"/popular/comments",
	}

	// Check exact matches for list endpoints
	for _, pp := range publicGetPaths {
		if path == pp {
			return true
		}
	}

	// Check id based paths
	idBasedPaths := []string{
		"/communities/",
		"/threads/",
		"/comments/",
	}

	for _, pp := range idBasedPaths {
		if strings.HasPrefix(path, pp) && path != pp {
			return true
		}
	}

	return false
}

func hasRequiredRole(path string, claims *auth.TokenClaims) bool {
	roleRequirements := map[string]string{
		// Communities
		"POST /communities": "user",
		"PATCH /communities/": "moderator",
		"DELETE /communities/": "moderator",

		// Threads
		"POST /threads": "user",
		"PATCH /threads/": "user",
		"DELETE /threads/": "user",

		// Comment sdpoints
		"POST /comments": "user",
		"PATCH /comments/": "user",
		"DELETE /comments/": "user",

		// Votes
		"POST /votes/thread/": "user",
		"POST /votes/comment/": "user",

		// Admin
		"POST /admin/": "admin",
		"PUT /admin/": "admin",
		"DELETE /admin/": "admin",
	}

	// Check each role requirement
	for pathPattern, requiredRole := range roleRequirements {
		parts := strings.SplitN(pathPattern, " ", 2)
		method, pattern := parts[0], parts[1]
		if strings.HasPrefix(path, pattern) {
			return claims.RealmAccess.Roles != nil && contains(claims.RealmAccess.Roles, requiredRole)
		}
	}

	// If no specific role requirement, allow access
	return true
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
} 
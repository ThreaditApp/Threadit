package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrNoToken          = errors.New("no token provided")
	ErrInvalidToken     = errors.New("invalid token")
	ErrInsufficientRole = errors.New("insufficient role")
)

type KeycloakConfig struct {
	Realm        string
	ClientID     string
	ClientSecret string
	KeycloakURL  string
}

type TokenClaims struct {
	jwt.StandardClaims
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
}

type KeycloakClient struct {
	config KeycloakConfig
	keys   map[string]interface{}
}

func NewKeycloakClient(config KeycloakConfig) (*KeycloakClient, error) {
	kc := &KeycloakClient{
		config: config,
		keys:   make(map[string]interface{}),
	}
	if err := kc.fetchKeys(); err != nil {
		return nil, err
	}
	return kc, nil
}

func (kc *KeycloakClient) fetchKeys() error {
	resp, err := http.Get(fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", kc.config.KeycloakURL, kc.config.Realm))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var jwks struct {
		Keys []struct {
			Kid string                 `json:"kid"`
			Kty string                 `json:"kty"`
			Alg string                 `json:"alg"`
			Use string                 `json:"use"`
			N   string                 `json:"n"`
			E   string                 `json:"e"`
		} `json:"keys"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return err
	}

	for _, key := range jwks.Keys {
		kc.keys[key.Kid] = key
	}

	return nil
}

func (kc *KeycloakClient) ValidateToken(ctx context.Context, tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if kid, ok := token.Header["kid"].(string); ok {
			if key, exists := kc.keys[kid]; exists {
				return key, nil
			}
		}
		return nil, ErrInvalidToken
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

func (kc *KeycloakClient) HasRole(claims *TokenClaims, requiredRole string) bool {
	for _, role := range claims.RealmAccess.Roles {
		if role == requiredRole {
			return true
		}
	}
	return false
}

func ExtractBearerToken(header string) (string, error) {
	if header == "" {
		return "", ErrNoToken
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", ErrInvalidToken
	}

	return parts[1], nil
} 
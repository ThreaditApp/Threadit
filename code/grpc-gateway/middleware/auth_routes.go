package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type AuthHandler struct {
	keycloakURL    string
	clientID       string
	clientSecret   string
	realm          string
	adminClientID  string
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"error_description,omitempty"`
}

func NewAuthHandler(keycloakURL, clientID, clientSecret, realm string) *AuthHandler {
	return &AuthHandler{
		keycloakURL:    keycloakURL,
		clientID:       clientID,
		clientSecret:   clientSecret,
		realm:         realm,
		adminClientID: "admin-cli",
	}
}

func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/auth/register", h.handleRegister)
	mux.HandleFunc("/auth/login", h.handleLogin)
	mux.HandleFunc("/auth/logout", h.handleLogout)
}

func (h *AuthHandler) validateRegistration(req *RegisterRequest) error {
	// Username validation
	if len(req.Username) < 3 || len(req.Username) > 30 {
		return fmt.Errorf("username must be between 3 and 30 characters")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(req.Username) {
		return fmt.Errorf("username can only contain letters, numbers, underscores, and hyphens")
	}

	// Email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return fmt.Errorf("invalid email format")
	}

	// Password validation
	if len(req.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(req.Password) {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(req.Password) {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(req.Password) {
		return fmt.Errorf("password must contain at least one number")
	}

	return nil
}

func (h *AuthHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate registration data
	if err := h.validateRegistration(&req); err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create user in Keycloak
	keycloakURL := fmt.Sprintf("%s/auth/admin/realms/%s/users", h.keycloakURL, h.realm)
	userData := map[string]interface{}{
		"username": req.Username,
		"email":    req.Email,
		"enabled":  true,
		"credentials": []map[string]interface{}{
			{
				"type":      "password",
				"value":     req.Password,
				"temporary": false,
			},
		},
	}

	jsonData, err := json.Marshal(userData)
	if err != nil {
		h.sendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get admin token
	adminToken, err := h.getAdminToken()
	if err != nil {
		h.sendError(w, "Failed to authenticate with Keycloak", http.StatusInternalServerError)
		return
	}

	request, err := http.NewRequest(http.MethodPost, keycloakURL, strings.NewReader(string(jsonData)))
	if err != nil {
		h.sendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+adminToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(request)
	if err != nil {
		h.sendError(w, "Failed to register user", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		h.sendError(w, fmt.Sprintf("Failed to register user: %s", string(body)), resp.StatusCode)
		return
	}

	// Auto login after registration
	h.performLogin(w, req.Username, req.Password)
}

func (h *AuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.performLogin(w, req.Username, req.Password)
}

func (h *AuthHandler) performLogin(w http.ResponseWriter, username, password string) {
	tokenURL := fmt.Sprintf("%s/auth/realms/%s/protocol/openid-connect/token", h.keycloakURL, h.realm)
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", h.clientID)
	data.Set("client_secret", h.clientSecret)
	data.Set("username", username)
	data.Set("password", password)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.PostForm(tokenURL, data)
	if err != nil {
		h.sendError(w, "Failed to login", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		h.sendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Forward Keycloak response (tokens) to client
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}

func (h *AuthHandler) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		h.sendError(w, "No token provided", http.StatusBadRequest)
		return
	}

	logoutURL := fmt.Sprintf("%s/auth/realms/%s/protocol/openid-connect/logout", h.keycloakURL, h.realm)
	request, err := http.NewRequest(http.MethodPost, logoutURL, nil)
	if err != nil {
		h.sendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	request.Header.Set("Authorization", token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(request)
	if err != nil {
		h.sendError(w, "Failed to logout", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		h.sendError(w, "Failed to logout", resp.StatusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

func (h *AuthHandler) getAdminToken() (string, error) {
	tokenURL := fmt.Sprintf("%s/auth/realms/master/protocol/openid-connect/token", h.keycloakURL)
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", h.adminClientID)
	data.Set("client_secret", h.clientSecret)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.PostForm(tokenURL, data)
	if err != nil {
		return "", fmt.Errorf("failed to get admin token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get admin token: status %d", resp.StatusCode)
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode admin token response: %v", err)
	}

	return result.AccessToken, nil
}

func (h *AuthHandler) sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:       http.StatusText(status),
		Description: message,
	})
} 
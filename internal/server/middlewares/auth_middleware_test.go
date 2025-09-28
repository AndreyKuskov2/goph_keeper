package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"goph_keeper/internal/server/config"
	"goph_keeper/pkg/jwt"
)

func TestJwtAuthValidator(t *testing.T) {
	cfg := &config.Config{
		Application: &config.Application{
			SecretToken: "test-secret-token",
		},
	}

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "no authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "no authorization token",
		},
		{
			name:           "invalid authorization header format",
			authHeader:     "InvalidToken",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "invalid authorization header format",
		},
		{
			name:           "invalid bearer format",
			authHeader:     "Bearer",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "invalid authorization header format",
		},
		{
			name:           "empty token",
			authHeader:     "Bearer ",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "empty token",
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "invalid token",
		},
		{
			name:           "valid token",
			authHeader:     "Bearer " + createValidToken(t, cfg.Application.SecretToken),
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create middleware
			middleware := JwtAuthValidator(cfg)

			// Create test handler
			handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check if claims are in context
				claims := r.Context().Value(jwt.ContextClaims)
				if claims == nil {
					t.Error("Expected claims to be in context")
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			}))

			// Create request
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Execute handler
			handler.ServeHTTP(w, req)

			// Check status
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check body
			if w.Body.String() != tt.expectedBody {
				t.Errorf("Expected body '%s', got '%s'", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestJwtAuthValidator_WithValidToken(t *testing.T) {
	cfg := &config.Config{
		Application: &config.Application{
			SecretToken: "test-secret-token",
		},
	}

	// Create middleware
	middleware := JwtAuthValidator(cfg)

	// Create test handler that checks claims
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(jwt.ContextClaims)
		if claims == nil {
			t.Error("Expected claims to be in context")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Verify claims structure - JWT claims are typically a struct
		t.Logf("Claims type: %T, value: %+v", claims, claims)
		// For now, just check that claims exist
		// In real implementation, you would check specific fields

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))

	// Create request with valid token
	validToken := createValidToken(t, cfg.Application.SecretToken)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)

	// Create response recorder
	w := httptest.NewRecorder()

	// Execute handler
	handler.ServeHTTP(w, req)

	// Check status
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Check body
	if w.Body.String() != "success" {
		t.Errorf("Expected body 'success', got '%s'", w.Body.String())
	}
}

func TestJwtAuthValidator_TokenWithWhitespace(t *testing.T) {
	cfg := &config.Config{
		Application: &config.Application{
			SecretToken: "test-secret-token",
		},
	}

	// Create middleware
	middleware := JwtAuthValidator(cfg)

	// Create test handler
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))

	// Create request with token that has whitespace
	validToken := createValidToken(t, cfg.Application.SecretToken)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer  "+validToken+"  ")

	// Create response recorder
	w := httptest.NewRecorder()

	// Execute handler
	handler.ServeHTTP(w, req)

	// Should still work because we trim whitespace
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

// Helper function to create a valid JWT token for testing
func createValidToken(t *testing.T, secret string) string {
	token, err := jwt.CreateJwtToken(secret, 1)
	if err != nil {
		t.Fatalf("Failed to create JWT token: %v", err)
	}
	return token
}

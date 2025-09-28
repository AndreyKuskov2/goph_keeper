package http_handlers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"goph_keeper/internal/models"
	"goph_keeper/internal/server/config"
	custom_errors "goph_keeper/internal/server/errors"
)

// Mock service for testing
type mockUserService struct {
	registerUserID int
	registerError  error
	getUserID      int
	getUserError   error
}

func (m *mockUserService) RegisterUserService(ctx context.Context, user models.User) (int, error) {
	return m.registerUserID, m.registerError
}

func (m *mockUserService) GetUserService(ctx context.Context, user models.UserLoginRequest) (int, error) {
	return m.getUserID, m.getUserError
}

func TestNewGophKeeperUserHandlers(t *testing.T) {
	service := &mockUserService{}
	cfg := &config.Config{
		Application: &config.Application{
			SecretToken: "test-secret",
		},
	}

	handler := NewGophKeeperUserHandlers(service, cfg)

	if handler == nil {
		t.Error("Expected handler to be created")
	}

	if handler.service != service {
		t.Error("Expected service to be set correctly")
	}

	if handler.cfg != cfg {
		t.Error("Expected config to be set correctly")
	}
}

func TestGophKeeperUserHandlers_RegisterUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		service        *mockUserService
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "successful registration",
			requestBody: `{"login":"testuser","password":"testpass"}`,
			service: &mockUserService{
				registerUserID: 1,
				registerError:  nil,
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"code":201,"data":{"token":"`,
		},
		{
			name:        "user already exists",
			requestBody: `{"login":"testuser","password":"testpass"}`,
			service: &mockUserService{
				registerUserID: 0,
				registerError:  custom_errors.ErrUserIsExist,
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"code":409,"error":"user is exist"}`,
		},
		{
			name:        "service error",
			requestBody: `{"login":"testuser","password":"testpass"}`,
			service: &mockUserService{
				registerUserID: 0,
				registerError:  errors.New("database error"),
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"code":500,"error":"database error"}`,
		},
		{
			name:        "invalid JSON",
			requestBody: `{"login":"testuser","password":}`,
			service: &mockUserService{
				registerUserID: 0,
				registerError:  nil,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":400,"error":`,
		},
		{
			name:        "missing login",
			requestBody: `{"password":"testpass"}`,
			service: &mockUserService{
				registerUserID: 0,
				registerError:  nil,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":400,"error":"login is required"}`,
		},
		{
			name:        "missing password",
			requestBody: `{"login":"testuser"}`,
			service: &mockUserService{
				registerUserID: 0,
				registerError:  nil,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":400,"error":"password is required"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create handler
			cfg := &config.Config{
				Application: &config.Application{
					SecretToken: "test-secret-token",
				},
			}
			handler := NewGophKeeperUserHandlers(tt.service, cfg)

			// Create request
			req := httptest.NewRequest("POST", "/register", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Call RegisterUserHandler
			handler.RegisterUserHandler(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check content type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
			}

			// Check response body (partial match for token)
			body := w.Body.String()
			if !strings.Contains(body, tt.expectedBody) {
				t.Errorf("Expected body to contain '%s', got '%s'", tt.expectedBody, body)
			}
		})
	}
}

func TestGophKeeperUserHandlers_LoginUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		service        *mockUserService
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "successful login",
			requestBody: `{"login":"testuser","password":"testpass"}`,
			service: &mockUserService{
				getUserID:    1,
				getUserError: nil,
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"code":200,"data":{"token":"`,
		},
		{
			name:        "user not found",
			requestBody: `{"login":"testuser","password":"testpass"}`,
			service: &mockUserService{
				getUserID:    0,
				getUserError: sql.ErrNoRows,
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"code":401,"error":"sql: no rows in result set"}`,
		},
		{
			name:        "invalid data",
			requestBody: `{"login":"testuser","password":"testpass"}`,
			service: &mockUserService{
				getUserID:    0,
				getUserError: custom_errors.ErrInvalidData,
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"code":401,"error":"invalid data"}`,
		},
		{
			name:        "service error",
			requestBody: `{"login":"testuser","password":"testpass"}`,
			service: &mockUserService{
				getUserID:    0,
				getUserError: errors.New("database error"),
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"code":500,"error":"database error"}`,
		},
		{
			name:        "invalid JSON",
			requestBody: `{"login":"testuser","password":}`,
			service: &mockUserService{
				getUserID:    0,
				getUserError: nil,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":400,"error":`,
		},
		{
			name:        "missing login",
			requestBody: `{"password":"testpass"}`,
			service: &mockUserService{
				getUserID:    0,
				getUserError: nil,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":400,"error":"login is required"}`,
		},
		{
			name:        "missing password",
			requestBody: `{"login":"testuser"}`,
			service: &mockUserService{
				getUserID:    0,
				getUserError: nil,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":400,"error":"password is required"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create handler
			cfg := &config.Config{
				Application: &config.Application{
					SecretToken: "test-secret-token",
				},
			}
			handler := NewGophKeeperUserHandlers(tt.service, cfg)

			// Create request
			req := httptest.NewRequest("POST", "/login", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Call LoginUserHandler
			handler.LoginUserHandler(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check content type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
			}

			// Check response body (partial match for token)
			body := w.Body.String()
			if !strings.Contains(body, tt.expectedBody) {
				t.Errorf("Expected body to contain '%s', got '%s'", tt.expectedBody, body)
			}
		})
	}
}

func TestGophKeeperUserHandlers_WithNilService(t *testing.T) {
	cfg := &config.Config{
		Application: &config.Application{
			SecretToken: "test-secret",
		},
	}

	handler := NewGophKeeperUserHandlers(nil, cfg)

	// Test RegisterUserHandler with nil service
	req := httptest.NewRequest("POST", "/register", strings.NewReader(`{"login":"testuser","password":"testpass"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered from panic: %v", r)
		}
	}()

	handler.RegisterUserHandler(w, req)
}

func TestGophKeeperUserHandlers_WithNilConfig(t *testing.T) {
	service := &mockUserService{}

	handler := NewGophKeeperUserHandlers(service, nil)

	// Test RegisterUserHandler with nil config
	req := httptest.NewRequest("POST", "/register", strings.NewReader(`{"login":"testuser","password":"testpass"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered from panic: %v", r)
		}
	}()

	handler.RegisterUserHandler(w, req)
}

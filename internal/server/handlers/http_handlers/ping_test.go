package http_handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Mock service for testing
type mockPingService struct {
	shouldError bool
	errorMsg    string
}

func (m *mockPingService) Ping(ctx context.Context) error {
	if m.shouldError {
		return errors.New(m.errorMsg)
	}
	return nil
}

func TestNewGophKeeperHTTPHandler(t *testing.T) {
	service := &mockPingService{}
	handler := NewGophKeeperHTTPHandler(service)

	if handler == nil {
		t.Error("Expected handler to be created")
	}

	if handler.service != service {
		t.Error("Expected service to be set correctly")
	}
}

func TestGophKeeperHTTPHandler_Ping(t *testing.T) {
	tests := []struct {
		name           string
		service        *mockPingService
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful ping",
			service: &mockPingService{
				shouldError: false,
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"code":200,"data":{"message":"pong"},"message":""}`,
		},
		{
			name: "ping service error",
			service: &mockPingService{
				shouldError: true,
				errorMsg:    "database connection failed",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"code":500,"error":"database connection failed"}`,
		},
		{
			name: "ping service error with empty message",
			service: &mockPingService{
				shouldError: true,
				errorMsg:    "",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"code":500,"error":""}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create handler
			handler := NewGophKeeperHTTPHandler(tt.service)

			// Create request
			req := httptest.NewRequest("GET", "/ping", nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Call Ping method
			handler.Ping(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check content type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
			}

			// Check response body (trim whitespace for comparison)
			body := strings.TrimSpace(w.Body.String())
			expected := strings.TrimSpace(tt.expectedBody)
			if body != expected {
				t.Errorf("Expected body '%s', got '%s'", expected, body)
			}
		})
	}
}

func TestGophKeeperHTTPHandler_Ping_WithContext(t *testing.T) {
	// Create handler with successful service
	service := &mockPingService{shouldError: false}
	handler := NewGophKeeperHTTPHandler(service)

	// Create request with context
	req := httptest.NewRequest("GET", "/ping", nil)

	// Add some context values
	ctx := context.WithValue(req.Context(), "test-key", "test-value")
	req = req.WithContext(ctx)

	// Create response recorder
	w := httptest.NewRecorder()

	// Call Ping method
	handler.Ping(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check response body
	expectedBody := `{"code":200,"data":{"message":"pong"},"message":""}`
	body := w.Body.String()
	if body != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, body)
	}
}

func TestGophKeeperHTTPHandler_Ping_DifferentHTTPMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run("method_"+method, func(t *testing.T) {
			// Create handler with successful service
			service := &mockPingService{shouldError: false}
			handler := NewGophKeeperHTTPHandler(service)

			// Create request with specific method
			req := httptest.NewRequest(method, "/ping", nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Call Ping method
			handler.Ping(w, req)

			// Check status code (should always be OK for successful ping)
			if w.Code != http.StatusOK {
				t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
			}

			// Check response body
			expectedBody := `{"code":200,"data":{"message":"pong"},"message":""}`
			body := w.Body.String()
			if body != expectedBody {
				t.Errorf("Expected body '%s', got '%s'", expectedBody, body)
			}
		})
	}
}

func TestGophKeeperHTTPHandler_Ping_WithNilService(t *testing.T) {
	// Create handler with nil service
	handler := NewGophKeeperHTTPHandler(nil)

	// Create request
	req := httptest.NewRequest("GET", "/ping", nil)

	// Create response recorder
	w := httptest.NewRecorder()

	// Call Ping method - this should panic or handle gracefully
	defer func() {
		if r := recover(); r != nil {
			// Expected behavior when service is nil
			t.Logf("Recovered from panic: %v", r)
		}
	}()

	handler.Ping(w, req)
}

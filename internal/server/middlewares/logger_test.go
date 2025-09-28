package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLoggerMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "GET request",
			method:         "GET",
			path:           "/test",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name:           "POST request",
			method:         "POST",
			path:           "/api/v1/users",
			expectedStatus: http.StatusCreated,
			expectedBody:   "created",
		},
		{
			name:           "PUT request",
			method:         "PUT",
			path:           "/api/v1/credentials/1",
			expectedStatus: http.StatusOK,
			expectedBody:   "updated",
		},
		{
			name:           "DELETE request",
			method:         "DELETE",
			path:           "/api/v1/credentials/1",
			expectedStatus: http.StatusOK,
			expectedBody:   "deleted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create middleware
			middleware := LoggerMiddleware

			// Create test handler
			handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.expectedStatus)
				w.Write([]byte(tt.expectedBody))
			}))

			// Create request
			req := httptest.NewRequest(tt.method, tt.path, nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Record start time
			start := time.Now()

			// Execute handler
			handler.ServeHTTP(w, req)

			// Record end time
			duration := time.Since(start)

			// Check status
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check body
			if w.Body.String() != tt.expectedBody {
				t.Errorf("Expected body '%s', got '%s'", tt.expectedBody, w.Body.String())
			}

			// Check that duration is reasonable (should be very fast for test)
			if duration > time.Second {
				t.Errorf("Request took too long: %v", duration)
			}
		})
	}
}

func TestLoggerMiddleware_WithDifferentStatusCodes(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "200 OK",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name:           "201 Created",
			expectedStatus: http.StatusCreated,
			expectedBody:   "created",
		},
		{
			name:           "400 Bad Request",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "bad request",
		},
		{
			name:           "401 Unauthorized",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "unauthorized",
		},
		{
			name:           "404 Not Found",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "not found",
		},
		{
			name:           "500 Internal Server Error",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create middleware
			middleware := LoggerMiddleware

			// Create test handler
			handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.expectedStatus)
				w.Write([]byte(tt.expectedBody))
			}))

			// Create request
			req := httptest.NewRequest("GET", "/test", nil)

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

func TestLoggerMiddleware_WithRequestBody(t *testing.T) {
	// Create middleware
	middleware := LoggerMiddleware

	// Create test handler that reads request body
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read request body
		body := make([]byte, 1024)
		_, _ = r.Body.Read(body)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))

	// Create request with body
	req := httptest.NewRequest("POST", "/test", nil)
	req.Body = http.NoBody // Use NoBody for testing

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

func TestLoggerMiddleware_WithHeaders(t *testing.T) {
	// Create middleware
	middleware := LoggerMiddleware

	// Create test handler
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if headers are preserved
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))

	// Create request with headers
	req := httptest.NewRequest("POST", "/test", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token")

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

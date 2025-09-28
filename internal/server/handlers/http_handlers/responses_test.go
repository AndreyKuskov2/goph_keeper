package http_handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestResponseError(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		message  string
		expected string
	}{
		{
			name:     "400 Bad Request",
			code:     http.StatusBadRequest,
			message:  "invalid request",
			expected: `{"code":400,"error":"invalid request"}`,
		},
		{
			name:     "401 Unauthorized",
			code:     http.StatusUnauthorized,
			message:  "unauthorized access",
			expected: `{"code":401,"error":"unauthorized access"}`,
		},
		{
			name:     "404 Not Found",
			code:     http.StatusNotFound,
			message:  "resource not found",
			expected: `{"code":404,"error":"resource not found"}`,
		},
		{
			name:     "500 Internal Server Error",
			code:     http.StatusInternalServerError,
			message:  "internal server error",
			expected: `{"code":500,"error":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest("GET", "/test", nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Call responseError
			responseError(w, req, tt.code, tt.message)

			// Check status code
			if w.Code != tt.code {
				t.Errorf("Expected status code %d, got %d", tt.code, w.Code)
			}

			// Check content type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
			}

			// Check response body (trim whitespace for comparison)
			body := strings.TrimSpace(w.Body.String())
			expected := strings.TrimSpace(tt.expected)
			if body != expected {
				t.Errorf("Expected body '%s', got '%s'", expected, body)
			}
		})
	}
}

func TestResponseOK(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		data     any
		message  string
		expected string
	}{
		{
			name:     "200 OK with string data",
			code:     http.StatusOK,
			data:     "success",
			message:  "operation completed",
			expected: `{"code":200,"data":"success","message":"operation completed"}`,
		},
		{
			name:     "201 Created with object data",
			code:     http.StatusCreated,
			data:     map[string]string{"id": "123"},
			message:  "resource created",
			expected: `{"code":201,"data":{"id":"123"},"message":"resource created"}`,
		},
		{
			name:     "200 OK with nil data",
			code:     http.StatusOK,
			data:     nil,
			message:  "operation completed",
			expected: `{"code":200,"data":null,"message":"operation completed"}`,
		},
		{
			name:     "200 OK with array data",
			code:     http.StatusOK,
			data:     []string{"item1", "item2"},
			message:  "list retrieved",
			expected: `{"code":200,"data":["item1","item2"],"message":"list retrieved"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest("GET", "/test", nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Call responseOK
			responseOK(w, req, tt.code, tt.data, tt.message)

			// Check status code
			if w.Code != tt.code {
				t.Errorf("Expected status code %d, got %d", tt.code, w.Code)
			}

			// Check content type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
			}

			// Check response body (trim whitespace for comparison)
			body := strings.TrimSpace(w.Body.String())
			expected := strings.TrimSpace(tt.expected)
			if body != expected {
				t.Errorf("Expected body '%s', got '%s'", expected, body)
			}
		})
	}
}

func TestResponseStructs(t *testing.T) {
	// Test response struct
	resp := response{
		Code:    200,
		Data:    "test data",
		Message: "test message",
	}

	if resp.Code != 200 {
		t.Errorf("Expected Code 200, got %d", resp.Code)
	}
	if resp.Data != "test data" {
		t.Errorf("Expected Data 'test data', got '%v'", resp.Data)
	}
	if resp.Message != "test message" {
		t.Errorf("Expected Message 'test message', got '%s'", resp.Message)
	}

	// Test responseBad struct
	respBad := responseBad{
		Code:  400,
		Error: "test error",
	}

	if respBad.Code != 400 {
		t.Errorf("Expected Code 400, got %d", respBad.Code)
	}
	if respBad.Error != "test error" {
		t.Errorf("Expected Error 'test error', got '%s'", respBad.Error)
	}
}

func TestResponseError_WithEmptyMessage(t *testing.T) {
	// Create request
	req := httptest.NewRequest("GET", "/test", nil)

	// Create response recorder
	w := httptest.NewRecorder()

	// Call responseError with empty message
	responseError(w, req, http.StatusBadRequest, "")

	// Check status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Check response body
	expected := `{"code":400,"error":""}`
	body := w.Body.String()
	if body != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, body)
	}
}

func TestResponseOK_WithEmptyMessage(t *testing.T) {
	// Create request
	req := httptest.NewRequest("GET", "/test", nil)

	// Create response recorder
	w := httptest.NewRecorder()

	// Call responseOK with empty message
	responseOK(w, req, http.StatusOK, "data", "")

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check response body
	expected := `{"code":200,"data":"data","message":""}`
	body := w.Body.String()
	if body != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, body)
	}
}

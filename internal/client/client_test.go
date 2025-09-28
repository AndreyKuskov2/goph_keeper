package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewCLIClient(t *testing.T) {
	baseURL := "http://localhost:8080"
	client := NewCLIClient(baseURL)

	if client.baseURL != baseURL {
		t.Errorf("Expected baseURL %s, got %s", baseURL, client.baseURL)
	}

	if client.client == nil {
		t.Error("Expected HTTP client to be initialized")
	}

	if client.client.Timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", client.client.Timeout)
	}
}

func TestCLIClient_SetToken(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")
	token := "test-token-123"

	client.SetToken(token)

	if client.token != token {
		t.Errorf("Expected token %s, got %s", token, client.token)
	}
}

func TestCLIClient_SaveToken(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")
	token := "test-token-123"

	// Create temporary directory for testing
	tempDir := t.TempDir()
	originalHomeDir := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHomeDir)

	err := client.SaveToken(token)
	if err != nil {
		t.Fatalf("SaveToken() error = %v", err)
	}

	// Verify token was saved
	tokenFile := filepath.Join(tempDir, ".goph_keeper_token")
	data, err := os.ReadFile(tokenFile)
	if err != nil {
		t.Fatalf("Failed to read token file: %v", err)
	}

	if string(data) != token {
		t.Errorf("Expected token %s, got %s", token, string(data))
	}
}

func TestCLIClient_LoadToken(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")
	token := "test-token-123"

	// Create temporary directory for testing
	tempDir := t.TempDir()
	originalHomeDir := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHomeDir)

	// Save token first
	err := client.SaveToken(token)
	if err != nil {
		t.Fatalf("SaveToken() error = %v", err)
	}

	// Load token
	err = client.LoadToken()
	if err != nil {
		t.Fatalf("LoadToken() error = %v", err)
	}

	if client.token != token {
		t.Errorf("Expected token %s, got %s", token, client.token)
	}
}

func TestCLIClient_LoadToken_FileNotFound(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")

	// Create temporary directory for testing
	tempDir := t.TempDir()
	originalHomeDir := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHomeDir)

	err := client.LoadToken()
	if err == nil {
		t.Error("Expected error when token file doesn't exist")
	}
}

func TestCLIClient_makeRequest(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/test" {
			t.Errorf("Expected path /test, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	client := NewCLIClient(server.URL)

	resp, err := client.makeRequest("GET", "/test", nil, false)
	if err != nil {
		t.Fatalf("makeRequest() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestCLIClient_makeRequest_WithAuth(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got %s", authHeader)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	client := NewCLIClient(server.URL)
	client.SetToken("test-token")

	resp, err := client.makeRequest("GET", "/test", nil, true)
	if err != nil {
		t.Fatalf("makeRequest() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestCLIClient_makeRequest_NoToken(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")

	_, err := client.makeRequest("GET", "/test", nil, true)
	if err == nil {
		t.Error("Expected error when no token is set")
	}

	expectedErr := "authentication token required"
	if err.Error() != expectedErr {
		t.Errorf("Expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestCLIClient_makeRequest_WithBody(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if body["test"] != "data" {
			t.Errorf("Expected body test=data, got %v", body)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	client := NewCLIClient(server.URL)

	body := map[string]string{"test": "data"}
	bodyBytes, _ := json.Marshal(body)

	resp, err := client.makeRequest("POST", "/test", bodyBytes, false)
	if err != nil {
		t.Fatalf("makeRequest() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestCLIClient_handleErrorResponse(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")

	// Create a response with error
	errorResp := ServerErrorResponse{
		Code:  400,
		Error: "test error",
	}
	body, _ := json.Marshal(errorResp)

	// Create mock response
	resp := &http.Response{
		StatusCode: 400,
		Body:       &mockReadCloser{bytes.NewReader(body)},
	}

	err := client.handleErrorResponse(resp)
	if err == nil {
		t.Error("Expected error from handleErrorResponse")
	}

	expectedErr := "server error (status 400): test error"
	if err.Error() != expectedErr {
		t.Errorf("Expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestCLIClient_handleErrorResponse_InvalidJSON(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")

	// Create a response with invalid JSON
	body := []byte("invalid json")

	// Create mock response
	resp := &http.Response{
		StatusCode: 400,
		Body:       &mockReadCloser{bytes.NewReader(body)},
	}

	err := client.handleErrorResponse(resp)
	if err == nil {
		t.Error("Expected error from handleErrorResponse")
	}

	expectedPrefix := "server error (status 400): invalid json"
	if !contains(err.Error(), expectedPrefix) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedPrefix, err.Error())
	}
}

func TestCLIClient_parseServerResponseData(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverResp := ServerResponse{
			Code:    200,
			Data:    Credentials{CredentialsID: 1, Login: "test", Password: "pass", Meta: "meta"},
			Message: "success",
		}
		json.NewEncoder(w).Encode(serverResp)
	}))
	defer server.Close()

	client := NewCLIClient(server.URL)

	resp, err := client.makeRequest("GET", "/test", nil, false)
	if err != nil {
		t.Fatalf("makeRequest() error = %v", err)
	}
	defer resp.Body.Close()

	var credentials Credentials
	err = client.parseServerResponseData(resp, &credentials)
	if err != nil {
		t.Fatalf("parseServerResponseData() error = %v", err)
	}

	if credentials.CredentialsID != 1 {
		t.Errorf("Expected CredentialsID 1, got %d", credentials.CredentialsID)
	}
	if credentials.Login != "test" {
		t.Errorf("Expected Login 'test', got '%s'", credentials.Login)
	}
}

// Helper functions and types for testing
type mockReadCloser struct {
	*bytes.Reader
}

func (m *mockReadCloser) Close() error {
	return nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}

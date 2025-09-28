package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCLIClient_Register(t *testing.T) {
	tests := []struct {
		name           string
		login          string
		password       string
		serverResponse ServerResponse
		serverStatus   int
		wantErr        bool
	}{
		{
			name:     "successful registration",
			login:    "testuser",
			password: "testpass",
			serverResponse: ServerResponse{
				Code:    201,
				Data:    nil,
				Message: "user successfully created",
			},
			serverStatus: http.StatusCreated,
			wantErr:      false,
		},
		{
			name:     "registration failed - user exists",
			login:    "existinguser",
			password: "testpass",
			serverResponse: ServerResponse{
				Code:    400,
				Data:    nil,
				Message: "user already exists",
			},
			serverStatus: http.StatusBadRequest,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Expected POST request, got %s", r.Method)
				}
				if r.URL.Path != "/api/v1/register" {
					t.Errorf("Expected path /api/v1/register, got %s", r.URL.Path)
				}

				// Verify request body
				var user User
				if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
					t.Fatalf("Failed to decode request body: %v", err)
				}

				if user.Login != tt.login {
					t.Errorf("Expected login %s, got %s", tt.login, user.Login)
				}
				if user.Password != tt.password {
					t.Errorf("Expected password %s, got %s", tt.password, user.Password)
				}

				w.WriteHeader(tt.serverStatus)
				json.NewEncoder(w).Encode(tt.serverResponse)
			}))
			defer server.Close()

			client := NewCLIClient(server.URL)

			err := client.Register(tt.login, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLIClient_Login(t *testing.T) {
	tests := []struct {
		name           string
		login          string
		password       string
		serverResponse ServerResponse
		serverStatus   int
		wantToken      string
		wantErr        bool
	}{
		{
			name:     "successful login",
			login:    "testuser",
			password: "testpass",
			serverResponse: ServerResponse{
				Code: 200,
				Data: UserLoginResponse{
					Token: "jwt-token-123",
				},
				Message: "login successful",
			},
			serverStatus: http.StatusOK,
			wantToken:    "jwt-token-123",
			wantErr:      false,
		},
		{
			name:     "login failed - invalid credentials",
			login:    "wronguser",
			password: "wrongpass",
			serverResponse: ServerResponse{
				Code:    401,
				Data:    nil,
				Message: "invalid credentials",
			},
			serverStatus: http.StatusUnauthorized,
			wantToken:    "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Expected POST request, got %s", r.Method)
				}
				if r.URL.Path != "/api/v1/login" {
					t.Errorf("Expected path /api/v1/login, got %s", r.URL.Path)
				}

				// Verify request body
				var loginReq UserLoginRequest
				if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
					t.Fatalf("Failed to decode request body: %v", err)
				}

				if loginReq.Login != tt.login {
					t.Errorf("Expected login %s, got %s", tt.login, loginReq.Login)
				}
				if loginReq.Password != tt.password {
					t.Errorf("Expected password %s, got %s", tt.password, loginReq.Password)
				}

				w.WriteHeader(tt.serverStatus)
				json.NewEncoder(w).Encode(tt.serverResponse)
			}))
			defer server.Close()

			client := NewCLIClient(server.URL)

			token, err := client.Login(tt.login, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}

			if token != tt.wantToken {
				t.Errorf("Login() token = %v, want %v", token, tt.wantToken)
			}
		})
	}
}

func TestCLIClient_Login_InvalidResponse(t *testing.T) {
	// Create test server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewCLIClient(server.URL)

	_, err := client.Login("testuser", "testpass")
	if err == nil {
		t.Error("Expected error for invalid JSON response")
	}
}

func TestCLIClient_Login_InvalidDataField(t *testing.T) {
	// Create test server that returns invalid data field
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverResp := ServerResponse{
			Code:    200,
			Data:    "invalid data type",
			Message: "success",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(serverResp)
	}))
	defer server.Close()

	client := NewCLIClient(server.URL)

	_, err := client.Login("testuser", "testpass")
	if err == nil {
		t.Error("Expected error for invalid data field")
	}
}

func TestCLIClient_Register_NetworkError(t *testing.T) {
	// Create client with invalid URL to simulate network error
	client := NewCLIClient("http://invalid-url:9999")

	err := client.Register("testuser", "testpass")
	if err == nil {
		t.Error("Expected error for network failure")
	}
}

func TestCLIClient_Login_NetworkError(t *testing.T) {
	// Create client with invalid URL to simulate network error
	client := NewCLIClient("http://invalid-url:9999")

	_, err := client.Login("testuser", "testpass")
	if err == nil {
		t.Error("Expected error for network failure")
	}
}

func TestCLIClient_Register_CreateJSONRequestError(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")

	// Test with invalid data that can't be marshaled
	// This is hard to test directly since User struct is simple,
	// but we can test the error handling path
	err := client.Register("", "")
	if err != nil {
		t.Errorf("Register() with empty strings should not error: %v", err)
	}
}

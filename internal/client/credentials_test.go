package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCLIClient_CreateCredentials(t *testing.T) {
	tests := []struct {
		name         string
		login        string
		password     string
		meta         string
		serverStatus int
		wantErr      bool
	}{
		{
			name:         "successful creation",
			login:        "testuser",
			password:     "testpass",
			meta:         "test meta",
			serverStatus: http.StatusCreated,
			wantErr:      false,
		},
		{
			name:         "creation failed",
			login:        "testuser",
			password:     "testpass",
			meta:         "test meta",
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
				if r.URL.Path != "/api/v1/credentials" {
					t.Errorf("Expected path /api/v1/credentials, got %s", r.URL.Path)
				}

				// Verify Authorization header
				authHeader := r.Header.Get("Authorization")
				if authHeader != "Bearer test-token" {
					t.Errorf("Expected Authorization header 'Bearer test-token', got %s", authHeader)
				}

				// Verify request body
				var cred Credentials
				if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
					t.Fatalf("Failed to decode request body: %v", err)
				}

				if cred.Login != tt.login {
					t.Errorf("Expected login %s, got %s", tt.login, cred.Login)
				}
				if cred.Password != tt.password {
					t.Errorf("Expected password %s, got %s", tt.password, cred.Password)
				}
				if cred.Meta != tt.meta {
					t.Errorf("Expected meta %s, got %s", tt.meta, cred.Meta)
				}

				w.WriteHeader(tt.serverStatus)
			}))
			defer server.Close()

			client := NewCLIClient(server.URL)
			client.SetToken("test-token")

			err := client.CreateCredentials(tt.login, tt.password, tt.meta)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLIClient_ListCredentials(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse ServerResponse
		serverStatus   int
		wantErr        bool
	}{
		{
			name: "successful list",
			serverResponse: ServerResponse{
				Code: 200,
				Data: []Credentials{
					{CredentialsID: 1, Login: "user1", Password: "pass1", Meta: "meta1", CreatedAt: time.Now()},
					{CredentialsID: 2, Login: "user2", Password: "pass2", Meta: "meta2", CreatedAt: time.Now()},
				},
				Message: "success",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name: "empty list",
			serverResponse: ServerResponse{
				Code:    200,
				Data:    []Credentials{},
				Message: "success",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name: "server error",
			serverResponse: ServerResponse{
				Code:    500,
				Data:    nil,
				Message: "internal server error",
			},
			serverStatus: http.StatusInternalServerError,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected GET request, got %s", r.Method)
				}
				if r.URL.Path != "/api/v1/credentials" {
					t.Errorf("Expected path /api/v1/credentials, got %s", r.URL.Path)
				}

				w.WriteHeader(tt.serverStatus)
				json.NewEncoder(w).Encode(tt.serverResponse)
			}))
			defer server.Close()

			client := NewCLIClient(server.URL)
			client.SetToken("test-token")

			err := client.ListCredentials()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLIClient_GetCredentials(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		serverResponse ServerResponse
		serverStatus   int
		wantErr        bool
	}{
		{
			name: "successful get",
			id:   "1",
			serverResponse: ServerResponse{
				Code: 200,
				Data: Credentials{
					CredentialsID: 1,
					Login:         "testuser",
					Password:      "testpass",
					Meta:          "test meta",
					CreatedAt:     time.Now(),
				},
				Message: "success",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name: "not found",
			id:   "999",
			serverResponse: ServerResponse{
				Code:    404,
				Data:    nil,
				Message: "not found",
			},
			serverStatus: http.StatusNotFound,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected GET request, got %s", r.Method)
				}
				expectedPath := "/api/v1/credentials/" + tt.id
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.serverStatus)
				json.NewEncoder(w).Encode(tt.serverResponse)
			}))
			defer server.Close()

			client := NewCLIClient(server.URL)
			client.SetToken("test-token")

			err := client.GetCredentials(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLIClient_UpdateCredentials(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		login        string
		password     string
		meta         string
		serverStatus int
		wantErr      bool
	}{
		{
			name:         "successful update",
			id:           "1",
			login:        "updateduser",
			password:     "updatedpass",
			meta:         "updated meta",
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:         "update failed",
			id:           "1",
			login:        "updateduser",
			password:     "updatedpass",
			meta:         "updated meta",
			serverStatus: http.StatusBadRequest,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "PUT" {
					t.Errorf("Expected PUT request, got %s", r.Method)
				}
				expectedPath := "/api/v1/credentials/" + tt.id
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				// Verify request body
				var updateReq UpdateCredentialsRequest
				if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
					t.Fatalf("Failed to decode request body: %v", err)
				}

				if updateReq.Login != tt.login {
					t.Errorf("Expected login %s, got %s", tt.login, updateReq.Login)
				}
				if updateReq.Password != tt.password {
					t.Errorf("Expected password %s, got %s", tt.password, updateReq.Password)
				}
				if updateReq.Meta != tt.meta {
					t.Errorf("Expected meta %s, got %s", tt.meta, updateReq.Meta)
				}

				w.WriteHeader(tt.serverStatus)
			}))
			defer server.Close()

			client := NewCLIClient(server.URL)
			client.SetToken("test-token")

			err := client.UpdateCredentials(tt.id, tt.login, tt.password, tt.meta)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLIClient_DeleteCredentials(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		serverStatus int
		wantErr      bool
	}{
		{
			name:         "successful delete",
			id:           "1",
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:         "delete failed",
			id:           "999",
			serverStatus: http.StatusNotFound,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("Expected DELETE request, got %s", r.Method)
				}
				expectedPath := "/api/v1/credentials/" + tt.id
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.serverStatus)
			}))
			defer server.Close()

			client := NewCLIClient(server.URL)
			client.SetToken("test-token")

			err := client.DeleteCredentials(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLIClient_CreateCredentials_NoToken(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")

	err := client.CreateCredentials("test", "pass", "meta")
	if err == nil {
		t.Error("Expected error when no token is available")
	}
}

func TestCLIClient_ListCredentials_NoToken(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")

	err := client.ListCredentials()
	if err == nil {
		t.Error("Expected error when no token is available")
	}
}

func TestCLIClient_GetCredentials_NoToken(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")

	err := client.GetCredentials("1")
	if err == nil {
		t.Error("Expected error when no token is available")
	}
}

func TestCLIClient_UpdateCredentials_NoToken(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")

	err := client.UpdateCredentials("1", "test", "pass", "meta")
	if err == nil {
		t.Error("Expected error when no token is available")
	}
}

func TestCLIClient_DeleteCredentials_NoToken(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")

	err := client.DeleteCredentials("1")
	if err == nil {
		t.Error("Expected error when no token is available")
	}
}

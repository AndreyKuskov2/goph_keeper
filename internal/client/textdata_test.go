package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCLIClient_CreateTextData(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		meta         string
		serverStatus int
		wantErr      bool
	}{
		{
			name:         "successful creation",
			text:         "This is test text data",
			meta:         "test meta",
			serverStatus: http.StatusCreated,
			wantErr:      false,
		},
		{
			name:         "creation failed",
			text:         "This is test text data",
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
				if r.URL.Path != "/api/v1/text-data" {
					t.Errorf("Expected path /api/v1/text-data, got %s", r.URL.Path)
				}

				// Verify Authorization header
				authHeader := r.Header.Get("Authorization")
				if authHeader != "Bearer test-token" {
					t.Errorf("Expected Authorization header 'Bearer test-token', got %s", authHeader)
				}

				// Verify request body
				var textData TextData
				if err := json.NewDecoder(r.Body).Decode(&textData); err != nil {
					t.Fatalf("Failed to decode request body: %v", err)
				}

				if textData.Text != tt.text {
					t.Errorf("Expected text %s, got %s", tt.text, textData.Text)
				}
				if textData.Meta != tt.meta {
					t.Errorf("Expected meta %s, got %s", tt.meta, textData.Meta)
				}

				w.WriteHeader(tt.serverStatus)
			}))
			defer server.Close()

			client := NewCLIClient(server.URL)
			client.SetToken("test-token")

			err := client.CreateTextData(tt.text, tt.meta)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTextData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLIClient_ListTextData(t *testing.T) {
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
				Data: []TextData{
					{TextDataID: 1, Text: "Text 1", Meta: "meta1", CreatedAt: time.Now()},
					{TextDataID: 2, Text: "Text 2", Meta: "meta2", CreatedAt: time.Now()},
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
				Data:    []TextData{},
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
				if r.URL.Path != "/api/v1/text-data" {
					t.Errorf("Expected path /api/v1/text-data, got %s", r.URL.Path)
				}

				w.WriteHeader(tt.serverStatus)
				json.NewEncoder(w).Encode(tt.serverResponse)
			}))
			defer server.Close()

			client := NewCLIClient(server.URL)
			client.SetToken("test-token")

			err := client.ListTextData()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListTextData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLIClient_GetTextData(t *testing.T) {
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
				Data: TextData{
					TextDataID: 1,
					Text:       "This is test text",
					Meta:       "test meta",
					CreatedAt:  time.Now(),
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
				expectedPath := "/api/v1/text-data/" + tt.id
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.serverStatus)
				json.NewEncoder(w).Encode(tt.serverResponse)
			}))
			defer server.Close()

			client := NewCLIClient(server.URL)
			client.SetToken("test-token")

			err := client.GetTextData(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTextData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLIClient_UpdateTextData(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		text         string
		meta         string
		serverStatus int
		wantErr      bool
	}{
		{
			name:         "successful update",
			id:           "1",
			text:         "Updated text",
			meta:         "Updated meta",
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:         "update failed",
			id:           "1",
			text:         "Updated text",
			meta:         "Updated meta",
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
				expectedPath := "/api/v1/text-data/" + tt.id
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				// Verify request body
				var updateReq UpdateTextDataRequest
				if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
					t.Fatalf("Failed to decode request body: %v", err)
				}

				if updateReq.Text != tt.text {
					t.Errorf("Expected text %s, got %s", tt.text, updateReq.Text)
				}
				if updateReq.Meta != tt.meta {
					t.Errorf("Expected meta %s, got %s", tt.meta, updateReq.Meta)
				}

				w.WriteHeader(tt.serverStatus)
			}))
			defer server.Close()

			client := NewCLIClient(server.URL)
			client.SetToken("test-token")

			err := client.UpdateTextData(tt.id, tt.text, tt.meta)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateTextData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLIClient_DeleteTextData(t *testing.T) {
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
				expectedPath := "/api/v1/text-data/" + tt.id
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.serverStatus)
			}))
			defer server.Close()

			client := NewCLIClient(server.URL)
			client.SetToken("test-token")

			err := client.DeleteTextData(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteTextData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLIClient_CreateTextData_NoToken(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")

	err := client.CreateTextData("test text", "meta")
	if err == nil {
		t.Error("Expected error when no token is available")
	}
}

func TestCLIClient_ListTextData_NoToken(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")

	err := client.ListTextData()
	if err == nil {
		t.Error("Expected error when no token is available")
	}
}

func TestCLIClient_GetTextData_NoToken(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")

	err := client.GetTextData("1")
	if err == nil {
		t.Error("Expected error when no token is available")
	}
}

func TestCLIClient_UpdateTextData_NoToken(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")

	err := client.UpdateTextData("1", "text", "meta")
	if err == nil {
		t.Error("Expected error when no token is available")
	}
}

func TestCLIClient_DeleteTextData_NoToken(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")

	err := client.DeleteTextData("1")
	if err == nil {
		t.Error("Expected error when no token is available")
	}
}

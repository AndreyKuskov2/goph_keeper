package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCLIClient_CreateBinaryData(t *testing.T) {
	tests := []struct {
		name         string
		filePath     string
		meta         string
		serverStatus int
		wantErr      bool
		setupFile    bool
	}{
		{
			name:         "successful creation",
			filePath:     "test_file.txt",
			meta:         "test meta",
			serverStatus: http.StatusCreated,
			wantErr:      false,
			setupFile:    true,
		},
		{
			name:         "creation failed",
			filePath:     "test_file.txt",
			meta:         "test meta",
			serverStatus: http.StatusBadRequest,
			wantErr:      true,
			setupFile:    true,
		},
		{
			name:         "file not found",
			filePath:     "nonexistent_file.txt",
			meta:         "test meta",
			serverStatus: http.StatusCreated,
			wantErr:      true,
			setupFile:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test file if needed
			if tt.setupFile {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, tt.filePath)
				err := os.WriteFile(testFile, []byte("test file content"), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				tt.filePath = testFile
			}

			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Expected POST request, got %s", r.Method)
				}
				if r.URL.Path != "/api/v1/binaries" {
					t.Errorf("Expected path /api/v1/binaries, got %s", r.URL.Path)
				}

				// Verify Authorization header
				authHeader := r.Header.Get("Authorization")
				if authHeader != "Bearer test-token" {
					t.Errorf("Expected Authorization header 'Bearer test-token', got %s", authHeader)
				}

				// Verify request body
				var binaryData BinariesData
				if err := json.NewDecoder(r.Body).Decode(&binaryData); err != nil {
					t.Fatalf("Failed to decode request body: %v", err)
				}

				if binaryData.Meta != tt.meta {
					t.Errorf("Expected meta %s, got %s", tt.meta, binaryData.Meta)
				}

				if tt.setupFile && len(binaryData.BinaryData) == 0 {
					t.Error("Expected binary data to be present")
				}

				w.WriteHeader(tt.serverStatus)
			}))
			defer server.Close()

			client := NewCLIClient(server.URL)
			client.SetToken("test-token")

			err := client.CreateBinaryData(tt.filePath, tt.meta)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateBinaryData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLIClient_ListBinaryData(t *testing.T) {
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
				Data: []BinariesData{
					{BinariesDataID: 1, BinaryData: []byte("test data 1"), Meta: "meta1", CreatedAt: time.Now()},
					{BinariesDataID: 2, BinaryData: []byte("test data 2"), Meta: "meta2", CreatedAt: time.Now()},
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
				Data:    []BinariesData{},
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
				if r.URL.Path != "/api/v1/binaries" {
					t.Errorf("Expected path /api/v1/binaries, got %s", r.URL.Path)
				}

				w.WriteHeader(tt.serverStatus)
				json.NewEncoder(w).Encode(tt.serverResponse)
			}))
			defer server.Close()

			client := NewCLIClient(server.URL)
			client.SetToken("test-token")

			err := client.ListBinaryData()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListBinaryData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLIClient_GetBinaryData(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		outputFile     string
		serverResponse ServerResponse
		serverStatus   int
		wantErr        bool
	}{
		{
			name:       "successful get with custom output file",
			id:         "1",
			outputFile: "downloaded_file.bin",
			serverResponse: ServerResponse{
				Code: 200,
				Data: BinariesData{
					BinariesDataID: 1,
					BinaryData:     []byte("test binary data"),
					Meta:           "test meta",
					CreatedAt:      time.Now(),
				},
				Message: "success",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:       "successful get with default output file",
			id:         "1",
			outputFile: "",
			serverResponse: ServerResponse{
				Code: 200,
				Data: BinariesData{
					BinariesDataID: 1,
					BinaryData:     []byte("test binary data"),
					Meta:           "test meta",
					CreatedAt:      time.Now(),
				},
				Message: "success",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:       "not found",
			id:         "999",
			outputFile: "downloaded_file.bin",
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
			tempDir := t.TempDir()
			outputPath := ""
			if tt.outputFile != "" {
				outputPath = filepath.Join(tempDir, tt.outputFile)
			}

			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected GET request, got %s", r.Method)
				}
				expectedPath := "/api/v1/binaries/" + tt.id
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.serverStatus)
				json.NewEncoder(w).Encode(tt.serverResponse)
			}))
			defer server.Close()

			client := NewCLIClient(server.URL)
			client.SetToken("test-token")

			err := client.GetBinaryData(tt.id, outputPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBinaryData() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify file was created if successful
			if !tt.wantErr && tt.serverStatus == http.StatusOK {
				expectedFile := outputPath
				if expectedFile == "" {
					expectedFile = "binary_" + tt.id + ".bin"
				}

				if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
					t.Errorf("Expected file %s to be created", expectedFile)
				}
			}
		})
	}
}

func TestCLIClient_UpdateBinaryData(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		filePath     string
		meta         string
		serverStatus int
		wantErr      bool
		setupFile    bool
	}{
		{
			name:         "successful update",
			id:           "1",
			filePath:     "updated_file.txt",
			meta:         "updated meta",
			serverStatus: http.StatusOK,
			wantErr:      false,
			setupFile:    true,
		},
		{
			name:         "update failed",
			id:           "1",
			filePath:     "updated_file.txt",
			meta:         "updated meta",
			serverStatus: http.StatusBadRequest,
			wantErr:      true,
			setupFile:    true,
		},
		{
			name:         "file not found",
			id:           "1",
			filePath:     "nonexistent_file.txt",
			meta:         "updated meta",
			serverStatus: http.StatusOK,
			wantErr:      true,
			setupFile:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test file if needed
			if tt.setupFile {
				tempDir := t.TempDir()
				testFile := filepath.Join(tempDir, tt.filePath)
				err := os.WriteFile(testFile, []byte("updated file content"), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				tt.filePath = testFile
			}

			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "PUT" {
					t.Errorf("Expected PUT request, got %s", r.Method)
				}
				expectedPath := "/api/v1/binaries/" + tt.id
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				// Verify request body
				var updateReq UpdateBinariesDataRequest
				if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
					t.Fatalf("Failed to decode request body: %v", err)
				}

				if updateReq.Meta != tt.meta {
					t.Errorf("Expected meta %s, got %s", tt.meta, updateReq.Meta)
				}

				if tt.setupFile && len(updateReq.BinaryData) == 0 {
					t.Error("Expected binary data to be present")
				}

				w.WriteHeader(tt.serverStatus)
			}))
			defer server.Close()

			client := NewCLIClient(server.URL)
			client.SetToken("test-token")

			err := client.UpdateBinaryData(tt.id, tt.filePath, tt.meta)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateBinaryData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLIClient_DeleteBinaryData(t *testing.T) {
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
				expectedPath := "/api/v1/binaries/" + tt.id
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.serverStatus)
			}))
			defer server.Close()

			client := NewCLIClient(server.URL)
			client.SetToken("test-token")

			err := client.DeleteBinaryData(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteBinaryData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLIClient_CreateBinaryData_NoToken(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")

	err := client.CreateBinaryData("test.txt", "meta")
	if err == nil {
		t.Error("Expected error when no token is available")
	}
}

func TestCLIClient_ListBinaryData_NoToken(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")

	err := client.ListBinaryData()
	if err == nil {
		t.Error("Expected error when no token is available")
	}
}

func TestCLIClient_GetBinaryData_NoToken(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")

	err := client.GetBinaryData("1", "output.bin")
	if err == nil {
		t.Error("Expected error when no token is available")
	}
}

func TestCLIClient_UpdateBinaryData_NoToken(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")

	err := client.UpdateBinaryData("1", "test.txt", "meta")
	if err == nil {
		t.Error("Expected error when no token is available")
	}
}

func TestCLIClient_DeleteBinaryData_NoToken(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")

	err := client.DeleteBinaryData("1")
	if err == nil {
		t.Error("Expected error when no token is available")
	}
}

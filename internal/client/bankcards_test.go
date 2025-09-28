package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCLIClient_CreateBankCard(t *testing.T) {
	tests := []struct {
		name           string
		cardNumber     string
		holder         string
		cvc            string
		expirationDate string
		meta           string
		serverStatus   int
		wantErr        bool
	}{
		{
			name:           "successful creation",
			cardNumber:     "1234567890123456",
			holder:         "John Doe",
			cvc:            "123",
			expirationDate: "2025-12-31",
			meta:           "test meta",
			serverStatus:   http.StatusCreated,
			wantErr:        false,
		},
		{
			name:           "creation failed",
			cardNumber:     "1234567890123456",
			holder:         "John Doe",
			cvc:            "123",
			expirationDate: "2025-12-31",
			meta:           "test meta",
			serverStatus:   http.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:           "invalid expiration date",
			cardNumber:     "1234567890123456",
			holder:         "John Doe",
			cvc:            "123",
			expirationDate: "invalid-date",
			meta:           "test meta",
			serverStatus:   http.StatusCreated,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Expected POST request, got %s", r.Method)
				}
				if r.URL.Path != "/api/v1/bank-cards" {
					t.Errorf("Expected path /api/v1/bank-cards, got %s", r.URL.Path)
				}

				// Verify Authorization header
				authHeader := r.Header.Get("Authorization")
				if authHeader != "Bearer test-token" {
					t.Errorf("Expected Authorization header 'Bearer test-token', got %s", authHeader)
				}

				// Verify request body
				var bankCard BankCards
				if err := json.NewDecoder(r.Body).Decode(&bankCard); err != nil {
					t.Fatalf("Failed to decode request body: %v", err)
				}

				if bankCard.CardNumber != tt.cardNumber {
					t.Errorf("Expected cardNumber %s, got %s", tt.cardNumber, bankCard.CardNumber)
				}
				if bankCard.Holder != tt.holder {
					t.Errorf("Expected holder %s, got %s", tt.holder, bankCard.Holder)
				}
				if bankCard.CVC != tt.cvc {
					t.Errorf("Expected CVC %s, got %s", tt.cvc, bankCard.CVC)
				}
				if bankCard.Meta != tt.meta {
					t.Errorf("Expected meta %s, got %s", tt.meta, bankCard.Meta)
				}

				w.WriteHeader(tt.serverStatus)
			}))
			defer server.Close()

			client := NewCLIClient(server.URL)
			client.SetToken("test-token")

			err := client.CreateBankCard(tt.cardNumber, tt.holder, tt.cvc, tt.expirationDate, tt.meta)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateBankCard() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLIClient_ListBankCards(t *testing.T) {
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
				Data: []BankCards{
					{BankCardsID: 1, CardNumber: "1234567890123456", Holder: "John Doe", CVC: "123", Meta: "meta1", CreatedAt: time.Now()},
					{BankCardsID: 2, CardNumber: "9876543210987654", Holder: "Jane Doe", CVC: "456", Meta: "meta2", CreatedAt: time.Now()},
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
				Data:    []BankCards{},
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
				if r.URL.Path != "/api/v1/bank-cards" {
					t.Errorf("Expected path /api/v1/bank-cards, got %s", r.URL.Path)
				}

				w.WriteHeader(tt.serverStatus)
				json.NewEncoder(w).Encode(tt.serverResponse)
			}))
			defer server.Close()

			client := NewCLIClient(server.URL)
			client.SetToken("test-token")

			err := client.ListBankCards()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListBankCards() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLIClient_GetBankCard(t *testing.T) {
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
				Data: BankCards{
					BankCardsID:    1,
					CardNumber:     "1234567890123456",
					Holder:         "John Doe",
					CVC:            "123",
					ExpirationDate: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
					Meta:           "test meta",
					CreatedAt:      time.Now(),
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
				expectedPath := "/api/v1/bank-cards/" + tt.id
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.serverStatus)
				json.NewEncoder(w).Encode(tt.serverResponse)
			}))
			defer server.Close()

			client := NewCLIClient(server.URL)
			client.SetToken("test-token")

			err := client.GetBankCard(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBankCard() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLIClient_UpdateBankCard(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		cardNumber     string
		holder         string
		cvc            string
		expirationDate string
		meta           string
		serverStatus   int
		wantErr        bool
	}{
		{
			name:           "successful update",
			id:             "1",
			cardNumber:     "9876543210987654",
			holder:         "Jane Doe",
			cvc:            "456",
			expirationDate: "2026-06-30",
			meta:           "updated meta",
			serverStatus:   http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "update failed",
			id:             "1",
			cardNumber:     "9876543210987654",
			holder:         "Jane Doe",
			cvc:            "456",
			expirationDate: "2026-06-30",
			meta:           "updated meta",
			serverStatus:   http.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:           "invalid expiration date",
			id:             "1",
			cardNumber:     "9876543210987654",
			holder:         "Jane Doe",
			cvc:            "456",
			expirationDate: "invalid-date",
			meta:           "updated meta",
			serverStatus:   http.StatusOK,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "PUT" {
					t.Errorf("Expected PUT request, got %s", r.Method)
				}
				expectedPath := "/api/v1/bank-cards/" + tt.id
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				// Verify request body
				var updateReq UpdateBankCardsRequest
				if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
					t.Fatalf("Failed to decode request body: %v", err)
				}

				if updateReq.CardNumber != tt.cardNumber {
					t.Errorf("Expected cardNumber %s, got %s", tt.cardNumber, updateReq.CardNumber)
				}
				if updateReq.Holder != tt.holder {
					t.Errorf("Expected holder %s, got %s", tt.holder, updateReq.Holder)
				}
				if updateReq.CVC != tt.cvc {
					t.Errorf("Expected CVC %s, got %s", tt.cvc, updateReq.CVC)
				}
				if updateReq.Meta != tt.meta {
					t.Errorf("Expected meta %s, got %s", tt.meta, updateReq.Meta)
				}

				w.WriteHeader(tt.serverStatus)
			}))
			defer server.Close()

			client := NewCLIClient(server.URL)
			client.SetToken("test-token")

			err := client.UpdateBankCard(tt.id, tt.cardNumber, tt.holder, tt.cvc, tt.expirationDate, tt.meta)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateBankCard() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLIClient_DeleteBankCard(t *testing.T) {
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
				expectedPath := "/api/v1/bank-cards/" + tt.id
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.serverStatus)
			}))
			defer server.Close()

			client := NewCLIClient(server.URL)
			client.SetToken("test-token")

			err := client.DeleteBankCard(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteBankCard() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLIClient_CreateBankCard_NoToken(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")

	err := client.CreateBankCard("1234567890123456", "John Doe", "123", "2025-12-31", "meta")
	if err == nil {
		t.Error("Expected error when no token is available")
	}
}

func TestCLIClient_ListBankCards_NoToken(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")

	err := client.ListBankCards()
	if err == nil {
		t.Error("Expected error when no token is available")
	}
}

func TestCLIClient_GetBankCard_NoToken(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")

	err := client.GetBankCard("1")
	if err == nil {
		t.Error("Expected error when no token is available")
	}
}

func TestCLIClient_UpdateBankCard_NoToken(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")

	err := client.UpdateBankCard("1", "1234567890123456", "John Doe", "123", "2025-12-31", "meta")
	if err == nil {
		t.Error("Expected error when no token is available")
	}
}

func TestCLIClient_DeleteBankCard_NoToken(t *testing.T) {
	client := NewCLIClient("http://localhost:8080")

	err := client.DeleteBankCard("1")
	if err == nil {
		t.Error("Expected error when no token is available")
	}
}

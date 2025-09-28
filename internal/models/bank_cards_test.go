package models

import (
	"net/http"
	"testing"
	"time"
)

func TestBankCards_Bind(t *testing.T) {
	tests := []struct {
		name    string
		card    BankCards
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid bank card",
			card: BankCards{
				CardNumber:     "1234567890123456",
				Holder:         "John Doe",
				CVC:            "123",
				ExpirationDate: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
				Meta:           "test meta",
			},
			wantErr: false,
		},
		{
			name: "empty card number",
			card: BankCards{
				CardNumber:     "",
				Holder:         "John Doe",
				CVC:            "123",
				ExpirationDate: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
				Meta:           "test meta",
			},
			wantErr: true,
			errMsg:  "card_number is required",
		},
		{
			name: "empty holder",
			card: BankCards{
				CardNumber:     "1234567890123456",
				Holder:         "",
				CVC:            "123",
				ExpirationDate: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
				Meta:           "test meta",
			},
			wantErr: true,
			errMsg:  "holder is required",
		},
		{
			name: "empty CVC",
			card: BankCards{
				CardNumber:     "1234567890123456",
				Holder:         "John Doe",
				CVC:            "",
				ExpirationDate: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
				Meta:           "test meta",
			},
			wantErr: true,
			errMsg:  "cvc is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.card.Bind(&http.Request{})
			if (err != nil) != tt.wantErr {
				t.Errorf("BankCards.Bind() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("BankCards.Bind() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestUpdateBankCardsRequest_Bind(t *testing.T) {
	tests := []struct {
		name    string
		card    UpdateBankCardsRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid update request",
			card: UpdateBankCardsRequest{
				CardNumber:     "9876543210987654",
				Holder:         "Jane Doe",
				CVC:            "456",
				ExpirationDate: time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC),
				Meta:           "updated meta",
			},
			wantErr: false,
		},
		{
			name: "empty card number",
			card: UpdateBankCardsRequest{
				CardNumber:     "",
				Holder:         "Jane Doe",
				CVC:            "456",
				ExpirationDate: time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC),
				Meta:           "updated meta",
			},
			wantErr: true,
			errMsg:  "card_number is required",
		},
		{
			name: "empty holder",
			card: UpdateBankCardsRequest{
				CardNumber:     "9876543210987654",
				Holder:         "",
				CVC:            "456",
				ExpirationDate: time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC),
				Meta:           "updated meta",
			},
			wantErr: true,
			errMsg:  "holder is required",
		},
		{
			name: "empty CVC",
			card: UpdateBankCardsRequest{
				CardNumber:     "9876543210987654",
				Holder:         "Jane Doe",
				CVC:            "",
				ExpirationDate: time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC),
				Meta:           "updated meta",
			},
			wantErr: true,
			errMsg:  "cvc is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.card.Bind(&http.Request{})
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateBankCardsRequest.Bind() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("UpdateBankCardsRequest.Bind() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestBankCards_Struct(t *testing.T) {
	now := time.Now()
	expDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	card := BankCards{
		BankCardsID:    1,
		CardNumber:     "1234567890123456",
		Holder:         "John Doe",
		CVC:            "123",
		ExpirationDate: expDate,
		UserID:         1,
		CreatedAt:      now,
		Meta:           "test meta",
	}

	if card.BankCardsID != 1 {
		t.Errorf("Expected BankCardsID 1, got %d", card.BankCardsID)
	}
	if card.CardNumber != "1234567890123456" {
		t.Errorf("Expected CardNumber '1234567890123456', got '%s'", card.CardNumber)
	}
	if card.Holder != "John Doe" {
		t.Errorf("Expected Holder 'John Doe', got '%s'", card.Holder)
	}
	if card.CVC != "123" {
		t.Errorf("Expected CVC '123', got '%s'", card.CVC)
	}
	if !card.ExpirationDate.Equal(expDate) {
		t.Errorf("Expected ExpirationDate %v, got %v", expDate, card.ExpirationDate)
	}
	if card.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", card.UserID)
	}
	if !card.CreatedAt.Equal(now) {
		t.Errorf("Expected CreatedAt %v, got %v", now, card.CreatedAt)
	}
	if card.Meta != "test meta" {
		t.Errorf("Expected Meta 'test meta', got '%s'", card.Meta)
	}
}

func TestUpdateBankCardsRequest_Struct(t *testing.T) {
	expDate := time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)
	updateReq := UpdateBankCardsRequest{
		CardNumber:     "9876543210987654",
		Holder:         "Jane Doe",
		CVC:            "456",
		ExpirationDate: expDate,
		Meta:           "updated meta",
	}

	if updateReq.CardNumber != "9876543210987654" {
		t.Errorf("Expected CardNumber '9876543210987654', got '%s'", updateReq.CardNumber)
	}
	if updateReq.Holder != "Jane Doe" {
		t.Errorf("Expected Holder 'Jane Doe', got '%s'", updateReq.Holder)
	}
	if updateReq.CVC != "456" {
		t.Errorf("Expected CVC '456', got '%s'", updateReq.CVC)
	}
	if !updateReq.ExpirationDate.Equal(expDate) {
		t.Errorf("Expected ExpirationDate %v, got %v", expDate, updateReq.ExpirationDate)
	}
	if updateReq.Meta != "updated meta" {
		t.Errorf("Expected Meta 'updated meta', got '%s'", updateReq.Meta)
	}
}

package models

import (
	"net/http"
	"testing"
	"time"
)

func TestTextData_Bind(t *testing.T) {
	tests := []struct {
		name     string
		textData TextData
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid text data",
			textData: TextData{
				Text: "This is test text",
				Meta: "test meta",
			},
			wantErr: false,
		},
		{
			name: "empty text",
			textData: TextData{
				Text: "",
				Meta: "test meta",
			},
			wantErr: true,
			errMsg:  "login is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.textData.Bind(&http.Request{})
			if (err != nil) != tt.wantErr {
				t.Errorf("TextData.Bind() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("TextData.Bind() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestUpdateTextDataRequest_Bind(t *testing.T) {
	tests := []struct {
		name     string
		textData UpdateTextDataRequest
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid update request",
			textData: UpdateTextDataRequest{
				Text: "Updated text",
				Meta: "Updated meta",
			},
			wantErr: false,
		},
		{
			name: "empty text",
			textData: UpdateTextDataRequest{
				Text: "",
				Meta: "test meta",
			},
			wantErr: true,
			errMsg:  "login is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.textData.Bind(&http.Request{})
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateTextDataRequest.Bind() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("UpdateTextDataRequest.Bind() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestTextData_Struct(t *testing.T) {
	now := time.Now()
	textData := TextData{
		TextDataID: 1,
		Text:       "This is test text",
		UserID:     1,
		CreatedAt:  now,
		Meta:       "test meta",
	}

	if textData.TextDataID != 1 {
		t.Errorf("Expected TextDataID 1, got %d", textData.TextDataID)
	}
	if textData.Text != "This is test text" {
		t.Errorf("Expected Text 'This is test text', got '%s'", textData.Text)
	}
	if textData.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", textData.UserID)
	}
	if !textData.CreatedAt.Equal(now) {
		t.Errorf("Expected CreatedAt %v, got %v", now, textData.CreatedAt)
	}
	if textData.Meta != "test meta" {
		t.Errorf("Expected Meta 'test meta', got '%s'", textData.Meta)
	}
}

func TestUpdateTextDataRequest_Struct(t *testing.T) {
	updateReq := UpdateTextDataRequest{
		Text: "Updated text",
		Meta: "Updated meta",
	}

	if updateReq.Text != "Updated text" {
		t.Errorf("Expected Text 'Updated text', got '%s'", updateReq.Text)
	}
	if updateReq.Meta != "Updated meta" {
		t.Errorf("Expected Meta 'Updated meta', got '%s'", updateReq.Meta)
	}
}

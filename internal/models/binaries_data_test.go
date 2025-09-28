package models

import (
	"net/http"
	"testing"
	"time"
)

func TestBinariesData_Bind(t *testing.T) {
	tests := []struct {
		name    string
		binary  BinariesData
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid binary data with nil data",
			binary: BinariesData{
				BinaryData: nil,
				Meta:       "test meta",
			},
			wantErr: false,
		},
		{
			name: "binary data with actual data",
			binary: BinariesData{
				BinaryData: []byte("test data"),
				Meta:       "test meta",
			},
			wantErr: true,
			errMsg:  "card_number is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.binary.Bind(&http.Request{})
			if (err != nil) != tt.wantErr {
				t.Errorf("BinariesData.Bind() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("BinariesData.Bind() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestUpdateBinariesDataRequest_Bind(t *testing.T) {
	tests := []struct {
		name    string
		binary  UpdateBinariesDataRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid update request with nil data",
			binary: UpdateBinariesDataRequest{
				BinaryData: nil,
				Meta:       "updated meta",
			},
			wantErr: false,
		},
		{
			name: "update request with actual data",
			binary: UpdateBinariesDataRequest{
				BinaryData: []byte("updated data"),
				Meta:       "updated meta",
			},
			wantErr: true,
			errMsg:  "card_number is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.binary.Bind(&http.Request{})
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateBinariesDataRequest.Bind() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("UpdateBinariesDataRequest.Bind() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestBinariesData_Struct(t *testing.T) {
	now := time.Now()
	binaryData := []byte("test binary data")
	binary := BinariesData{
		BinariesDataID: 1,
		BinaryData:     binaryData,
		UserID:         1,
		CreatedAt:      now,
		Meta:           "test meta",
	}

	if binary.BinariesDataID != 1 {
		t.Errorf("Expected BinariesDataID 1, got %d", binary.BinariesDataID)
	}
	if string(binary.BinaryData) != string(binaryData) {
		t.Errorf("Expected BinaryData %v, got %v", binaryData, binary.BinaryData)
	}
	if binary.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", binary.UserID)
	}
	if !binary.CreatedAt.Equal(now) {
		t.Errorf("Expected CreatedAt %v, got %v", now, binary.CreatedAt)
	}
	if binary.Meta != "test meta" {
		t.Errorf("Expected Meta 'test meta', got '%s'", binary.Meta)
	}
}

func TestUpdateBinariesDataRequest_Struct(t *testing.T) {
	updatedData := []byte("updated binary data")
	updateReq := UpdateBinariesDataRequest{
		BinaryData: updatedData,
		Meta:       "updated meta",
	}

	if string(updateReq.BinaryData) != string(updatedData) {
		t.Errorf("Expected BinaryData %v, got %v", updatedData, updateReq.BinaryData)
	}
	if updateReq.Meta != "updated meta" {
		t.Errorf("Expected Meta 'updated meta', got '%s'", updateReq.Meta)
	}
}

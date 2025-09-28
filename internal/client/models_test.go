package client

import (
	"encoding/json"
	"testing"
	"time"
)

func TestParseJSONResponse(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		target  interface{}
		wantErr bool
	}{
		{
			name:    "valid JSON to struct",
			data:    []byte(`{"credentials_id":1,"login":"test","password":"pass","user_id":1,"created_at":"2023-01-01T00:00:00Z","meta":"test meta"}`),
			target:  &Credentials{},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			data:    []byte(`{"invalid": json}`),
			target:  &Credentials{},
			wantErr: true,
		},
		{
			name:    "empty data",
			data:    []byte(``),
			target:  &Credentials{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ParseJSONResponse(tt.data, tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJSONResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateJSONRequest(t *testing.T) {
	tests := []struct {
		name    string
		data    interface{}
		wantErr bool
	}{
		{
			name:    "valid struct",
			data:    Credentials{Login: "test", Password: "pass", Meta: "meta"},
			wantErr: false,
		},
		{
			name:    "nil data",
			data:    nil,
			wantErr: false,
		},
		{
			name:    "channel (unmarshalable)",
			data:    make(chan int),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateJSONRequest(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateJSONRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFormatTime(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		want string
	}{
		{
			name: "standard time",
			time: time.Date(2023, 1, 15, 14, 30, 45, 0, time.UTC),
			want: "2023-01-15 14:30:45",
		},
		{
			name: "zero time",
			time: time.Time{},
			want: "0001-01-01 00:00:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatTime(tt.time); got != tt.want {
				t.Errorf("FormatTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseTime(t *testing.T) {
	tests := []struct {
		name    string
		timeStr string
		wantErr bool
	}{
		{
			name:    "valid date",
			timeStr: "2023-01-15",
			wantErr: false,
		},
		{
			name:    "invalid format",
			timeStr: "15-01-2023",
			wantErr: true,
		},
		{
			name:    "empty string",
			timeStr: "",
			wantErr: true,
		},
		{
			name:    "invalid date",
			timeStr: "2023-13-45",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseTime(tt.timeStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCredentialsStruct(t *testing.T) {
	cred := Credentials{
		CredentialsID: 1,
		Login:         "test",
		Password:      "pass",
		UserID:        1,
		CreatedAt:     time.Now(),
		Meta:          "test meta",
	}

	// Test JSON marshaling/unmarshaling
	data, err := json.Marshal(cred)
	if err != nil {
		t.Fatalf("Failed to marshal credentials: %v", err)
	}

	var unmarshaled Credentials
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal credentials: %v", err)
	}

	if unmarshaled.CredentialsID != cred.CredentialsID {
		t.Errorf("CredentialsID mismatch: got %d, want %d", unmarshaled.CredentialsID, cred.CredentialsID)
	}
	if unmarshaled.Login != cred.Login {
		t.Errorf("Login mismatch: got %s, want %s", unmarshaled.Login, cred.Login)
	}
}

func TestUserLoginRequestStruct(t *testing.T) {
	req := UserLoginRequest{
		Login:    "testuser",
		Password: "testpass",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal UserLoginRequest: %v", err)
	}

	var unmarshaled UserLoginRequest
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal UserLoginRequest: %v", err)
	}

	if unmarshaled.Login != req.Login {
		t.Errorf("Login mismatch: got %s, want %s", unmarshaled.Login, req.Login)
	}
	if unmarshaled.Password != req.Password {
		t.Errorf("Password mismatch: got %s, want %s", unmarshaled.Password, req.Password)
	}
}

func TestServerResponseStruct(t *testing.T) {
	resp := ServerResponse{
		Code:    200,
		Data:    "test data",
		Message: "success",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal ServerResponse: %v", err)
	}

	var unmarshaled ServerResponse
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ServerResponse: %v", err)
	}

	if unmarshaled.Code != resp.Code {
		t.Errorf("Code mismatch: got %d, want %d", unmarshaled.Code, resp.Code)
	}
	if unmarshaled.Message != resp.Message {
		t.Errorf("Message mismatch: got %s, want %s", unmarshaled.Message, resp.Message)
	}
}

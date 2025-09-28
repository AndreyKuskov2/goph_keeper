package models

import (
	"net/http"
	"testing"
	"time"
)

func TestCredentials_Bind(t *testing.T) {
	tests := []struct {
		name    string
		creds   Credentials
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid credentials",
			creds: Credentials{
				Login:    "testuser",
				Password: "testpass",
				Meta:     "test meta",
			},
			wantErr: false,
		},
		{
			name: "empty login",
			creds: Credentials{
				Login:    "",
				Password: "testpass",
				Meta:     "test meta",
			},
			wantErr: true,
			errMsg:  "login is required",
		},
		{
			name: "empty password",
			creds: Credentials{
				Login:    "testuser",
				Password: "",
				Meta:     "test meta",
			},
			wantErr: true,
			errMsg:  "password is required",
		},
		{
			name: "both empty",
			creds: Credentials{
				Login:    "",
				Password: "",
				Meta:     "test meta",
			},
			wantErr: true,
			errMsg:  "login is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.creds.Bind(&http.Request{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Credentials.Bind() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("Credentials.Bind() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestUpdateCredentialsRequest_Bind(t *testing.T) {
	tests := []struct {
		name    string
		creds   UpdateCredentialsRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid update request",
			creds: UpdateCredentialsRequest{
				Login:    "testuser",
				Password: "testpass",
				Meta:     "test meta",
			},
			wantErr: false,
		},
		{
			name: "empty login",
			creds: UpdateCredentialsRequest{
				Login:    "",
				Password: "testpass",
				Meta:     "test meta",
			},
			wantErr: true,
			errMsg:  "login is required",
		},
		{
			name: "empty password",
			creds: UpdateCredentialsRequest{
				Login:    "testuser",
				Password: "",
				Meta:     "test meta",
			},
			wantErr: true,
			errMsg:  "password is required",
		},
		{
			name: "both empty",
			creds: UpdateCredentialsRequest{
				Login:    "",
				Password: "",
				Meta:     "test meta",
			},
			wantErr: true,
			errMsg:  "login is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.creds.Bind(&http.Request{})
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCredentialsRequest.Bind() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("UpdateCredentialsRequest.Bind() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestCredentials_Struct(t *testing.T) {
	now := time.Now()
	creds := Credentials{
		CredentialsID: 1,
		Login:         "testuser",
		Password:      "testpass",
		UserID:        1,
		CreatedAt:     now,
		Meta:          "test meta",
	}

	if creds.CredentialsID != 1 {
		t.Errorf("Expected CredentialsID 1, got %d", creds.CredentialsID)
	}
	if creds.Login != "testuser" {
		t.Errorf("Expected Login testuser, got %s", creds.Login)
	}
	if creds.Password != "testpass" {
		t.Errorf("Expected Password testpass, got %s", creds.Password)
	}
	if creds.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", creds.UserID)
	}
	if !creds.CreatedAt.Equal(now) {
		t.Errorf("Expected CreatedAt %v, got %v", now, creds.CreatedAt)
	}
	if creds.Meta != "test meta" {
		t.Errorf("Expected Meta test meta, got %s", creds.Meta)
	}
}

func TestUpdateCredentialsRequest_Struct(t *testing.T) {
	updateReq := UpdateCredentialsRequest{
		Login:    "testuser",
		Password: "testpass",
		Meta:     "test meta",
	}

	if updateReq.Login != "testuser" {
		t.Errorf("Expected Login testuser, got %s", updateReq.Login)
	}
	if updateReq.Password != "testpass" {
		t.Errorf("Expected Password testpass, got %s", updateReq.Password)
	}
	if updateReq.Meta != "test meta" {
		t.Errorf("Expected Meta test meta, got %s", updateReq.Meta)
	}
}

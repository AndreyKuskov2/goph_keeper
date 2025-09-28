package models

import (
	"net/http"
	"testing"
	"time"
)

func TestUser_Bind(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid user",
			user: User{
				Login:    "testuser",
				Password: "testpass",
			},
			wantErr: false,
		},
		{
			name: "empty login",
			user: User{
				Login:    "",
				Password: "testpass",
			},
			wantErr: true,
			errMsg:  "login is required",
		},
		{
			name: "empty password",
			user: User{
				Login:    "testuser",
				Password: "",
			},
			wantErr: true,
			errMsg:  "password is required",
		},
		{
			name: "both empty",
			user: User{
				Login:    "",
				Password: "",
			},
			wantErr: true,
			errMsg:  "login is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Bind(&http.Request{})
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Bind() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("User.Bind() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestUserLoginRequest_Bind(t *testing.T) {
	tests := []struct {
		name    string
		user    UserLoginRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid login request",
			user: UserLoginRequest{
				Login:    "testuser",
				Password: "testpass",
			},
			wantErr: false,
		},
		{
			name: "empty login",
			user: UserLoginRequest{
				Login:    "",
				Password: "testpass",
			},
			wantErr: true,
			errMsg:  "login is required",
		},
		{
			name: "empty password",
			user: UserLoginRequest{
				Login:    "testuser",
				Password: "",
			},
			wantErr: true,
			errMsg:  "password is required",
		},
		{
			name: "both empty",
			user: UserLoginRequest{
				Login:    "",
				Password: "",
			},
			wantErr: true,
			errMsg:  "login is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Bind(&http.Request{})
			if (err != nil) != tt.wantErr {
				t.Errorf("UserLoginRequest.Bind() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("UserLoginRequest.Bind() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestUser_Struct(t *testing.T) {
	now := time.Now()
	user := User{
		UserID:    1,
		Login:     "testuser",
		Password:  "testpass",
		CreatedAt: now,
	}

	if user.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", user.UserID)
	}
	if user.Login != "testuser" {
		t.Errorf("Expected Login testuser, got %s", user.Login)
	}
	if user.Password != "testpass" {
		t.Errorf("Expected Password testpass, got %s", user.Password)
	}
	if !user.CreatedAt.Equal(now) {
		t.Errorf("Expected CreatedAt %v, got %v", now, user.CreatedAt)
	}
}

func TestUserLoginRequest_Struct(t *testing.T) {
	loginReq := UserLoginRequest{
		Login:    "testuser",
		Password: "testpass",
	}

	if loginReq.Login != "testuser" {
		t.Errorf("Expected Login testuser, got %s", loginReq.Login)
	}
	if loginReq.Password != "testpass" {
		t.Errorf("Expected Password testpass, got %s", loginReq.Password)
	}
}

func TestUserLoginResponse_Struct(t *testing.T) {
	loginResp := UserLoginResponse{
		Token: "jwt-token-123",
	}

	if loginResp.Token != "jwt-token-123" {
		t.Errorf("Expected Token jwt-token-123, got %s", loginResp.Token)
	}
}

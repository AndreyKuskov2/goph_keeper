// Package models defines the data structures used throughout the GophKeeper application.
// It includes models for users, credentials, text data, bank cards, and binary data,
// along with their validation methods and request/response structures.
package models

import (
	"fmt"
	"net/http"
	"time"
)

// User represents a user account in the GophKeeper system.
type User struct {
	UserID    int       `json:"user_id"`    // Unique identifier for the user
	Login     string    `json:"login"`      // User's login name (must be unique)
	Password  string    `json:"password"`   // User's password (should be hashed)
	CreatedAt time.Time `json:"created_at"` // Timestamp when the user was created
}

// Bind validates the User struct fields from an HTTP request.
func (user *User) Bind(r *http.Request) error {
	if user.Login == "" {
		return fmt.Errorf("login is required")
	}
	if user.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

// UserLoginRequest represents the request structure for user authentication.
type UserLoginRequest struct {
	Login    string `json:"login"`    // User's login name
	Password string `json:"password"` // User's password
}

// Bind validates the UserLoginRequest struct fields from an HTTP request.
func (user *UserLoginRequest) Bind(r *http.Request) error {
	if user.Login == "" {
		return fmt.Errorf("login is required")
	}
	if user.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

// UserLoginResponse represents the response structure for successful user authentication.
type UserLoginResponse struct {
	Token string `json:"token"` // JWT token for authenticated requests
}

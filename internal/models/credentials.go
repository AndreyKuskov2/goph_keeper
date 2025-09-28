package models

import (
	"fmt"
	"net/http"
	"time"
)

// Credentials represents stored login credentials for a user.
type Credentials struct {
	CredentialsID int       `json:"credentials_id"` // Unique identifier for the credentials
	Login         string    `json:"login"`          // Login name for the stored credentials
	Password      string    `json:"password"`       // Password for the stored credentials
	UserID        int       `json:"user_id"`        // ID of the user who owns these credentials
	CreatedAt     time.Time `json:"created_at"`     // Timestamp when credentials were created
	Meta          string    `json:"meta"`           // Additional metadata or description
}

// Bind validates the Credentials struct fields from an HTTP request.
func (c *Credentials) Bind(r *http.Request) error {
	if c.Login == "" {
		return fmt.Errorf("login is required")
	}
	if c.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

// UpdateCredentialsRequest represents the request structure for updating existing credentials.
type UpdateCredentialsRequest struct {
	Login    string `json:"login"`    // Updated login name
	Password string `json:"password"` // Updated password
	Meta     string `json:"meta"`     // Updated metadata or description
}

// Bind validates the UpdateCredentialsRequest struct fields from an HTTP request.
func (c *UpdateCredentialsRequest) Bind(r *http.Request) error {
	if c.Login == "" {
		return fmt.Errorf("login is required")
	}
	if c.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

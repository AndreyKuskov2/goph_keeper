package models

import (
	"fmt"
	"net/http"
	"time"
)

type Credentials struct {
	CredentialsID int       `json:"credentials_id"`
	Login         string    `json:"login"`
	Password      string    `json:"password"`
	UserID        int       `json:"user_id"`
	CreatedAt     time.Time `json:"created_at"`
	Meta          string    `json:"meta"`
}

func (c *Credentials) Bind(r *http.Request) error {
	if c.Login == "" {
		return fmt.Errorf("login is required")
	}
	if c.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

type UpdateCredentialsRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Meta     string `json:"meta"`
}

func (c *UpdateCredentialsRequest) Bind(r *http.Request) error {
	if c.Login == "" {
		return fmt.Errorf("login is required")
	}
	if c.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

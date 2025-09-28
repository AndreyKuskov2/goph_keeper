package models

import (
	"fmt"
	"net/http"
	"time"
)

// TextData represents stored text information for a user.
type TextData struct {
	TextDataID int       `json:"text_data_id"` // Unique identifier for the text data
	Text       string    `json:"text"`         // The actual text content to store
	UserID     int       `json:"user_id"`      // ID of the user who owns this text data
	CreatedAt  time.Time `json:"created_at"`   // Timestamp when text data was created
	Meta       string    `json:"meta"`         // Additional metadata or description
}

// Bind validates the TextData struct fields from an HTTP request.
func (td *TextData) Bind(r *http.Request) error {
	if td.Text == "" {
		return fmt.Errorf("login is required")
	}
	return nil
}

// UpdateTextDataRequest represents the request structure for updating existing text data.
type UpdateTextDataRequest struct {
	Text string `json:"text"` // Updated text content
	Meta string `json:"meta"` // Updated metadata or description
}

// Bind validates the UpdateTextDataRequest struct fields from an HTTP request.
func (td *UpdateTextDataRequest) Bind(r *http.Request) error {
	if td.Text == "" {
		return fmt.Errorf("login is required")
	}
	return nil
}

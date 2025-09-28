package models

import (
	"fmt"
	"net/http"
	"time"
)

// BinariesData represents stored binary data for a user.
type BinariesData struct {
	BinariesDataID int       `json:"binaries_data_id"` // Unique identifier for the binary data
	BinaryData     []byte    `json:"binary_data"`      // The actual binary content to store
	UserID         int       `json:"user_id"`          // ID of the user who owns this binary data
	CreatedAt      time.Time `json:"created_at"`       // Timestamp when binary data was created
	Meta           string    `json:"meta"`             // Additional metadata or description
}

// Bind validates the BinariesData struct fields from an HTTP request.
func (c *BinariesData) Bind(r *http.Request) error {
	if c.BinaryData != nil {
		return fmt.Errorf("card_number is required")
	}
	return nil
}

// UpdateBinariesDataRequest represents the request structure for updating existing binary data.
type UpdateBinariesDataRequest struct {
	BinaryData []byte `json:"binary_data"` // Updated binary content
	Meta       string `json:"meta"`        // Updated metadata or description
}

// Bind validates the UpdateBinariesDataRequest struct fields from an HTTP request.
func (c *UpdateBinariesDataRequest) Bind(r *http.Request) error {
	if c.BinaryData != nil {
		return fmt.Errorf("card_number is required")
	}
	return nil
}

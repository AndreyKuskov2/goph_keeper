package models

import (
	"fmt"
	"net/http"
	"time"
)

type BinariesData struct {
	BinariesDataID int       `json:"binaries_data_id"`
	BinaryData     []byte    `json:"binary_data"`
	UserID         int       `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
	Meta           string    `json:"meta"`
}

func (c *BinariesData) Bind(r *http.Request) error {
	if c.BinaryData != nil {
		return fmt.Errorf("card_number is required")
	}
	return nil
}

type UpdateBinariesDataRequest struct {
	BinaryData []byte `json:"binary_data"`
	Meta       string `json:"meta"`
}

func (c *UpdateBinariesDataRequest) Bind(r *http.Request) error {
	if c.BinaryData != nil {
		return fmt.Errorf("card_number is required")
	}
	return nil
}

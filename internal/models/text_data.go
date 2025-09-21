package models

import (
	"fmt"
	"net/http"
	"time"
)

type TextData struct {
	TextDataID int       `json:"text_data_id"`
	Text       string    `json:"text"`
	UserID     int       `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (td *TextData) Bind(r *http.Request) error {
	if td.Text == "" {
		return fmt.Errorf("login is required")
	}
	return nil
}

type UpdateTextDataRequest struct {
	Text string `json:"text"`
}

func (td *UpdateTextDataRequest) Bind(r *http.Request) error {
	if td.Text == "" {
		return fmt.Errorf("login is required")
	}
	return nil
}

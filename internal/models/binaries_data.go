package models

import "time"

type BinariesData struct {
	BinariesDataID int       `json:"binaries_data_id"`
	BinaryData     []byte    `json:"binary_data"`
	UserID         int       `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
}

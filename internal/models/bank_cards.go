package models

import (
	"fmt"
	"net/http"
	"time"
)

type BankCards struct {
	BankCardsID    int       `json:"bank_cards_id"`
	CardNumber     string    `json:"card_number"`
	Holder         string    `json:"holder"`
	CVC            string    `json:"cvc"`
	ExpirationDate time.Time `json:"expiration_date"`
	UserID         int       `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
	Meta           string    `json:"meta"`
}

func (c *BankCards) Bind(r *http.Request) error {
	if c.CardNumber == "" {
		return fmt.Errorf("card_number is required")
	}
	if c.Holder == "" {
		return fmt.Errorf("holder is required")
	}
	if c.CVC == "" {
		return fmt.Errorf("cvc is required")
	}
	return nil
}

type UpdateBankCardsRequest struct {
	CardNumber     string    `json:"card_number"`
	Holder         string    `json:"holder"`
	CVC            string    `json:"cvc"`
	ExpirationDate time.Time `json:"expiration_date"`
	Meta           string    `json:"meta"`
}

func (c *UpdateBankCardsRequest) Bind(r *http.Request) error {
	if c.CardNumber == "" {
		return fmt.Errorf("card_number is required")
	}
	if c.Holder == "" {
		return fmt.Errorf("holder is required")
	}
	if c.CVC == "" {
		return fmt.Errorf("cvc is required")
	}
	return nil
}

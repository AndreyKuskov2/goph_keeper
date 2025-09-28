package models

import (
	"fmt"
	"goph_keeper/pkg/validator"
	"net/http"
	"time"
)

// BankCards represents stored bank card information for a user.
type BankCards struct {
	BankCardsID    int       `json:"bank_cards_id"`   // Unique identifier for the bank card
	CardNumber     string    `json:"card_number"`     // Bank card number
	Holder         string    `json:"holder"`          // Cardholder name
	CVC            string    `json:"cvc"`             // Card verification code
	ExpirationDate time.Time `json:"expiration_date"` // Card expiration date
	UserID         int       `json:"user_id"`         // ID of the user who owns this card
	CreatedAt      time.Time `json:"created_at"`      // Timestamp when card was stored
	Meta           string    `json:"meta"`            // Additional metadata or description
}

// Bind validates the BankCards struct fields from an HTTP request.
func (c *BankCards) Bind(r *http.Request) error {
	if c.CardNumber == "" {
		return fmt.Errorf("card_number is required")
	}
	if !validator.LuhnAlgorith(c.CardNumber) {
		return fmt.Errorf("card_number is invalid")
	}
	if c.Holder == "" {
		return fmt.Errorf("holder is required")
	}
	if c.CVC == "" {
		return fmt.Errorf("cvc is required")
	}
	return nil
}

// UpdateBankCardsRequest represents the request structure for updating existing bank card information.
type UpdateBankCardsRequest struct {
	CardNumber     string    `json:"card_number"`     // Updated bank card number
	Holder         string    `json:"holder"`          // Updated cardholder name
	CVC            string    `json:"cvc"`             // Updated card verification code
	ExpirationDate time.Time `json:"expiration_date"` // Updated card expiration date
	Meta           string    `json:"meta"`            // Updated metadata or description
}

// Bind validates the UpdateBankCardsRequest struct fields from an HTTP request.
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

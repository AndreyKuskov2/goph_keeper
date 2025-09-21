package models

import "time"

type BankCards struct {
	BankCardsID    int       `json:"bank_cards_id"`
	CardNumber     string    `json:"card_number"`
	Holder         string    `json:"holder"`
	CVC            string    `json:"cvc"`
	ExpirationDate time.Time `json:"expiration_date"`
	UserID         int       `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
}

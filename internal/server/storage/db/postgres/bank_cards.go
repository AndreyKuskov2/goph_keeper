package postgres

import (
	"context"
	"goph_keeper/internal/models"

	"github.com/jackc/pgx/v5"
)

// CreateBankCards creates a new bank card.
func (pg *Postgres) CreateBankCards(ctx context.Context, bankCards *models.BankCards) (int, error) {
	var bankCardID int
	if err := pg.DB.QueryRow(ctx, createBankCards, bankCards.CardNumber, bankCards.Holder, bankCards.ExpirationDate, bankCards.UserID, bankCards.Meta).Scan(&bankCardID); err != nil {
		return 0, err
	}
	return bankCardID, nil
}

// GetBankCardsByUserID gets all bank cards by user ID.
func (pg *Postgres) GetBankCardsByUserID(ctx context.Context, userID int) ([]models.BankCards, error) {
	rows, err := pg.DB.Query(ctx, getBankCardsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bankCards, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.BankCards])
	if err != nil {
		return nil, err
	}

	return bankCards, nil
}

// GetBankCardsByIDAndUserID gets a bank card by ID and user ID.
func (pg *Postgres) GetBankCardsByIDAndUserID(ctx context.Context, bankCardID, userID int) ([]models.BankCards, error) {
	rows, err := pg.DB.Query(ctx, getBankCardsByIDAndUserID, bankCardID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bankCards, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.BankCards])
	if err != nil {
		return nil, err
	}

	return bankCards, nil
}

// UpdateBankCardsByIDAndUserID updates a bank card by ID and user ID.
func (pg *Postgres) UpdateBankCardsByIDAndUserID(ctx context.Context, bankCard models.UpdateBankCardsRequest, bankCardID, userID int) error {
	if _, err := pg.DB.Exec(ctx, updateBankCardsByIDAndUserID, bankCard.CardNumber, bankCard.Holder, bankCard.CVC, bankCard.ExpirationDate, bankCard.Meta, bankCardID, userID); err != nil {
		return err
	}
	return nil
}

// DeleteBankCardsByIDAndUserID deletes a bank card by ID and user ID.
func (pg *Postgres) DeleteBankCardsByIDAndUserID(ctx context.Context, bankCardID, userID int) error {
	if _, err := pg.DB.Exec(ctx, deleteBankCardsByIDAndUserID, bankCardID, userID); err != nil {
		return err
	}
	return nil
}

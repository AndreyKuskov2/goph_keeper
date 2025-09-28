package postgres

import (
	"context"
	"goph_keeper/internal/models"

	"github.com/jackc/pgx/v5"
)

// CreateTextData creates a new text data.
func (pg *Postgres) CreateTextData(ctx context.Context, textData *models.TextData) (int, error) {
	var textDataID int
	if err := pg.DB.QueryRow(ctx, createTextData, textData.Text, textData.UserID, textData.Meta).Scan(&textDataID); err != nil {
		return 0, err
	}
	return textDataID, nil
}

// GetTextDataByUserID gets all text data by user ID.
func (pg *Postgres) GetTextDataByUserID(ctx context.Context, userID int) ([]models.TextData, error) {
	rows, err := pg.DB.Query(ctx, getTextDataByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	textData, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.TextData])
	if err != nil {
		return nil, err
	}

	return textData, nil
}

// GetTextDataByIDAndUserID gets a text data by ID and user ID.
func (pg *Postgres) GetTextDataByIDAndUserID(ctx context.Context, textDataID, userID int) ([]models.TextData, error) {
	rows, err := pg.DB.Query(ctx, getTextDataByIDAndUserID, textDataID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	textData, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.TextData])
	if err != nil {
		return nil, err
	}

	return textData, nil
}

// UpdateTextDataByIDAndUserID updates a text data by ID and user ID.
func (pg *Postgres) UpdateTextDataByIDAndUserID(ctx context.Context, textData models.UpdateTextDataRequest, textDataID, userID int) error {
	if _, err := pg.DB.Exec(ctx, updateTextDataByIDAndUserID, textData.Text, textData.Meta, textDataID, userID); err != nil {
		return err
	}
	return nil
}

// DeleteTextDataByIDAndUserID deletes a text data by ID and user ID.
func (pg *Postgres) DeleteTextDataByIDAndUserID(ctx context.Context, textDataID, userID int) error {
	if _, err := pg.DB.Exec(ctx, deleteTextDataByIDAndUserID, textDataID, userID); err != nil {
		return err
	}
	return nil
}

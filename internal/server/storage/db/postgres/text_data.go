package postgres

import (
	"context"
	"goph_keeper/internal/models"

	"github.com/jackc/pgx/v5"
)

func (pg *Postgres) CreateTextData(ctx context.Context, textData *models.TextData) (int, error) {
	var textDataID int
	if err := pg.DB.QueryRow(ctx, createTextData, textData.Text, textData.UserID).Scan(&textDataID); err != nil {
		return 0, err
	}
	return textDataID, nil
}

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

func (pg *Postgres) UpdateTextDataByIDAndUserID(ctx context.Context, textData models.UpdateTextDataRequest, textDataID, userID int) error {
	if _, err := pg.DB.Exec(ctx, updateTextDataByIDAndUserID, textData.Text, textDataID, userID); err != nil {
		return err
	}
	return nil
}

func (pg *Postgres) DeleteTextDataByIDAndUserID(ctx context.Context, textDataID, userID int) error {
	if _, err := pg.DB.Exec(ctx, deleteTextDataByIDAndUserID, textDataID, userID); err != nil {
		return err
	}
	return nil
}

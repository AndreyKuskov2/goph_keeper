package postgres

import (
	"context"
	"goph_keeper/internal/models"

	"github.com/jackc/pgx/v5"
)

func (pg *Postgres) CreateBinariesData(ctx context.Context, binariesData *models.BinariesData) (int, error) {
	var binariesDataID int
	if err := pg.DB.QueryRow(ctx, createBinariesData, binariesData.BinaryData, binariesData.UserID, binariesData.Meta).Scan(&binariesDataID); err != nil {
		return 0, err
	}
	return binariesDataID, nil
}

func (pg *Postgres) GetBinariesDataByUserID(ctx context.Context, userID int) ([]models.BinariesData, error) {
	rows, err := pg.DB.Query(ctx, getBinariesDataByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	binariesData, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.BinariesData])
	if err != nil {
		return nil, err
	}

	return binariesData, nil
}

func (pg *Postgres) GetBinariesDataByIDAndUserID(ctx context.Context, binariesDataID, userID int) ([]models.BinariesData, error) {
	rows, err := pg.DB.Query(ctx, getBinariesDataByIDAndUserID, binariesDataID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	binariesData, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.BinariesData])
	if err != nil {
		return nil, err
	}

	return binariesData, nil
}

func (pg *Postgres) UpdateBinariesDataByIDAndUserID(ctx context.Context, binaryData models.UpdateBinariesDataRequest, binariesDataID, userID int) error {
	if _, err := pg.DB.Exec(ctx, updateBinariesDataByIDAndUserID, binaryData.BinaryData, binaryData.Meta, binariesDataID, userID); err != nil {
		return err
	}
	return nil
}

func (pg *Postgres) DeleteBinariesDataByIDAndUserID(ctx context.Context, binariesDataID, userID int) error {
	if _, err := pg.DB.Exec(ctx, deleteBinariesDataByIDAndUserID, binariesDataID, userID); err != nil {
		return err
	}
	return nil
}

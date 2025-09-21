package postgres

import (
	"context"
	"goph_keeper/internal/models"

	"github.com/jackc/pgx/v5"
)

func (pg *Postgres) CreateCredentials(ctx context.Context, credentials *models.Credentials) (int, error) {
	var credentialsID int
	if err := pg.DB.QueryRow(ctx, createCredentials, credentials.Login, credentials.Password, credentials.UserID).Scan(&credentialsID); err != nil {
		return 0, err
	}
	return credentialsID, nil
}

func (pg *Postgres) GetCredentialsByUserID(ctx context.Context, userID int) ([]models.Credentials, error) {
	rows, err := pg.DB.Query(ctx, getCredentialsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	credentials, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Credentials])
	if err != nil {
		return nil, err
	}

	return credentials, nil
}

func (pg *Postgres) GetCredentialsByIDAndUserID(ctx context.Context, credentialsID, userID int) ([]models.Credentials, error) {
	rows, err := pg.DB.Query(ctx, getCredentialsByIDAndUserID, credentialsID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	credentials, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Credentials])
	if err != nil {
		return nil, err
	}

	return credentials, nil
}

func (pg *Postgres) UpdateCredentialsByIDAndUserID(ctx context.Context, credentials models.UpdateCredentialsRequest, credentialsID, userID int) error {
	if _, err := pg.DB.Exec(ctx, updateCredentialsByIDAndUserID, credentials.Login, credentials.Password, credentialsID, userID); err != nil {
		return err
	}
	return nil
}

func (pg *Postgres) DeleteCredentialsByIDAndUserID(ctx context.Context, credentialsID, userID int) error {
	if _, err := pg.DB.Exec(ctx, deleteCredentialsByIDAndUserID, credentialsID, userID); err != nil {
		return err
	}
	return nil
}

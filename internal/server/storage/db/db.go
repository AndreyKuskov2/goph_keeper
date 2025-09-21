package db

import (
	"context"
	"fmt"
	"goph_keeper/internal/models"
	"goph_keeper/internal/server/config"
	"goph_keeper/internal/server/storage/db/postgres"
)

type DBStorager interface {
	CreateUser(ctx context.Context, user models.User) (int, error)
	GetUserByLogin(ctx context.Context, user models.UserLoginRequest) (int, error)

	CreateCredentials(ctx context.Context, credentials *models.Credentials) (int, error)
	GetCredentialsByUserID(ctx context.Context, userID int) ([]models.Credentials, error)
	GetCredentialsByIDAndUserID(ctx context.Context, credentialsID, userID int) ([]models.Credentials, error)
	UpdateCredentialsByIDAndUserID(ctx context.Context, credentials models.UpdateCredentialsRequest, credentialsID, userID int) error
	DeleteCredentialsByIDAndUserID(ctx context.Context, credentialsID, userID int) error

	CreateTextData(ctx context.Context, textData *models.TextData) (int, error)
	GetTextDataByUserID(ctx context.Context, userID int) ([]models.TextData, error)
	GetTextDataByIDAndUserID(ctx context.Context, textDataID, userID int) ([]models.TextData, error)
	UpdateTextDataByIDAndUserID(ctx context.Context, textData models.UpdateTextDataRequest, textDataID, userID int) error
	DeleteTextDataByIDAndUserID(ctx context.Context, textDataID, userID int) error

	Ping(ctx context.Context) error
	Close()
}

func NewDBStorage(ctx context.Context, cfg *config.Database) (DBStorager, error) {
	if cfg.Type == "postgres" {
		return postgres.NewPostgres(ctx, cfg)
	}
	return nil, fmt.Errorf("available database drivers: postgres")
}

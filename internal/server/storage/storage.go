package storage

import (
	"context"
	"fmt"
	"goph_keeper/internal/models"
	"goph_keeper/internal/server/config"
	"goph_keeper/internal/server/storage/db"
)

type Storager interface {
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

	CreateBankCards(ctx context.Context, bankCards *models.BankCards) (int, error)
	GetBankCardsByUserID(ctx context.Context, userID int) ([]models.BankCards, error)
	GetBankCardsByIDAndUserID(ctx context.Context, bankCardID, userID int) ([]models.BankCards, error)
	UpdateBankCardsByIDAndUserID(ctx context.Context, bankCard models.UpdateBankCardsRequest, bankCardID, userID int) error
	DeleteBankCardsByIDAndUserID(ctx context.Context, bankCardID, userID int) error

	CreateBinariesData(ctx context.Context, binariesData *models.BinariesData) (int, error)
	GetBinariesDataByUserID(ctx context.Context, userID int) ([]models.BinariesData, error)
	GetBinariesDataByIDAndUserID(ctx context.Context, binariesDataID, userID int) ([]models.BinariesData, error)
	UpdateBinariesDataByIDAndUserID(ctx context.Context, binaryData models.UpdateBinariesDataRequest, binariesDataID, userID int) error
	DeleteBinariesDataByIDAndUserID(ctx context.Context, binariesDataID, userID int) error

	Ping(ctx context.Context) error
	Close()
}

func NewStorager(ctx context.Context, cfg *config.Config) (Storager, error) {
	if cfg.Database != nil {
		return db.NewDBStorage(ctx, cfg.Database)
	}
	return nil, fmt.Errorf("available storages: database")
}

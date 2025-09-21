package service

import (
	"context"
	"goph_keeper/internal/models"
)

type GophKeeperCredentialsStorager interface {
	CreateCredentials(ctx context.Context, credentials *models.Credentials) (int, error)
	GetCredentialsByUserID(ctx context.Context, userID int) ([]models.Credentials, error)
	GetCredentialsByIDAndUserID(ctx context.Context, credentialsID, userID int) ([]models.Credentials, error)
	UpdateCredentialsByIDAndUserID(ctx context.Context, credentials models.UpdateCredentialsRequest, credentialsID, userID int) error
	DeleteCredentialsByIDAndUserID(ctx context.Context, credentialsID, userID int) error
}

type GophKeeperCredentialsService struct {
	storage GophKeeperCredentialsStorager
}

func NewGophKeeperCredentialsService(storage GophKeeperCredentialsStorager) *GophKeeperCredentialsService {
	return &GophKeeperCredentialsService{
		storage: storage,
	}
}

func (s *GophKeeperCredentialsService) CreateCredentials(ctx context.Context, credentials *models.Credentials) (int, error) {
	return s.storage.CreateCredentials(ctx, credentials)
}
func (s *GophKeeperCredentialsService) GetCredentialsByUserID(ctx context.Context, userID int) ([]models.Credentials, error) {
	return s.storage.GetCredentialsByUserID(ctx, userID)
}
func (s *GophKeeperCredentialsService) GetCredentialsByIDAndUserID(ctx context.Context, credentialsID, userID int) ([]models.Credentials, error) {
	return s.storage.GetCredentialsByIDAndUserID(ctx, credentialsID, userID)
}
func (s *GophKeeperCredentialsService) UpdateCredentialsByIDAndUserID(ctx context.Context, credentials models.UpdateCredentialsRequest, credentialsID, userID int) error {
	return s.storage.UpdateCredentialsByIDAndUserID(ctx, credentials, credentialsID, userID)
}
func (s *GophKeeperCredentialsService) DeleteCredentialsByIDAndUserID(ctx context.Context, credentialsID, userID int) error {
	return s.storage.DeleteCredentialsByIDAndUserID(ctx, credentialsID, userID)
}

package service

import (
	"context"
	"goph_keeper/internal/models"
)

// GophKeeperCredentialsStorager is the interface that wraps the basic methods for working with credentials.
type GophKeeperCredentialsStorager interface {
	CreateCredentials(ctx context.Context, credentials *models.Credentials) (int, error)
	GetCredentialsByUserID(ctx context.Context, userID int) ([]models.Credentials, error)
	GetCredentialsByIDAndUserID(ctx context.Context, credentialsID, userID int) ([]models.Credentials, error)
	UpdateCredentialsByIDAndUserID(ctx context.Context, credentials models.UpdateCredentialsRequest, credentialsID, userID int) error
	DeleteCredentialsByIDAndUserID(ctx context.Context, credentialsID, userID int) error
}

// GophKeeperCredentialsService is the service that wraps the basic methods for working with credentials.
type GophKeeperCredentialsService struct {
	storage GophKeeperCredentialsStorager
}

// NewGophKeeperCredentialsService creates a new GophKeeperCredentialsService.
func NewGophKeeperCredentialsService(storage GophKeeperCredentialsStorager) *GophKeeperCredentialsService {
	return &GophKeeperCredentialsService{
		storage: storage,
	}
}

// CreateCredentials creates a new credentials.
func (s *GophKeeperCredentialsService) CreateCredentials(ctx context.Context, credentials *models.Credentials) (int, error) {
	return s.storage.CreateCredentials(ctx, credentials)
}

// GetCredentialsByUserID gets all credentials by user ID.
func (s *GophKeeperCredentialsService) GetCredentialsByUserID(ctx context.Context, userID int) ([]models.Credentials, error) {
	return s.storage.GetCredentialsByUserID(ctx, userID)
}

// GetCredentialsByIDAndUserID gets a credentials by ID and user ID.
func (s *GophKeeperCredentialsService) GetCredentialsByIDAndUserID(ctx context.Context, credentialsID, userID int) ([]models.Credentials, error) {
	return s.storage.GetCredentialsByIDAndUserID(ctx, credentialsID, userID)
}

// UpdateCredentialsByIDAndUserID updates a credentials by ID and user ID.
func (s *GophKeeperCredentialsService) UpdateCredentialsByIDAndUserID(ctx context.Context, credentials models.UpdateCredentialsRequest, credentialsID, userID int) error {
	return s.storage.UpdateCredentialsByIDAndUserID(ctx, credentials, credentialsID, userID)
}

// DeleteCredentialsByIDAndUserID deletes a credentials by ID and user ID.
func (s *GophKeeperCredentialsService) DeleteCredentialsByIDAndUserID(ctx context.Context, credentialsID, userID int) error {
	return s.storage.DeleteCredentialsByIDAndUserID(ctx, credentialsID, userID)
}

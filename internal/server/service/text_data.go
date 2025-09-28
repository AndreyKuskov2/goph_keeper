package service

import (
	"context"
	"goph_keeper/internal/models"
)

// GophKeeperTextDataStorager is the interface that wraps the basic methods for working with text data.
type GophKeeperTextDataStorager interface {
	CreateTextData(ctx context.Context, textData *models.TextData) (int, error)
	GetTextDataByUserID(ctx context.Context, userID int) ([]models.TextData, error)
	GetTextDataByIDAndUserID(ctx context.Context, textDataID, userID int) ([]models.TextData, error)
	UpdateTextDataByIDAndUserID(ctx context.Context, textData models.UpdateTextDataRequest, textDataID, userID int) error
	DeleteTextDataByIDAndUserID(ctx context.Context, textDataID, userID int) error
}

// GophKeeperTextDataService is the service that wraps the basic methods for working with text data.
type GophKeeperTextDataService struct {
	storage GophKeeperTextDataStorager
}

// NewGophKeeperTextDataService creates a new GophKeeperTextDataService.
func NewGophKeeperTextDataService(storage GophKeeperTextDataStorager) *GophKeeperTextDataService {
	return &GophKeeperTextDataService{
		storage: storage,
	}
}

// CreateTextData creates a new text data.
func (s *GophKeeperTextDataService) CreateTextData(ctx context.Context, textData *models.TextData) (int, error) {
	return s.storage.CreateTextData(ctx, textData)
}

// GetTextDataByUserID gets all text data by user ID.
func (s *GophKeeperTextDataService) GetTextDataByUserID(ctx context.Context, userID int) ([]models.TextData, error) {
	return s.storage.GetTextDataByUserID(ctx, userID)
}

// GetTextDataByIDAndUserID gets a text data by ID and user ID.
func (s *GophKeeperTextDataService) GetTextDataByIDAndUserID(ctx context.Context, credentialsID, userID int) ([]models.TextData, error) {
	return s.storage.GetTextDataByIDAndUserID(ctx, credentialsID, userID)
}

// UpdateTextDataByIDAndUserID updates a text data by ID and user ID.
func (s *GophKeeperTextDataService) UpdateTextDataByIDAndUserID(ctx context.Context, credentials models.UpdateTextDataRequest, credentialsID, userID int) error {
	return s.storage.UpdateTextDataByIDAndUserID(ctx, credentials, credentialsID, userID)
}

// DeleteTextDataByIDAndUserID deletes a text data by ID and user ID.
func (s *GophKeeperTextDataService) DeleteTextDataByIDAndUserID(ctx context.Context, credentialsID, userID int) error {
	return s.storage.DeleteTextDataByIDAndUserID(ctx, credentialsID, userID)
}

package service

import (
	"context"
	"goph_keeper/internal/models"
)

type GophKeeperTextDataStorager interface {
	CreateTextData(ctx context.Context, textData *models.TextData) (int, error)
	GetTextDataByUserID(ctx context.Context, userID int) ([]models.TextData, error)
	GetTextDataByIDAndUserID(ctx context.Context, textDataID, userID int) ([]models.TextData, error)
	UpdateTextDataByIDAndUserID(ctx context.Context, textData models.UpdateTextDataRequest, textDataID, userID int) error
	DeleteTextDataByIDAndUserID(ctx context.Context, textDataID, userID int) error
}

type GophKeeperTextDataService struct {
	storage GophKeeperTextDataStorager
}

func NewGophKeeperTextDataService(storage GophKeeperTextDataStorager) *GophKeeperTextDataService {
	return &GophKeeperTextDataService{
		storage: storage,
	}
}

func (s *GophKeeperTextDataService) CreateTextData(ctx context.Context, textData *models.TextData) (int, error) {
	return s.storage.CreateTextData(ctx, textData)
}
func (s *GophKeeperTextDataService) GetTextDataByUserID(ctx context.Context, userID int) ([]models.TextData, error) {
	return s.storage.GetTextDataByUserID(ctx, userID)
}
func (s *GophKeeperTextDataService) GetTextDataByIDAndUserID(ctx context.Context, credentialsID, userID int) ([]models.TextData, error) {
	return s.storage.GetTextDataByIDAndUserID(ctx, credentialsID, userID)
}
func (s *GophKeeperTextDataService) UpdateTextDataByIDAndUserID(ctx context.Context, credentials models.UpdateTextDataRequest, credentialsID, userID int) error {
	return s.storage.UpdateTextDataByIDAndUserID(ctx, credentials, credentialsID, userID)
}
func (s *GophKeeperTextDataService) DeleteTextDataByIDAndUserID(ctx context.Context, credentialsID, userID int) error {
	return s.storage.DeleteTextDataByIDAndUserID(ctx, credentialsID, userID)
}

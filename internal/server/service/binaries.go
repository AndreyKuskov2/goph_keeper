package service

import (
	"context"
	"goph_keeper/internal/models"
)

type GophKeeperBinariesDataStorager interface {
	CreateBinariesData(ctx context.Context, binariesData *models.BinariesData) (int, error)
	GetBinariesDataByUserID(ctx context.Context, userID int) ([]models.BinariesData, error)
	GetBinariesDataByIDAndUserID(ctx context.Context, binariesDataID, userID int) ([]models.BinariesData, error)
	UpdateBinariesDataByIDAndUserID(ctx context.Context, binaryData models.UpdateBinariesDataRequest, binariesDataID, userID int) error
	DeleteBinariesDataByIDAndUserID(ctx context.Context, binariesDataID, userID int) error
}

type GophKeeperBinariesDataService struct {
	storage GophKeeperBinariesDataStorager
}

func NewGophKeeperBinariesDataService(storage GophKeeperBinariesDataStorager) *GophKeeperBinariesDataService {
	return &GophKeeperBinariesDataService{
		storage: storage,
	}
}

func (s *GophKeeperBinariesDataService) CreateBinariesData(ctx context.Context, binariesData *models.BinariesData) (int, error) {
	return s.storage.CreateBinariesData(ctx, binariesData)
}
func (s *GophKeeperBinariesDataService) GetBinariesDataByUserID(ctx context.Context, userID int) ([]models.BinariesData, error) {
	return s.storage.GetBinariesDataByUserID(ctx, userID)
}
func (s *GophKeeperBinariesDataService) GetBinariesDataByIDAndUserID(ctx context.Context, binariesDataID, userID int) ([]models.BinariesData, error) {
	return s.storage.GetBinariesDataByIDAndUserID(ctx, binariesDataID, userID)
}
func (s *GophKeeperBinariesDataService) UpdateBinariesDataByIDAndUserID(ctx context.Context, binaryData models.UpdateBinariesDataRequest, binariesDataID, userID int) error {
	return s.storage.UpdateBinariesDataByIDAndUserID(ctx, binaryData, binariesDataID, userID)
}
func (s *GophKeeperBinariesDataService) DeleteBinariesDataByIDAndUserID(ctx context.Context, binariesDataID, userID int) error {
	return s.storage.DeleteBinariesDataByIDAndUserID(ctx, binariesDataID, userID)
}

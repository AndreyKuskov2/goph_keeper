package service

import (
	"context"
	"goph_keeper/internal/models"
)

// GophKeeperBinariesDataStorager is the interface that wraps the basic methods for working with binaries data.
type GophKeeperBinariesDataStorager interface {
	CreateBinariesData(ctx context.Context, binariesData *models.BinariesData) (int, error)
	GetBinariesDataByUserID(ctx context.Context, userID int) ([]models.BinariesData, error)
	GetBinariesDataByIDAndUserID(ctx context.Context, binariesDataID, userID int) ([]models.BinariesData, error)
	UpdateBinariesDataByIDAndUserID(ctx context.Context, binaryData models.UpdateBinariesDataRequest, binariesDataID, userID int) error
	DeleteBinariesDataByIDAndUserID(ctx context.Context, binariesDataID, userID int) error
}

// GophKeeperBinariesDataService is the service that wraps the basic methods for working with binaries data.
type GophKeeperBinariesDataService struct {
	storage GophKeeperBinariesDataStorager
}

// NewGophKeeperBinariesDataService creates a new GophKeeperBinariesDataService.
func NewGophKeeperBinariesDataService(storage GophKeeperBinariesDataStorager) *GophKeeperBinariesDataService {
	return &GophKeeperBinariesDataService{
		storage: storage,
	}
}

// CreateBinariesData creates a new binaries data.
func (s *GophKeeperBinariesDataService) CreateBinariesData(ctx context.Context, binariesData *models.BinariesData) (int, error) {
	return s.storage.CreateBinariesData(ctx, binariesData)
}

// GetBinariesDataByUserID gets all binaries data by user ID.
func (s *GophKeeperBinariesDataService) GetBinariesDataByUserID(ctx context.Context, userID int) ([]models.BinariesData, error) {
	return s.storage.GetBinariesDataByUserID(ctx, userID)
}

// GetBinariesDataByIDAndUserID gets a binaries data by ID and user ID.
func (s *GophKeeperBinariesDataService) GetBinariesDataByIDAndUserID(ctx context.Context, binariesDataID, userID int) ([]models.BinariesData, error) {
	return s.storage.GetBinariesDataByIDAndUserID(ctx, binariesDataID, userID)
}

// UpdateBinariesDataByIDAndUserID updates a binaries data by ID and user ID.
func (s *GophKeeperBinariesDataService) UpdateBinariesDataByIDAndUserID(ctx context.Context, binaryData models.UpdateBinariesDataRequest, binariesDataID, userID int) error {
	return s.storage.UpdateBinariesDataByIDAndUserID(ctx, binaryData, binariesDataID, userID)
}

// DeleteBinariesDataByIDAndUserID deletes a binaries data by ID and user ID.
func (s *GophKeeperBinariesDataService) DeleteBinariesDataByIDAndUserID(ctx context.Context, binariesDataID, userID int) error {
	return s.storage.DeleteBinariesDataByIDAndUserID(ctx, binariesDataID, userID)
}

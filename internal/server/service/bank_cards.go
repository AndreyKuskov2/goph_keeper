package service

import (
	"context"
	"goph_keeper/internal/models"
)

// GophKeeperBankCardsStorager is the interface that wraps the basic methods for working with bank cards.
type GophKeeperBankCardsStorager interface {
	CreateBankCards(ctx context.Context, bankCards *models.BankCards) (int, error)
	GetBankCardsByUserID(ctx context.Context, userID int) ([]models.BankCards, error)
	GetBankCardsByIDAndUserID(ctx context.Context, bankCardID, userID int) ([]models.BankCards, error)
	UpdateBankCardsByIDAndUserID(ctx context.Context, bankCard models.UpdateBankCardsRequest, bankCardID, userID int) error
	DeleteBankCardsByIDAndUserID(ctx context.Context, bankCardID, userID int) error
}

// GophKeeperBankCardsService is the service that wraps the basic methods for working with bank cards.
type GophKeeperBankCardsService struct {
	storage GophKeeperBankCardsStorager
}

// NewGophKeeperBankCardsService creates a new GophKeeperBankCardsService.
func NewGophKeeperBankCardsService(storage GophKeeperBankCardsStorager) *GophKeeperBankCardsService {
	return &GophKeeperBankCardsService{
		storage: storage,
	}
}

// CreateBankCards creates a new bank card.
func (s *GophKeeperBankCardsService) CreateBankCards(ctx context.Context, bankCards *models.BankCards) (int, error) {
	return s.storage.CreateBankCards(ctx, bankCards)
}

// GetBankCardsByUserID gets all bank cards by user ID.
func (s *GophKeeperBankCardsService) GetBankCardsByUserID(ctx context.Context, userID int) ([]models.BankCards, error) {
	return s.storage.GetBankCardsByUserID(ctx, userID)
}

// GetBankCardsByIDAndUserID gets a bank card by ID and user ID.
func (s *GophKeeperBankCardsService) GetBankCardsByIDAndUserID(ctx context.Context, bankCardID, userID int) ([]models.BankCards, error) {
	return s.storage.GetBankCardsByIDAndUserID(ctx, bankCardID, userID)
}

// UpdateBankCardsByIDAndUserID updates a bank card by ID and user ID.
func (s *GophKeeperBankCardsService) UpdateBankCardsByIDAndUserID(ctx context.Context, bankCard models.UpdateBankCardsRequest, bankCardID, userID int) error {
	return s.storage.UpdateBankCardsByIDAndUserID(ctx, bankCard, bankCardID, userID)
}

// DeleteBankCardsByIDAndUserID deletes a bank card by ID and user ID.
func (s *GophKeeperBankCardsService) DeleteBankCardsByIDAndUserID(ctx context.Context, bankCardID, userID int) error {
	return s.storage.DeleteBankCardsByIDAndUserID(ctx, bankCardID, userID)
}

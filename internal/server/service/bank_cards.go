package service

import (
	"context"
	"goph_keeper/internal/models"
)

type GophKeeperBankCardsStorager interface {
	CreateBankCards(ctx context.Context, bankCards *models.BankCards) (int, error)
	GetBankCardsByUserID(ctx context.Context, userID int) ([]models.BankCards, error)
	GetBankCardsByIDAndUserID(ctx context.Context, bankCardID, userID int) ([]models.BankCards, error)
	UpdateBankCardsByIDAndUserID(ctx context.Context, bankCard models.UpdateBankCardsRequest, bankCardID, userID int) error
	DeleteBankCardsByIDAndUserID(ctx context.Context, bankCardID, userID int) error
}

type GophKeeperBankCardsService struct {
	storage GophKeeperBankCardsStorager
}

func NewGophKeeperBankCardsService(storage GophKeeperBankCardsStorager) *GophKeeperBankCardsService {
	return &GophKeeperBankCardsService{
		storage: storage,
	}
}

func (s *GophKeeperBankCardsService) CreateBankCards(ctx context.Context, bankCards *models.BankCards) (int, error) {
	return s.storage.CreateBankCards(ctx, bankCards)
}
func (s *GophKeeperBankCardsService) GetBankCardsByUserID(ctx context.Context, userID int) ([]models.BankCards, error) {
	return s.storage.GetBankCardsByUserID(ctx, userID)
}
func (s *GophKeeperBankCardsService) GetBankCardsByIDAndUserID(ctx context.Context, bankCardID, userID int) ([]models.BankCards, error) {
	return s.storage.GetBankCardsByIDAndUserID(ctx, bankCardID, userID)
}
func (s *GophKeeperBankCardsService) UpdateBankCardsByIDAndUserID(ctx context.Context, bankCard models.UpdateBankCardsRequest, bankCardID, userID int) error {
	return s.storage.UpdateBankCardsByIDAndUserID(ctx, bankCard, bankCardID, userID)
}
func (s *GophKeeperBankCardsService) DeleteBankCardsByIDAndUserID(ctx context.Context, bankCardID, userID int) error {
	return s.storage.DeleteBankCardsByIDAndUserID(ctx, bankCardID, userID)
}

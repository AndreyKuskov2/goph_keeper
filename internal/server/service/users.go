package service

import (
	"context"
	"goph_keeper/internal/models"
)

type GophKeeperUserStorager interface {
	CreateUser(ctx context.Context, user models.User) (int, error)
	GetUserByLogin(ctx context.Context, user models.UserLoginRequest) (int, error)
}

type GophKeeperUserService struct {
	storage GophKeeperUserStorager
}

func NewGophermartUserService(storage GophKeeperUserStorager) *GophKeeperUserService {
	return &GophKeeperUserService{
		storage: storage,
	}
}

func (gs *GophKeeperUserService) RegisterUserService(ctx context.Context, user models.User) (int, error) {
	return gs.storage.CreateUser(ctx, user)
}

func (gs *GophKeeperUserService) GetUserService(ctx context.Context, user models.UserLoginRequest) (int, error) {
	return gs.storage.GetUserByLogin(ctx, user)
}

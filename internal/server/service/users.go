package service

import (
	"context"
	"goph_keeper/internal/models"
)

// GophKeeperUserStorager is the interface that wraps the basic methods for working with users.
type GophKeeperUserStorager interface {
	CreateUser(ctx context.Context, user models.User) (int, error)
	GetUserByLogin(ctx context.Context, user models.UserLoginRequest) (int, error)
}

// GophKeeperUserService is the service that wraps the basic methods for working with users.
type GophKeeperUserService struct {
	storage GophKeeperUserStorager
}

// NewGophermartUserService creates a new GophKeeperUserService.
func NewGophermartUserService(storage GophKeeperUserStorager) *GophKeeperUserService {
	return &GophKeeperUserService{
		storage: storage,
	}
}

// RegisterUserService creates a new user.
func (gs *GophKeeperUserService) RegisterUserService(ctx context.Context, user models.User) (int, error) {
	return gs.storage.CreateUser(ctx, user)
}

// GetUserService gets a user by login.
func (gs *GophKeeperUserService) GetUserService(ctx context.Context, user models.UserLoginRequest) (int, error) {
	return gs.storage.GetUserByLogin(ctx, user)
}

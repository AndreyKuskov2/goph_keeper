package postgres

import (
	"context"
	"fmt"
	"goph_keeper/internal/models"
	custom_errors "goph_keeper/internal/server/errors"

	"golang.org/x/crypto/bcrypt"
)

// CreateUser creates a new user.
func (pg *Postgres) CreateUser(ctx context.Context, user models.User) (int, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return 0, fmt.Errorf("cannot hashing password: %v", err)
	}

	var userID int
	if err := pg.DB.QueryRow(ctx, checkUserIsExists, user.Login).Scan(&userID); err == nil {
		return 0, custom_errors.ErrUserIsExist
	}

	if err := pg.DB.QueryRow(ctx, createNewUser, user.Login, string(passwordHash)).Scan(&userID); err != nil {
		fmt.Println(err)
		return 0, fmt.Errorf("cannot create user: %v", err)
	}
	return userID, nil
}

// GetUserByLogin gets a user by login.
func (db *Postgres) GetUserByLogin(ctx context.Context, user models.UserLoginRequest) (int, error) {
	var userID int
	var passwordHash string

	if err := db.DB.QueryRow(ctx, getUserPasswordByLogin, user.Login).Scan(&userID, &passwordHash); err != nil {
		return 0, fmt.Errorf("user not found: %v", err)
	}

	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(user.Password))
	if err != nil {
		return 0, custom_errors.ErrInvalidData
	}
	return userID, nil
}

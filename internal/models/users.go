package models

import (
	"fmt"
	"net/http"
	"time"
)

type User struct {
	UserID    int       `json:"user_id"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

func (user *User) Bind(r *http.Request) error {
	if user.Login == "" {
		return fmt.Errorf("login is required")
	}
	if user.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

type UserLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (user *UserLoginRequest) Bind(r *http.Request) error {
	if user.Login == "" {
		return fmt.Errorf("login is required")
	}
	if user.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

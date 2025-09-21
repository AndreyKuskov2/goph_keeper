package http_handlers

import (
	"context"
	"database/sql"
	"errors"
	"goph_keeper/internal/server/config"
	"goph_keeper/pkg/jwt"
	"log/slog"
	"net/http"

	"goph_keeper/internal/models"
	custom_errors "goph_keeper/internal/server/errors"

	"github.com/go-chi/render"
)

type GophKeeperUserServicer interface {
	RegisterUserService(ctx context.Context, user models.User) (int, error)
	GetUserService(ctx context.Context, user models.UserLoginRequest) (int, error)
}

type GophKeeperUserHandlers struct {
	service GophKeeperUserServicer
	cfg     *config.Config
}

func NewGophKeeperUserHandlers(service GophKeeperUserServicer, cfg *config.Config) *GophKeeperUserHandlers {
	return &GophKeeperUserHandlers{
		service: service,
		cfg:     cfg,
	}
}

func (gh *GophKeeperUserHandlers) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := render.Bind(r, &user); err != nil {
		slog.Error("cannot parse body", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := gh.service.RegisterUserService(r.Context(), user)
	if err != nil {
		slog.Error("some error", slog.String("error", err.Error()))
		if errors.Is(err, custom_errors.ErrUserIsExist) {
			responseError(w, r, http.StatusConflict, err.Error())
			return
		}
		responseError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	jwtToken, err := jwt.CreateJwtToken(gh.cfg.Application.SecretToken, userID)
	if err != nil {
		slog.Error("cannot create jwt token", slog.String("error", err.Error()))
		responseError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	resp := models.UserLoginResponse{
		Token: jwtToken,
	}
	responseOK(w, r, http.StatusCreated, resp, "user succesfully created")
}

func (gh *GophKeeperUserHandlers) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserLoginRequest

	if err := render.Bind(r, &user); err != nil {
		slog.Error("cannot parse body", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := gh.service.GetUserService(r.Context(), user)
	if err != nil {
		slog.Error("some error", slog.String("error", err.Error()))
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, custom_errors.ErrInvalidData) {
			responseError(w, r, http.StatusUnauthorized, err.Error())
			return
		}
		responseError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	jwtToken, err := jwt.CreateJwtToken(gh.cfg.Application.SecretToken, userID)
	if err != nil {
		slog.Error("cannot create jwt token", slog.String("error", err.Error()))
		responseError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	resp := models.UserLoginResponse{
		Token: jwtToken,
	}
	responseOK(w, r, http.StatusOK, resp, "user succesfully login")
}

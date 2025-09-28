package http_handlers

import (
	"context"
	"goph_keeper/internal/models"
	"goph_keeper/internal/server/config"
	"goph_keeper/pkg/jwt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// GophKeeperCredentialsServicer defines the interface for credentials service operations.
// It provides methods for creating, retrieving, updating, and deleting credentials.
type GophKeeperCredentialsServicer interface {
	CreateCredentials(ctx context.Context, credentials *models.Credentials) (int, error)
	GetCredentialsByUserID(ctx context.Context, userID int) ([]models.Credentials, error)
	GetCredentialsByIDAndUserID(ctx context.Context, credentialsID, userID int) ([]models.Credentials, error)
	UpdateCredentialsByIDAndUserID(ctx context.Context, credentials models.UpdateCredentialsRequest, credentialsID, userID int) error
	DeleteCredentialsByIDAndUserID(ctx context.Context, credentialsID, userID int) error
}

// GophKeeperCredentialsHandlers handles HTTP requests related to credentials management.
// It provides endpoints for creating, retrieving, updating, and deleting credentials.
type GophKeeperCredentialsHandlers struct {
	service GophKeeperCredentialsServicer
	cfg     *config.Config
}

// NewGophKeeperCredentialsHandlers creates a new instance of GophKeeperCredentialsHandlers.
func NewGophKeeperCredentialsHandlers(service GophKeeperCredentialsServicer, cfg *config.Config) *GophKeeperCredentialsHandlers {
	return &GophKeeperCredentialsHandlers{
		service: service,
		cfg:     cfg,
	}
}

// CreateCredentials handles HTTP requests for creating new credentials.
func (gh *GophKeeperCredentialsHandlers) CreateCredentials(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials
	if err := render.Bind(r, &credentials); err != nil {
		slog.Error("cannot parse body", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	claims, err := jwt.GetJwtClaims(r)
	if err != nil {
		slog.Error("some error", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		slog.Error("cannot convert user_id from token", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	credentials.UserID = userID

	credentialsID, err := gh.service.CreateCredentials(r.Context(), &credentials)
	if err != nil {
		slog.Error("some error", slog.String("error", err.Error()))
		responseError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	responseOK(w, r, http.StatusCreated, map[string]int{"id": credentialsID}, "credentials succesfully created")
}

// GetCredentialsByUserID handles HTTP requests for retrieving credentials by user ID.
func (gh *GophKeeperCredentialsHandlers) GetCredentialsByUserID(w http.ResponseWriter, r *http.Request) {
	claims, err := jwt.GetJwtClaims(r)
	if err != nil {
		slog.Error("some error", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		slog.Error("cannot convert user_id from token", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	credentials, err := gh.service.GetCredentialsByUserID(r.Context(), userID)
	if err != nil {
		slog.Error("some error", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	responseOK(w, r, http.StatusOK, credentials, "OK")
}

// GetCredentialsByIDAndUserID handles HTTP requests for retrieving credentials by ID and user ID.
func (gh *GophKeeperCredentialsHandlers) GetCredentialsByIDAndUserID(w http.ResponseWriter, r *http.Request) {
	credentialsID, userID, err := gh.getCredentialsIDAndUserID(w, r)
	if err != nil {
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	credentials, err := gh.service.GetCredentialsByIDAndUserID(r.Context(), credentialsID, userID)
	if err != nil {
		slog.Error("some error", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	responseOK(w, r, http.StatusOK, credentials, "OK")
}

// UpdateCredentialsByIDAndUserID handles HTTP requests for updating credentials by ID and user ID.
func (gh *GophKeeperCredentialsHandlers) UpdateCredentialsByIDAndUserID(w http.ResponseWriter, r *http.Request) {
	credentialsID, userID, err := gh.getCredentialsIDAndUserID(w, r)
	if err != nil {
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var credentials models.UpdateCredentialsRequest
	if err := render.Bind(r, &credentials); err != nil {
		slog.Error("cannot parse body", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := gh.service.UpdateCredentialsByIDAndUserID(r.Context(), credentials, credentialsID, userID); err != nil {
		slog.Error("some error", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	responseOK(w, r, http.StatusOK, "", "updated")
}

// DeleteCredentialsByIDAndUserID handles HTTP requests for deleting credentials by ID and user ID.
func (gh *GophKeeperCredentialsHandlers) DeleteCredentialsByIDAndUserID(w http.ResponseWriter, r *http.Request) {
	credentialsID, userID, err := gh.getCredentialsIDAndUserID(w, r)
	if err != nil {
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := gh.service.DeleteCredentialsByIDAndUserID(r.Context(), credentialsID, userID); err != nil {
		slog.Error("some error", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	responseOK(w, r, http.StatusNoContent, "", "deleted")
}

// getCredentialsIDAndUserID extracts credentials ID and user ID from the request.
func (gh *GophKeeperCredentialsHandlers) getCredentialsIDAndUserID(w http.ResponseWriter, r *http.Request) (int, int, error) {
	claims, err := jwt.GetJwtClaims(r)
	if err != nil {
		slog.Error("some error", slog.String("error", err.Error()))
		return 0, 0, err
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		slog.Error("cannot convert user_id from token", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return 0, 0, err
	}

	credentialsIDRaw := chi.URLParam(r, "credentials_id")
	if credentialsIDRaw == "" {
		slog.Error("credentials_id is required in url param")
		responseError(w, r, http.StatusBadRequest, "credentials_id is required in url param")
		return 0, 0, err
	}

	credentialsID, err := strconv.Atoi(credentialsIDRaw)
	if err != nil {
		slog.Error("cannot convert credentials_id", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return 0, 0, err
	}

	return credentialsID, userID, nil
}

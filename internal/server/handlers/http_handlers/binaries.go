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

// GophKeeperBinariesDataServicer defines the interface for binaries data service operations.
// It provides methods for creating, retrieving, updating, and deleting binaries data.
type GophKeeperBinariesDataServicer interface {
	CreateBinariesData(ctx context.Context, binariesData *models.BinariesData) (int, error)
	GetBinariesDataByUserID(ctx context.Context, userID int) ([]models.BinariesData, error)
	GetBinariesDataByIDAndUserID(ctx context.Context, binariesDataID, userID int) ([]models.BinariesData, error)
	UpdateBinariesDataByIDAndUserID(ctx context.Context, binaryData models.UpdateBinariesDataRequest, binariesDataID, userID int) error
	DeleteBinariesDataByIDAndUserID(ctx context.Context, binariesDataID, userID int) error
}

// GophKeeperBinariesDataHandlers handles HTTP requests related to binaries data management.
// It provides endpoints for creating, retrieving, updating, and deleting binaries data.
type GophKeeperBinariesDataHandlers struct {
	service GophKeeperBinariesDataServicer
	cfg     *config.Config
}

// NewGophKeeperBinariesDataHandlers creates a new instance of GophKeeperBinariesDataHandlers.
func NewGophKeeperBinariesDataHandlers(service GophKeeperBinariesDataServicer, cfg *config.Config) *GophKeeperBinariesDataHandlers {
	return &GophKeeperBinariesDataHandlers{
		service: service,
		cfg:     cfg,
	}
}

// CreateBinariesData handles HTTP requests for creating new binaries data.
func (gh *GophKeeperBinariesDataHandlers) CreateBinariesData(w http.ResponseWriter, r *http.Request) {
	var binariesData models.BinariesData
	if err := render.Bind(r, &binariesData); err != nil {
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

	binariesData.UserID = userID

	binariesDataID, err := gh.service.CreateBinariesData(r.Context(), &binariesData)
	if err != nil {
		slog.Error("some error", slog.String("error", err.Error()))
		responseError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	responseOK(w, r, http.StatusCreated, map[string]int{"id": binariesDataID}, "binaries data succesfully created")
}

// GetBinariesDataByUserID handles HTTP requests for retrieving binaries data by user ID.
func (gh *GophKeeperBinariesDataHandlers) GetBinariesDataByUserID(w http.ResponseWriter, r *http.Request) {
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

	binariesData, err := gh.service.GetBinariesDataByUserID(r.Context(), userID)
	if err != nil {
		slog.Error("some error", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	responseOK(w, r, http.StatusOK, binariesData, "OK")
}

// GetBinariesDataByIDAndUserID handles HTTP requests for retrieving binaries data by ID and user ID.
func (gh *GophKeeperBinariesDataHandlers) GetBinariesDataByIDAndUserID(w http.ResponseWriter, r *http.Request) {
	binariesDataID, userID, err := gh.getBinariesDataIDAndUserID(w, r)
	if err != nil {
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	BinariesData, err := gh.service.GetBinariesDataByIDAndUserID(r.Context(), binariesDataID, userID)
	if err != nil {
		slog.Error("some error", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	responseOK(w, r, http.StatusOK, BinariesData, "OK")
}

// UpdateBinariesDataByIDAndUserID handles HTTP requests for updating binaries data by ID and user ID.
func (gh *GophKeeperBinariesDataHandlers) UpdateBinariesDataByIDAndUserID(w http.ResponseWriter, r *http.Request) {
	binariesDataID, userID, err := gh.getBinariesDataIDAndUserID(w, r)
	if err != nil {
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var binariesData models.UpdateBinariesDataRequest
	if err := render.Bind(r, &binariesData); err != nil {
		slog.Error("cannot parse body", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := gh.service.UpdateBinariesDataByIDAndUserID(r.Context(), binariesData, binariesDataID, userID); err != nil {
		slog.Error("some error", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	responseOK(w, r, http.StatusOK, "", "updated")
}

// DeleteBinariesDataByIDAndUserID handles HTTP requests for deleting binaries data by ID and user ID.
func (gh *GophKeeperBinariesDataHandlers) DeleteBinariesDataByIDAndUserID(w http.ResponseWriter, r *http.Request) {
	binariesDataID, userID, err := gh.getBinariesDataIDAndUserID(w, r)
	if err != nil {
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := gh.service.DeleteBinariesDataByIDAndUserID(r.Context(), binariesDataID, userID); err != nil {
		slog.Error("some error", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	responseOK(w, r, http.StatusNoContent, "", "deleted")
}

// getBinariesDataIDAndUserID extracts binaries data ID and user ID from the request.
func (gh *GophKeeperBinariesDataHandlers) getBinariesDataIDAndUserID(w http.ResponseWriter, r *http.Request) (int, int, error) {
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

	binariesDataIDRaw := chi.URLParam(r, "binaries_id")
	if binariesDataIDRaw == "" {
		slog.Error("binaries_id is required in url param")
		responseError(w, r, http.StatusBadRequest, "binaries_id is required in url param")
		return 0, 0, err
	}

	binariesDataID, err := strconv.Atoi(binariesDataIDRaw)
	if err != nil {
		slog.Error("cannot convert binaries_id", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return 0, 0, err
	}

	return binariesDataID, userID, nil
}

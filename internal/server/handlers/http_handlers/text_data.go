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

type GophKeeperTextDataServicer interface {
	CreateTextData(ctx context.Context, textData *models.TextData) (int, error)
	GetTextDataByUserID(ctx context.Context, userID int) ([]models.TextData, error)
	GetTextDataByIDAndUserID(ctx context.Context, textDataID, userID int) ([]models.TextData, error)
	UpdateTextDataByIDAndUserID(ctx context.Context, textData models.UpdateTextDataRequest, textDataID, userID int) error
	DeleteTextDataByIDAndUserID(ctx context.Context, textDataID, userID int) error
}

type GophKeeperTextDataHandlers struct {
	service GophKeeperTextDataServicer
	cfg     *config.Config
}

func NewGophKeeperTextDataHandlers(service GophKeeperTextDataServicer, cfg *config.Config) *GophKeeperTextDataHandlers {
	return &GophKeeperTextDataHandlers{
		service: service,
		cfg:     cfg,
	}
}

func (gh *GophKeeperTextDataHandlers) CreateTextData(w http.ResponseWriter, r *http.Request) {
	var textData models.TextData
	if err := render.Bind(r, &textData); err != nil {
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

	textData.UserID = userID

	textDataID, err := gh.service.CreateTextData(r.Context(), &textData)
	if err != nil {
		slog.Error("some error", slog.String("error", err.Error()))
		responseError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	responseOK(w, r, http.StatusCreated, map[string]int{"id": textDataID}, "text data succesfully created")
}

func (gh *GophKeeperTextDataHandlers) GetTextDataByUserID(w http.ResponseWriter, r *http.Request) {
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

	textData, err := gh.service.GetTextDataByUserID(r.Context(), userID)
	if err != nil {
		slog.Error("some error", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	responseOK(w, r, http.StatusOK, textData, "OK")
}

func (gh *GophKeeperTextDataHandlers) GetTextDataByIDAndUserID(w http.ResponseWriter, r *http.Request) {
	textDataID, userID, err := gh.getTextDataIDAndUserID(w, r)
	if err != nil {
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	textData, err := gh.service.GetTextDataByIDAndUserID(r.Context(), textDataID, userID)
	if err != nil {
		slog.Error("some error", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	responseOK(w, r, http.StatusOK, textData, "OK")
}

func (gh *GophKeeperTextDataHandlers) UpdateTextDataByIDAndUserID(w http.ResponseWriter, r *http.Request) {
	textDataID, userID, err := gh.getTextDataIDAndUserID(w, r)
	if err != nil {
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var textData models.UpdateTextDataRequest
	if err := render.Bind(r, &textData); err != nil {
		slog.Error("cannot parse body", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := gh.service.UpdateTextDataByIDAndUserID(r.Context(), textData, textDataID, userID); err != nil {
		slog.Error("some error", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	responseOK(w, r, http.StatusOK, "", "updated")
}

func (gh *GophKeeperTextDataHandlers) DeleteTextDataByIDAndUserID(w http.ResponseWriter, r *http.Request) {
	textDataID, userID, err := gh.getTextDataIDAndUserID(w, r)
	if err != nil {
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := gh.service.DeleteTextDataByIDAndUserID(r.Context(), textDataID, userID); err != nil {
		slog.Error("some error", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	responseOK(w, r, http.StatusNoContent, "", "deleted")
}

func (gh *GophKeeperTextDataHandlers) getTextDataIDAndUserID(w http.ResponseWriter, r *http.Request) (int, int, error) {
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

	textDataIDRaw := chi.URLParam(r, "text_data_id")
	if textDataIDRaw == "" {
		slog.Error("text_data_id is required in url param")
		responseError(w, r, http.StatusBadRequest, "text_data_id is required in url param")
		return 0, 0, err
	}

	textDataID, err := strconv.Atoi(textDataIDRaw)
	if err != nil {
		slog.Error("cannot convert text_data_id", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return 0, 0, err
	}

	return textDataID, userID, nil
}

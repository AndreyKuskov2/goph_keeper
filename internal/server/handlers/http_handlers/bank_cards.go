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

// GophKeeperBankCardsServicer defines the interface for bank cards service operations.
// It provides methods for creating, retrieving, updating, and deleting bank cards.
type GophKeeperBankCardsServicer interface {
	CreateBankCards(ctx context.Context, bankCards *models.BankCards) (int, error)
	GetBankCardsByUserID(ctx context.Context, userID int) ([]models.BankCards, error)
	GetBankCardsByIDAndUserID(ctx context.Context, bankCardID, userID int) ([]models.BankCards, error)
	UpdateBankCardsByIDAndUserID(ctx context.Context, bankCard models.UpdateBankCardsRequest, bankCardID, userID int) error
	DeleteBankCardsByIDAndUserID(ctx context.Context, bankCardID, userID int) error
}

// GophKeeperBankCardsHandlers handles HTTP requests related to bank cards management.
// It provides endpoints for creating, retrieving, updating, and deleting bank cards.
type GophKeeperBankCardsHandlers struct {
	service GophKeeperBankCardsServicer
	cfg     *config.Config
}

// NewGophKeeperBankCardsHandlers creates a new instance of GophKeeperBankCardsHandlers.
func NewGophKeeperBankCardsHandlers(service GophKeeperBankCardsServicer, cfg *config.Config) *GophKeeperBankCardsHandlers {
	return &GophKeeperBankCardsHandlers{
		service: service,
		cfg:     cfg,
	}
}

// CreateBankCards handles HTTP requests for creating new bank cards.
func (gh *GophKeeperBankCardsHandlers) CreateBankCards(w http.ResponseWriter, r *http.Request) {
	var bankCards models.BankCards
	if err := render.Bind(r, &bankCards); err != nil {
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

	bankCards.UserID = userID

	bankCardsID, err := gh.service.CreateBankCards(r.Context(), &bankCards)
	if err != nil {
		slog.Error("some error", slog.String("error", err.Error()))
		responseError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	responseOK(w, r, http.StatusCreated, map[string]int{"id": bankCardsID}, "bank card succesfully created")
}

// GetBankCardsByUserID handles HTTP requests for retrieving bank cards by user ID.
func (gh *GophKeeperBankCardsHandlers) GetBankCardsByUserID(w http.ResponseWriter, r *http.Request) {
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

	bankCards, err := gh.service.GetBankCardsByUserID(r.Context(), userID)
	if err != nil {
		slog.Error("some error", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	responseOK(w, r, http.StatusOK, bankCards, "OK")
}

// GetBankCardsByIDAndUserID handles HTTP requests for retrieving bank cards by ID and user ID.
func (gh *GophKeeperBankCardsHandlers) GetBankCardsByIDAndUserID(w http.ResponseWriter, r *http.Request) {
	bankCardsID, userID, err := gh.getBankCardsIDAndUserID(w, r)
	if err != nil {
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	bankCards, err := gh.service.GetBankCardsByIDAndUserID(r.Context(), bankCardsID, userID)
	if err != nil {
		slog.Error("some error", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	responseOK(w, r, http.StatusOK, bankCards, "OK")
}

// UpdateBankCardsByIDAndUserID handles HTTP requests for updating bank cards by ID and user ID.
func (gh *GophKeeperBankCardsHandlers) UpdateBankCardsByIDAndUserID(w http.ResponseWriter, r *http.Request) {
	bankCardID, userID, err := gh.getBankCardsIDAndUserID(w, r)
	if err != nil {
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var bankCard models.UpdateBankCardsRequest
	if err := render.Bind(r, &bankCard); err != nil {
		slog.Error("cannot parse body", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := gh.service.UpdateBankCardsByIDAndUserID(r.Context(), bankCard, bankCardID, userID); err != nil {
		slog.Error("some error", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	responseOK(w, r, http.StatusOK, "", "updated")
}

// DeleteBankCardsByIDAndUserID handles HTTP requests for deleting bank cards by ID and user ID.
func (gh *GophKeeperBankCardsHandlers) DeleteBankCardsByIDAndUserID(w http.ResponseWriter, r *http.Request) {
	bankCardID, userID, err := gh.getBankCardsIDAndUserID(w, r)
	if err != nil {
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := gh.service.DeleteBankCardsByIDAndUserID(r.Context(), bankCardID, userID); err != nil {
		slog.Error("some error", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	responseOK(w, r, http.StatusNoContent, "", "deleted")
}

// getBankCardsIDAndUserID extracts bank cards ID and user ID from the request.
func (gh *GophKeeperBankCardsHandlers) getBankCardsIDAndUserID(w http.ResponseWriter, r *http.Request) (int, int, error) {
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

	bankCardIDRaw := chi.URLParam(r, "bank_cards_id")
	if bankCardIDRaw == "" {
		slog.Error("bank_cards_id is required in url param")
		responseError(w, r, http.StatusBadRequest, "bank_cards_id is required in url param")
		return 0, 0, err
	}

	bankCardID, err := strconv.Atoi(bankCardIDRaw)
	if err != nil {
		slog.Error("cannot convert bank_cards_id", slog.String("error", err.Error()))
		responseError(w, r, http.StatusBadRequest, err.Error())
		return 0, 0, err
	}

	return bankCardID, userID, nil
}

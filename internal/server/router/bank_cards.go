package router

import (
	"goph_keeper/internal/server/config"
	"goph_keeper/internal/server/handlers/http_handlers"
	"goph_keeper/internal/server/middlewares"
	"goph_keeper/internal/server/service"
	"goph_keeper/internal/server/storage"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func newBankCardsRouter(cfg *config.Config, storage storage.Storager) http.Handler {
	r := chi.NewRouter()

	r.Use(middlewares.JwtAuthValidator(cfg))

	bankCardsService := service.NewGophKeeperBankCardsService(storage)
	bankCardsHandler := http_handlers.NewGophKeeperBankCardsHandlers(bankCardsService, cfg)

	r.Post("/", bankCardsHandler.CreateBankCards)
	r.Get("/", bankCardsHandler.GetBankCardsByUserID)
	r.Get("/{bank_cards_id}", bankCardsHandler.GetBankCardsByIDAndUserID)
	r.Put("/{bank_cards_id}", bankCardsHandler.UpdateBankCardsByIDAndUserID)
	r.Delete("/{bank_cards_id}", bankCardsHandler.DeleteBankCardsByIDAndUserID)

	return r
}

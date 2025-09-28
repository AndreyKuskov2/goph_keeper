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

// Main router for the GophKeeper server.
func GophKeeperHTTPRouter(cfg *config.Config, storage storage.Storager) http.Handler {
	r := chi.NewRouter()

	r.Use(middlewares.LoggerMiddleware)

	userService := service.NewGophermartUserService(storage)
	userHandler := http_handlers.NewGophKeeperUserHandlers(userService, cfg)

	pingService := service.NewGophKeeperService(storage)
	pingHandler := http_handlers.NewGophKeeperHTTPHandler(pingService)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/ping", pingHandler.Ping)

		r.Post("/register", userHandler.RegisterUserHandler)
		r.Post("/login", userHandler.LoginUserHandler)

		r.Mount("/credentials", newCredentialsRouter(cfg, storage))
		r.Mount("/text-data", newTextDataRouter(cfg, storage))
		r.Mount("/bank-cards", newBankCardsRouter(cfg, storage))
		r.Mount("/binaries", newBinariesRouter(cfg, storage))
	})

	return r
}

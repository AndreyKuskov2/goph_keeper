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

func newCredentialsRouter(cfg *config.Config, storage storage.Storager) http.Handler {
	r := chi.NewRouter()

	r.Use(middlewares.JwtAuthValidator(cfg))

	credentialsService := service.NewGophKeeperCredentialsService(storage)
	credentialsHandler := http_handlers.NewGophKeeperCredentialsHandlers(credentialsService, cfg)

	r.Post("/", credentialsHandler.CreateCredentials)
	r.Get("/", credentialsHandler.GetCredentialsByUserID)
	r.Get("/{credentials_id}", credentialsHandler.GetCredentialsByIDAndUserID)
	r.Put("/{credentials_id}", credentialsHandler.UpdateCredentialsByIDAndUserID)
	r.Delete("/{credentials_id}", credentialsHandler.DeleteCredentialsByIDAndUserID)

	return r
}

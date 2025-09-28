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

func newTextDataRouter(cfg *config.Config, storage storage.Storager) http.Handler {
	r := chi.NewRouter()

	r.Use(middlewares.JwtAuthValidator(cfg))

	textDataService := service.NewGophKeeperTextDataService(storage)
	textDataHandler := http_handlers.NewGophKeeperTextDataHandlers(textDataService, cfg)

	r.Post("/", textDataHandler.CreateTextData)
	r.Get("/", textDataHandler.GetTextDataByUserID)
	r.Get("/{text_data_id}", textDataHandler.GetTextDataByIDAndUserID)
	r.Put("/{text_data_id}", textDataHandler.UpdateTextDataByIDAndUserID)
	r.Delete("/{text_data_id}", textDataHandler.DeleteTextDataByIDAndUserID)

	return r
}

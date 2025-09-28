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

func newBinariesRouter(cfg *config.Config, storage storage.Storager) http.Handler {
	r := chi.NewRouter()

	r.Use(middlewares.JwtAuthValidator(cfg))

	binariesService := service.NewGophKeeperBinariesDataService(storage)
	binariesHandler := http_handlers.NewGophKeeperBinariesDataHandlers(binariesService, cfg)

	r.Post("/", binariesHandler.CreateBinariesData)
	r.Get("/", binariesHandler.GetBinariesDataByUserID)
	r.Get("/{binaries_id}", binariesHandler.GetBinariesDataByIDAndUserID)
	r.Put("/{binaries_id}", binariesHandler.UpdateBinariesDataByIDAndUserID)
	r.Delete("/{binaries_id}", binariesHandler.DeleteBinariesDataByIDAndUserID)

	return r
}

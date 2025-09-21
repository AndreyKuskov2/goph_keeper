package app

import (
	"context"
	"goph_keeper/internal/server/config"
	"goph_keeper/internal/server/router"
	"goph_keeper/internal/server/storage"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	storage, err := storage.NewStorager(context.Background(), cfg)
	if err != nil {
		slog.Error("cannot init storage", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer storage.Close()

	go func() {
		routerHTTP := router.GophKeeperHTTPRouter(cfg, storage)

		slog.Info("start web-server", slog.String("address", cfg.Application.Address))
		if err := http.ListenAndServe(cfg.Application.Address, routerHTTP); err != nil {
			slog.Error("failed to start server", slog.String("error", err.Error()))
			sigChan <- syscall.SIGTERM
		}
	}()

	<-sigChan
	slog.Info("shutting down server...")
}

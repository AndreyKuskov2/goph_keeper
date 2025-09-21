package main

import (
	"goph_keeper/internal/server/app"
	"goph_keeper/internal/server/config"
	"goph_keeper/pkg/logger"
	"log/slog"
	"os"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("failed parse config or flags", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.NewLogger(cfg.Application.Debug, cfg.Application.AddSourcePath)
	slog.Info("init config")
	slog.Info("init logger")

	app.Run(cfg)
}

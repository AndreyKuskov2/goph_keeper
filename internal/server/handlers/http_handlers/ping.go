package http_handlers

import (
	"context"
	"log/slog"
	"net/http"
)

// GophKeeperPingServicer defines the interface for ping service operations.
// It provides methods for checking system health and connectivity.
type GophKeeperPingServicer interface {
	Ping(ctx context.Context) error // Ping checks if the service is available
}

// GophKeeperHTTPHandler handles HTTP requests for the GophKeeper application.
// It provides endpoints for system health checks and utility operations.
type GophKeeperHTTPHandler struct {
	service GophKeeperPingServicer // Service layer for business logic
}

// NewGophKeeperHTTPHandler creates a new instance of GophKeeperHTTPHandler.
func NewGophKeeperHTTPHandler(service GophKeeperPingServicer) *GophKeeperHTTPHandler {
	return &GophKeeperHTTPHandler{
		service: service,
	}
}

// Ping handles HTTP requests to check system health and connectivity.
func (gh *GophKeeperHTTPHandler) Ping(w http.ResponseWriter, r *http.Request) {
	if err := gh.service.Ping(r.Context()); err != nil {
		slog.Error("cannot ping storage", slog.String("error", err.Error()))
		responseError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	responseOK(w, r, http.StatusOK, map[string]string{"message": "pong"}, "")
}

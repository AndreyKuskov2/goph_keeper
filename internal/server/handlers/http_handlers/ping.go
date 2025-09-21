package http_handlers

import (
	"context"
	"log/slog"
	"net/http"
)

type GophKeeperPingServicer interface {
	Ping(ctx context.Context) error
}

type GophKeeperHTTPHandler struct {
	service GophKeeperPingServicer
}

func NewGophKeeperHTTPHandler(service GophKeeperPingServicer) *GophKeeperHTTPHandler {
	return &GophKeeperHTTPHandler{
		service: service,
	}
}

func (gh *GophKeeperHTTPHandler) Ping(w http.ResponseWriter, r *http.Request) {
	if err := gh.service.Ping(r.Context()); err != nil {
		slog.Error("cannot ping storage", slog.String("error", err.Error()))
		responseError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	responseOK(w, r, http.StatusOK, map[string]string{"message": "pong"}, "")
}

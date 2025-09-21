package middlewares

import (
	"context"
	"goph_keeper/internal/server/config"
	"goph_keeper/pkg/jwt"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

func JwtAuthValidator(cfg *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				slog.Error("no authorization token")
				render.Status(r, http.StatusUnauthorized)
				render.PlainText(w, r, "no authorization token")
				return
			}
			claims, err := jwt.VerifyToken(tokenString, cfg.Application.SecretToken)
			if err != nil {
				slog.Error("invalid token", slog.String("error", err.Error()))
				render.Status(r, http.StatusUnauthorized)
				render.PlainText(w, r, "invalid token")
				return
			}

			r = r.Clone(context.WithValue(r.Context(), jwt.ContextClaims, claims))
			next.ServeHTTP(w, r)
		})
	}
}

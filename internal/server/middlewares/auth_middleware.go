package middlewares

import (
"context"
"goph_keeper/internal/server/config"
"goph_keeper/pkg/jwt"
"log/slog"
"net/http"
"strings"

"github.com/go-chi/render"
)

func JwtAuthValidator(cfg *config.Config) func(next http.Handler) http.Handler {
return func(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
authHeader := r.Header.Get("Authorization")
if authHeader == "" {
slog.Error("no authorization token")
render.Status(r, http.StatusUnauthorized)
render.PlainText(w, r, "no authorization token")
return
}

// Extract token from "Bearer <token>" format
parts := strings.SplitN(authHeader, " ", 2)
if len(parts) != 2 || parts[0] != "Bearer" {
slog.Error("invalid authorization header format")
render.Status(r, http.StatusUnauthorized)
render.PlainText(w, r, "invalid authorization header format")
return
}

tokenString := strings.TrimSpace(parts[1])
if tokenString == "" {
slog.Error("empty token")
render.Status(r, http.StatusUnauthorized)
render.PlainText(w, r, "empty token")
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

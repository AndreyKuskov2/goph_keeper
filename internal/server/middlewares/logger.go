package middlewares

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// LoggerMiddleware creates a middleware function that logs HTTP requests and responses.
// It measures request duration, captures response status and size, and logs structured
// information about each HTTP request using slog.
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		duration := time.Since(start)

		slog.Info("request",
			slog.String("uri", r.RequestURI),
			slog.String("method", r.Method),
			slog.Duration("duration", duration),
			slog.Int("status", ww.Status()),
			slog.Int("size", ww.BytesWritten()),
		)
	})
}

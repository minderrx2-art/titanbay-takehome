package handler

import (
	"log/slog"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Logs Method && Path
		next.ServeHTTP(w, r)
		slog.Info("connection",
			"method", r.Method,
			"path", r.URL.Path,
		)
	})
}

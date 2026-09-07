package httpapi

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		startedAt := time.Now()
		next.ServeHTTP(ww, r)
		duration := time.Since(startedAt)

		requestID := middleware.GetReqID(r.Context())

		slog.Info("HTTP request", "method", r.Method, "path", r.URL.Path, "duration", duration, "request_id", requestID, "status", ww.Status())
	})
}

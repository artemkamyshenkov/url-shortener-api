package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(handler *Handler) http.Handler {
	mux := chi.NewRouter()

	mux.Get("/health", HealthHandler)
	mux.Post("/api/v1/urls", handler.CreateHandler)
	mux.Get("/api/v1/urls/{shortCode}/stats", handler.GetURLStatsByShortCode)
	mux.Get("/api/v1/urls/{shortCode}", handler.GetURLByShortCode)
	mux.Get("/{shortCode}", handler.RedirectByShortCode)
	mux.Delete("/api/v1/urls/{shortCode}", handler.DeleteByShortCode)

	return mux
}

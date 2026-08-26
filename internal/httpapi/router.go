package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(handler *Handler) http.Handler {
	mux := chi.NewRouter()

	mux.Get("/health", HealthHandler)
	mux.Post("/api/v1/urls", handler.CreateHandler)

	return mux
}

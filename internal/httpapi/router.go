package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/health", HealthHandler)

	return mux

}

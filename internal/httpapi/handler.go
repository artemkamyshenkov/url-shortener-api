package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/artemkamyshenkov/url-shortener-api/internal/shortener"
)

type Handler struct {
	service *shortener.Service
	baseURL string
}

func NewHandler(service *shortener.Service,
	baseURL string) *Handler {
	return &Handler{service: service, baseURL: baseURL}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) CreateHandler(w http.ResponseWriter,
	r *http.Request) {

	var request createURLRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid request body",
		})
		return
	}

	createdURL, err := h.service.Create(r.Context(), request.URL)

	if errors.Is(err, shortener.ErrInvalidURL) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid URL",
		})

		return
	}

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "internal server error",
		})
		return
	}

	response := urlResponse{
		ID:        createdURL.ID,
		URL:       createdURL.OriginalURL,
		ShortCode: createdURL.ShortCode,
		ShortURL:  h.baseURL + "/" + createdURL.ShortCode,
		CreatedAt: createdURL.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}

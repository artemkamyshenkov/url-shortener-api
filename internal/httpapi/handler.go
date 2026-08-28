package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/artemkamyshenkov/url-shortener-api/internal/shortener"
	"github.com/go-chi/chi/v5"
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
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	createdURL, err := h.service.Create(r.Context(), request.URL)

	if errors.Is(err, shortener.ErrInvalidURL) {
		writeJSONError(w, http.StatusBadRequest, "invalid URL")
		return
	}

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
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

func (h *Handler) GetURLByShortCode(w http.ResponseWriter,
	r *http.Request) {
	shortCode := chi.URLParam(r, "shortCode")

	data, err := h.service.GetByShortCode(r.Context(), shortCode)

	if errors.Is(err, shortener.ErrNotFound) {
		writeJSONError(w, http.StatusNotFound, "short URL not found")
		return
	}

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response := urlResponse{ID: data.ID,
		URL:       data.OriginalURL,
		ShortCode: data.ShortCode,
		ShortURL:  h.baseURL + "/" + data.ShortCode,
		CreatedAt: data.CreatedAt}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) RedirectByShortCode(w http.ResponseWriter,
	r *http.Request) {
	shortCode := chi.URLParam(r, "shortCode")

	data, err := h.service.GetByShortCode(r.Context(), shortCode)

	if errors.Is(err, shortener.ErrNotFound) {
		writeJSONError(w, http.StatusNotFound, "short URL not found")
		return
	}

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	http.Redirect(w, r, data.OriginalURL, http.StatusFound)
}

func (h *Handler) DeleteByShortCode(w http.ResponseWriter,
	r *http.Request) {
	shortCode := chi.URLParam(r, "shortCode")

	err := h.service.DeleteByShortCode(r.Context(), shortCode)

	if errors.Is(err, shortener.ErrNotFound) {
		writeJSONError(w, http.StatusNotFound, "short URL not found")
		return
	}

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

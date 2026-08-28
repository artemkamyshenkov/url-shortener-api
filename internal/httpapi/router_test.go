package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/artemkamyshenkov/url-shortener-api/internal/shortener"
)

type fakeRepository struct {
	createResult shortener.URL
	createErr    error
	getResult    shortener.URL
	getErr       error
	deleteErr    error
}

type fakeGenerator struct {
	code string
	err  error
}

func (f *fakeGenerator) Generate(length int) (string, error) {
	return f.code, f.err
}

func (r *fakeRepository) Create(ctx context.Context, originalURL, shortCode string) (shortener.URL, error) {
	return r.createResult, r.createErr
}

func (r *fakeRepository) GetByShortCode(ctx context.Context, shortCode string) (shortener.URL, error) {
	return r.getResult, r.getErr
}

func (r *fakeRepository) DeleteByShortCode(ctx context.Context, shortCode string) error {
	return r.deleteErr
}

func TestHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	NewRouter(nil).ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("expected content type application/json")
	}

	var body map[string]string

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %q", body["status"])
	}
}

func TestCreateShortURL(t *testing.T) {
	body := strings.NewReader(`{"url":"https://example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/urls", body)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	repo := &fakeRepository{
		createResult: shortener.URL{
			ID:          1,
			OriginalURL: "https://example.com",
			ShortCode:   "abc123",
		},
	}
	generator := &fakeGenerator{
		code: "abc123",
	}

	service := shortener.NewService(repo, generator, 6, 3)
	handler := NewHandler(service, "https://baseurl.com")

	NewRouter(handler).ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("expected content type application/json")
	}

	var response urlResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if response.ID != 1 {
		t.Fatalf("incorrect response ID, expect 1, got: %d", response.ID)
	}

	if response.ShortCode != "abc123" {
		t.Fatalf("incorrect response short code expect abc123, got: %s", response.ShortCode)
	}

	if response.ShortURL != "https://baseurl.com/abc123" {
		t.Fatalf("incorrect response short url expect https://baseurl.com/abc123, got: %s", response.ShortURL)
	}

	if response.URL != "https://example.com" {
		t.Fatalf("incorrect response url expect https://example.com, got: %s", response.URL)
	}
}

func TestCreateShortURL_InvalidURL(t *testing.T) {
	body := strings.NewReader(`{"url":"not-url"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/urls", body)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	repo := &fakeRepository{
		createResult: shortener.URL{
			ID:          1,
			OriginalURL: "https://example.com",
			ShortCode:   "abc123",
		},
	}
	generator := &fakeGenerator{
		code: "abc123",
	}

	service := shortener.NewService(repo, generator, 6, 3)
	handler := NewHandler(service, "https://baseurl.com")

	NewRouter(handler).ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("expected content type application/json")
	}

	var response map[string]string

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if response["error"] != "invalid URL" {
		t.Fatalf("incorrect response error")
	}
}

func TestGetURLByShortCode(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/urls/abc123", nil)

	rec := httptest.NewRecorder()

	repo := &fakeRepository{
		getResult: shortener.URL{
			ID:          1,
			OriginalURL: "https://example.com",
			ShortCode:   "abc123",
		},
	}

	generator := &fakeGenerator{
		code: "abc123",
	}

	service := shortener.NewService(repo, generator, 6, 3)
	handler := NewHandler(service, "https://baseurl.com")

	NewRouter(handler).ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("expected content type application/json")
	}

	var response urlResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if response.ID != 1 {
		t.Fatalf("incorrect response ID, expect 1, got: %d", response.ID)
	}

	if response.ShortCode != "abc123" {
		t.Fatalf("incorrect response short code expect abc123, got: %s", response.ShortCode)
	}

	if response.ShortURL != "https://baseurl.com/abc123" {
		t.Fatalf("incorrect response short url expect https://baseurl.com/abc123, got: %s", response.ShortURL)
	}

	if response.URL != "https://example.com" {
		t.Fatalf("incorrect response url expect https://example.com, got: %s", response.URL)
	}
}

func TestRedirect(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)

	rec := httptest.NewRecorder()

	repo := &fakeRepository{
		getResult: shortener.URL{
			ID:          1,
			OriginalURL: "https://example.com",
			ShortCode:   "abc123",
		},
	}

	generator := &fakeGenerator{
		code: "abc123",
	}

	service := shortener.NewService(repo, generator, 6, 3)
	handler := NewHandler(service, "https://baseurl.com")

	NewRouter(handler).ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("expected status %d, got %d", http.StatusFound, resp.StatusCode)
	}

	if resp.Header.Get("Location") != "https://example.com" {
		t.Fatalf("expected header location https://example.com")
	}
}

func TestDeleteURLByShortCode(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/urls/abc123", nil)

	rec := httptest.NewRecorder()

	repo := &fakeRepository{
		deleteErr: nil,
	}

	generator := &fakeGenerator{
		code: "abc123",
	}

	service := shortener.NewService(repo, generator, 6, 3)
	handler := NewHandler(service, "https://baseurl.com")

	NewRouter(handler).ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, resp.StatusCode)
	}

	if rec.Body.Len() != 0 {
		t.Fatalf("expected empty body")
	}
}

func TestGetURLByShortCode_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/urls/abc123", nil)

	rec := httptest.NewRecorder()

	repo := &fakeRepository{
		getErr: shortener.ErrNotFound,
	}

	generator := &fakeGenerator{
		code: "abc123",
	}

	service := shortener.NewService(repo, generator, 6, 3)
	handler := NewHandler(service, "https://baseurl.com")

	NewRouter(handler).ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
	}

	var response map[string]string

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if response["error"] != "short URL not found" {
		t.Fatalf("incorrect response error")
	}

}

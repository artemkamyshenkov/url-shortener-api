package postgres

import (
	"context"
	"errors"
	"os"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/artemkamyshenkov/url-shortener-api/internal/shortener"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testShortCodeCounter atomic.Uint64

func uniqueShortCode() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36) +
		strconv.FormatUint(testShortCodeCounter.Add(1), 36)
}

func TestURLRepository_Create(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL is not set")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	defer pool.Close()

	repository := NewURLRepository(pool)

	originalURL := "https://example.com/test"
	shortCode := uniqueShortCode()

	createdURL, err := repository.Create(ctx, originalURL, shortCode)
	if err != nil {
		t.Fatalf("create URL: %v", err)
	}

	t.Cleanup(func() {
		_, _ = pool.Exec(
			context.Background(),
			"DELETE FROM urls WHERE id = $1",
			createdURL.ID,
		)
	})

	if createdURL.OriginalURL != originalURL {
		t.Fatalf("expected correct original url")
	}

	if createdURL.ShortCode != shortCode {
		t.Fatalf("expected correct short code")
	}

	if createdURL.CreatedAt.IsZero() {
		t.Fatalf("expected correct created at")
	}

	if createdURL.ID <= 0 {
		t.Fatalf("expected positive ID, got %d", createdURL.ID)
	}
}

func TestURLRepository_GetByShortCode(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL is not set")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	defer pool.Close()

	repository := NewURLRepository(pool)

	originalURL := "https://example.com/test"
	shortCode := uniqueShortCode()

	createdURL, err := repository.Create(ctx, originalURL, shortCode)
	if err != nil {
		t.Fatalf("create URL: %v", err)
	}

	t.Cleanup(func() {
		_, _ = pool.Exec(
			context.Background(),
			"DELETE FROM urls WHERE id = $1",
			createdURL.ID,
		)
	})

	urlByShortCode, err := repository.GetByShortCode(ctx, shortCode)
	if err != nil {
		t.Fatalf("get URL: %v", err)
	}

	if createdURL.OriginalURL != urlByShortCode.OriginalURL {
		t.Fatalf("expected correct original url")
	}

	if createdURL.ShortCode != urlByShortCode.ShortCode {
		t.Fatalf("expected correct short code")
	}

	if urlByShortCode.CreatedAt.IsZero() {
		t.Fatalf("expected correct created at")
	}

	if urlByShortCode.ID != createdURL.ID {
		t.Fatalf("expected ID %d, got %d", createdURL.ID, urlByShortCode.ID)
	}
}

func TestURLRepository_GetByShortCode_NotFound(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL is not set")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	defer pool.Close()

	repository := NewURLRepository(pool)

	_, err = repository.GetByShortCode(ctx, uniqueShortCode())

	if !errors.Is(err, shortener.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestURLRepository_DeleteByShortCode(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL is not set")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	defer pool.Close()

	repository := NewURLRepository(pool)

	createdURL, err := repository.Create(ctx, "https://example.com/delete", uniqueShortCode())
	if err != nil {
		t.Fatalf("create URL: %v", err)
	}

	t.Cleanup(func() {
		_, _ = pool.Exec(
			context.Background(),
			"DELETE FROM urls WHERE id = $1",
			createdURL.ID,
		)
	})

	if err := repository.DeleteByShortCode(ctx, createdURL.ShortCode); err != nil {
		t.Fatalf("delete URL: %v", err)
	}

	_, err = repository.GetByShortCode(ctx, createdURL.ShortCode)
	if !errors.Is(err, shortener.ErrNotFound) {
		t.Fatalf("expected deleted URL to be missing, got %v", err)
	}
}

func TestURLRepository_DeleteByShortCode_NotFound(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL is not set")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	defer pool.Close()

	repository := NewURLRepository(pool)

	err = repository.DeleteByShortCode(ctx, uniqueShortCode())
	if !errors.Is(err, shortener.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

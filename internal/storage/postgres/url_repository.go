package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/artemkamyshenkov/url-shortener-api/internal/shortener"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type URLRepository struct {
	pool *pgxpool.Pool
}

func NewURLRepository(pool *pgxpool.Pool) *URLRepository {
	return &URLRepository{
		pool: pool,
	}
}

func (r *URLRepository) Create(ctx context.Context, originalURL, shortCode string) (shortener.URL, error) {
	var createdURL shortener.URL
	var pgErr *pgconn.PgError

	query := `INSERT INTO urls (original_url, short_code)
						VALUES ($1, $2)
						RETURNING id, original_url, short_code, created_at`
	row := r.pool.QueryRow(ctx, query, originalURL, shortCode)
	err := row.Scan(&createdURL.ID, &createdURL.OriginalURL, &createdURL.ShortCode, &createdURL.CreatedAt)

	if errors.As(err, &pgErr) && pgErr.Code == "23505" &&
		pgErr.ConstraintName == "urls_short_code_key" {
		return shortener.URL{}, shortener.ErrConflict
	}

	if err != nil {
		return shortener.URL{}, fmt.Errorf("create short URL: %w", err)
	}

	return createdURL, nil
}

func (r *URLRepository) GetByShortCode(ctx context.Context, shortCode string) (shortener.URL, error) {
	var url shortener.URL

	query := `SELECT id, original_url, short_code, created_at, clicks, last_accessed_at FROM urls WHERE short_code = $1`
	row := r.pool.QueryRow(ctx, query, shortCode)
	err := row.Scan(&url.ID, &url.OriginalURL, &url.ShortCode, &url.CreatedAt, &url.Clicks, &url.LastAccessedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return shortener.URL{}, shortener.ErrNotFound
	}

	if err != nil {
		return shortener.URL{}, fmt.Errorf("get URL by short code: %w", err)
	}

	return url, nil
}

func (r *URLRepository) DeleteByShortCode(ctx context.Context, shortCode string) error {
	query := `DELETE FROM urls WHERE short_code = $1`
	result, err := r.pool.Exec(ctx, query, shortCode)

	if err != nil {
		return fmt.Errorf("delete URL by short code: %w", err)
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return shortener.ErrNotFound
	}

	return nil
}

func (r *URLRepository) RecordClick(ctx context.Context, shortCode string) error {
	query := `UPDATE urls
						SET clicks = clicks + 1,
    				last_accessed_at = NOW()
						WHERE short_code = $1;`

	result, err := r.pool.Exec(ctx, query, shortCode)

	if err != nil {
		return fmt.Errorf("update URL by short code: %w", err)
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return shortener.ErrNotFound
	}

	return nil
}

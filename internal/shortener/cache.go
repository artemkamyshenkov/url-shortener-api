package shortener

import "context"

type URLCache interface {
	Get(ctx context.Context, shortCode string) (string, error)
	Set(ctx context.Context, shortCode string, originalURL string) error
	Delete(ctx context.Context, shortCode string) error
}

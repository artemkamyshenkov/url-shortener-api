package shortener

import (
	"context"
	"errors"
	"testing"
)

type fakeRepository struct {
	createResult URL
	createErr    error
	getResult    URL
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

func (r *fakeRepository) Create(ctx context.Context, originalURL, shortCode string) (URL, error) {
	return r.createResult, r.createErr
}

func (r *fakeRepository) GetByShortCode(ctx context.Context, shortCode string) (URL, error) {
	return r.getResult, r.getErr
}

func (r *fakeRepository) DeleteByShortCode(ctx context.Context, shortCode string) error {
	return r.deleteErr
}

func (r *fakeRepository) RecordClick(ctx context.Context, shortCode string) error {
	return nil
}

func TestService_Create(t *testing.T) {
	repo := &fakeRepository{
		createResult: URL{
			ID:          1,
			OriginalURL: "https://example.com",
			ShortCode:   "abc123",
		},
	}
	generator := &fakeGenerator{
		code: "abc123",
	}

	service := NewService(repo, generator, 6, 5)

	result, err := service.Create(context.Background(), "https://example.com")

	if err != nil {
		t.Fatalf("create URL: %v", err)
	}

	if result.ShortCode != "abc123" {
		t.Fatalf("expected abc123, got %s", result.ShortCode)
	}
}

func TestService_CreateIncorrectRawURL(t *testing.T) {
	repo := &fakeRepository{
		createResult: URL{
			ID:          1,
			OriginalURL: "https://example.com",
			ShortCode:   "abc123",
		},
	}
	generator := &fakeGenerator{
		code: "abc123",
	}

	service := NewService(repo, generator, 6, 5)

	_, err := service.Create(context.Background(), "example.com")

	if !errors.Is(err, ErrInvalidURL) {
		t.Fatalf("expected ErrInvalidURL, got %v", err)
	}
}

func TestService_GeneratorError(t *testing.T) {
	repo := &fakeRepository{
		createResult: URL{
			ID:          1,
			OriginalURL: "https://example.com",
			ShortCode:   "abc123",
		},
	}

	generatorErr := errors.New("generator failed")

	generator := &fakeGenerator{
		code: "abc123",
		err:  generatorErr,
	}

	service := NewService(repo, generator, 6, 5)

	_, err := service.Create(context.Background(), "https://example.com")

	if !errors.Is(err, generatorErr) {
		t.Fatalf("expected generator error, got %v", err)
	}
}

func TestService_ExhaustedAttempts(t *testing.T) {
	repo := &fakeRepository{
		createErr: ErrConflict,
	}

	generator := &fakeGenerator{
		code: "abc123",
	}

	service := NewService(repo, generator, 6, 3)

	_, err := service.Create(context.Background(), "https://example.com")

	if !errors.Is(err, ErrCreateAttemptsExhausted) {
		t.Fatalf("expected attempts error, got %v", err)
	}
}

package shortener

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type Repository interface {
	Create(ctx context.Context, originalURL, shortCode string) (URL, error)
	GetByShortCode(ctx context.Context, shortCode string) (URL, error)
	DeleteByShortCode(ctx context.Context, shortCode string) error
	RecordClick(ctx context.Context, shortCode string) error
}

type CodeGenerator interface {
	Generate(length int) (string, error)
}

type Service struct {
	repo        Repository
	cache       URLCache
	generator   CodeGenerator
	codeLength  int
	maxAttempts int
}

func NewService(repo Repository, cache URLCache,
	generator CodeGenerator,
	codeLength int,
	maxAttempts int) *Service {
	return &Service{
		repo:        repo,
		cache:       cache,
		generator:   generator,
		codeLength:  codeLength,
		maxAttempts: maxAttempts,
	}
}

func validateURL(rawURL string) (string, error) {
	cleanURL := strings.TrimSpace(rawURL)

	if cleanURL == "" {
		return "", ErrInvalidURL
	}

	u, err := url.Parse(cleanURL)

	if err != nil {
		return "", ErrInvalidURL
	}

	scheme := u.Scheme
	host := u.Host

	if scheme != "https" && scheme != "http" {
		return "", ErrInvalidURL
	}

	if host == "" {
		return "", ErrInvalidURL
	}

	return cleanURL, nil
}

func (s *Service) Create(ctx context.Context, rawURL string) (URL, error) {
	originalURL, err := validateURL(rawURL)

	if err != nil {
		return URL{}, err
	}

	for attempt := 0; attempt < s.maxAttempts; attempt++ {
		shortCode, err := s.generator.Generate(s.codeLength)

		if err != nil {
			return URL{}, fmt.Errorf("generate short code: %w", err)
		}

		newURL, err := s.repo.Create(ctx, originalURL, shortCode)

		if err == nil {
			return newURL, nil
		}

		if errors.Is(err, ErrConflict) {
			continue
		}

		return URL{}, fmt.Errorf("create short URL: %w", err)

	}

	return URL{}, ErrCreateAttemptsExhausted
}

func (s *Service) GetByShortCode(ctx context.Context, shortCode string) (URL, error) {
	foundURL, err := s.repo.GetByShortCode(ctx, shortCode)

	if err != nil {
		return URL{}, fmt.Errorf("get short URL: %w", err)
	}

	return foundURL, nil
}

func (s *Service) DeleteByShortCode(ctx context.Context, shortCode string) error {
	err := s.repo.DeleteByShortCode(ctx, shortCode)

	if err != nil {
		return fmt.Errorf("delete short URL: %w", err)
	}

	_ = s.cache.Delete(ctx, shortCode)

	return nil
}

func (s *Service) RecordClickByShortCode(ctx context.Context, shortCode string) error {
	err := s.repo.RecordClick(ctx, shortCode)

	if err != nil {
		return fmt.Errorf("record click short URL: %w", err)
	}

	return nil
}

func (s *Service) ResolveByShortCode(ctx context.Context, shortCode string) (URL, error) {
	originalURL, cachedErr := s.cache.Get(ctx, shortCode)

	if cachedErr == nil {
		err := s.RecordClickByShortCode(ctx, shortCode)

		if err != nil {
			return URL{}, err
		}

		cachedURL := URL{
			OriginalURL: originalURL,
			ShortCode:   shortCode,
		}
		return cachedURL, nil
	}

	url, err := s.GetByShortCode(ctx, shortCode)

	if err != nil {
		return URL{}, err
	}

	err = s.RecordClickByShortCode(ctx, shortCode)

	if err != nil {
		return URL{}, err
	}

	_ = s.cache.Set(ctx, shortCode, url.OriginalURL)

	return url, nil
}

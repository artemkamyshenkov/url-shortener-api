package rediscache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/artemkamyshenkov/url-shortener-api/internal/shortener"
	"github.com/redis/go-redis/v9"
)

type URLCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewURLCache(client *redis.Client,
	ttl time.Duration) *URLCache {
	return &URLCache{
		client: client,
		ttl:    ttl,
	}
}

const keyPrefix = "shortener:url:"

func (c *URLCache) Get(ctx context.Context, shortCode string) (string, error) {
	key := keyPrefix + shortCode

	originalURL, err := c.client.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return "", shortener.ErrCacheMiss
	}

	if err != nil {
		return "", fmt.Errorf("get cached URL: %w", err)
	}

	return originalURL, nil
}

func (c *URLCache) Set(ctx context.Context, shortCode string, originalURL string) error {
	key := keyPrefix + shortCode

	err := c.client.Set(ctx, key, originalURL, c.ttl).Err()

	if err != nil {
		return fmt.Errorf("set cached URL: %w", err)
	}

	return nil
}

func (c *URLCache) Delete(ctx context.Context, shortCode string) error {
	key := keyPrefix + shortCode

	err := c.client.Del(ctx, key).Err()

	if err != nil {
		return fmt.Errorf("delete cached URL: %w", err)
	}

	return nil
}

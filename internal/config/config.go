package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	HTTPPort    int
	DatabaseURL string
	BaseURL     string
	RedisAddr   string
}

func Load() (Config, error) {
	var config Config

	port := os.Getenv("HTTP_PORT")
	defaultPort := "8080"

	databaseURL := os.Getenv("DATABASE_URL")
	baseURL := os.Getenv("BASE_URL")
	redisAddr := os.Getenv("REDIS_ADDR")

	if port == "" {
		port = defaultPort
	}

	parsedPort, err := strconv.Atoi(port)

	if err != nil {
		return Config{}, fmt.Errorf("parse HTTP_PORT: %w", err)
	}

	if parsedPort < 1 || parsedPort > 65535 {
		return Config{}, errors.New("incorrect http port")
	}

	if databaseURL == "" {
		return Config{}, errors.New("incorrect database url")
	}

	if baseURL == "" {
		return Config{}, errors.New("incorrect base url")
	}

	if redisAddr == "" {
		return Config{}, errors.New("incorrect redis address")
	}

	config.HTTPPort = parsedPort
	config.DatabaseURL = databaseURL
	config.BaseURL = baseURL
	config.RedisAddr = redisAddr

	return config, nil
}

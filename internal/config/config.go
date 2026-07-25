package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	HTTPPort int
}

func Load() (Config, error) {
	var config Config

	port := os.Getenv("HTTP_PORT")
	defaultPort := "8080"

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

	config.HTTPPort = parsedPort

	return config, nil
}

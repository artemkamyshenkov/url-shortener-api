package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/artemkamyshenkov/url-shortener-api/internal/config"
	"github.com/artemkamyshenkov/url-shortener-api/internal/httpapi"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		fmt.Println("error app config load: ", err)
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)

	if err != nil {
		cancel()
		panic(err)
	}

	if err := pool.Ping(ctx); err != nil {
		cancel()
		pool.Close()
		panic(err)
	}

	cancel()
	defer pool.Close()

	mux := httpapi.NewRouter()

	address := fmt.Sprintf(":%d", cfg.HTTPPort)

	fmt.Println("Server listening on: ", address)

	serveErrCh := make(chan error, 1)

	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	go func() {
		serveErrCh <- server.ListenAndServe()
	}()

	if err := <-serveErrCh; err != nil {
		fmt.Println("Error start server", err)
		panic(err)
	}

}

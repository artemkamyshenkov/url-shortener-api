package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/artemkamyshenkov/url-shortener-api/internal/config"
	"github.com/artemkamyshenkov/url-shortener-api/internal/httpapi"
	"github.com/artemkamyshenkov/url-shortener-api/internal/shortener"
	"github.com/artemkamyshenkov/url-shortener-api/internal/storage/postgres"
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

	repository := postgres.NewURLRepository(pool)
	generator := &shortener.RandomCodeGenerator{}

	service := shortener.NewService(repository, generator, 6, 5)
	handler := httpapi.NewHandler(service, cfg.BaseURL)

	mux := httpapi.NewRouter(handler)

	address := fmt.Sprintf(":%d", cfg.HTTPPort)

	fmt.Println("Server listening on: ", address)

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	defer stop()

	serveErrCh := make(chan error, 1)

	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	go func() {
		serveErrCh <- server.ListenAndServe()
	}()

	select {
	case err := <-serveErrCh:
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("server closed:", err)
			return
		}
		fmt.Println("server error:", err)
		panic(err)
	case <-signalCtx.Done():
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		err := server.Shutdown(shutdownCtx)

		if err != nil {
			fmt.Println("shutdown error:", err)
		}

	}

}

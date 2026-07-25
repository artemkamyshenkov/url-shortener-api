package main

import (
	"fmt"
	"net/http"

	"github.com/artemkamyshenkov/url-shortener-api/internal/config"
	"github.com/artemkamyshenkov/url-shortener-api/internal/httpapi"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		fmt.Println("error app config load: ", err)
		panic(err)
	}

	mux := httpapi.NewRouter()

	address := fmt.Sprintf(":%d", cfg.HTTPPort)

	fmt.Println("Server listening on: ", address)

	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error start server", err)
		panic(err)
	}

}

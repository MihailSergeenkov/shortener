package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/MihailSergeenkov/shortener/internal/app/handlers"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.AddHandler)
	mux.HandleFunc("/{hash}", handlers.FetchHandler)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("HTTP server has encoutenred an error: %w", err)
		}
	}

	return nil
}

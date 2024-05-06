package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/MihailSergeenkov/shortener/internal/app/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.AddHandler)
		r.Get("/{hash}", handlers.FetchHandler)
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("HTTP server has encoutenred an error: %w", err)
		}
	}

	return nil
}

package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/MihailSergeenkov/shortener/internal/app/config"
	"github.com/MihailSergeenkov/shortener/internal/app/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	config.ParseFlags()
	fmt.Println("Running server on", config.Params.SAddr)

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.AddHandler)
		r.Get("/{hash}", handlers.FetchHandler)
	})

	if err := http.ListenAndServe(config.Params.SAddr, r); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("HTTP server has encoutenred an error: %w", err)
		}
	}

	return nil
}

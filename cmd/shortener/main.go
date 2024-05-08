package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/MihailSergeenkov/shortener/internal/app/config"
	"github.com/MihailSergeenkov/shortener/internal/app/routes"
	"github.com/MihailSergeenkov/shortener/internal/app/storage"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	config.ParseFlags()
	log.Printf("Running server on: %s", config.Params.RunAddr)
	s := storage.Init()
	r := routes.Init(s)

	if err := http.ListenAndServe(config.Params.RunAddr, r); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("HTTP server has encoutenred an error: %w", err)
		}
	}

	return nil
}

package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/MihailSergeenkov/shortener/internal/app/config"
	"github.com/MihailSergeenkov/shortener/internal/app/data"
	"github.com/MihailSergeenkov/shortener/internal/app/logger"
	"github.com/MihailSergeenkov/shortener/internal/app/routes"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := config.ParseFlags(); err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	if err := logger.Init(config.Params.LogLevel); err != nil {
		return fmt.Errorf("logger error: %w", err)
	}

	log.Printf("Running server on: %s", config.Params.RunAddr)

	s, err := data.NewStorage(config.Params.FileStoragePath)
	if err != nil {
		return fmt.Errorf("storage error: %w", err)
	}

	r := routes.Init(s)

	if err := http.ListenAndServe(config.Params.RunAddr, r); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("HTTP server has encoutenred an error: %w", err)
		}
	}

	return nil
}

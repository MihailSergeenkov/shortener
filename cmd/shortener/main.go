package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"go.uber.org/zap"

	"github.com/MihailSergeenkov/shortener/internal/app/config"
	"github.com/MihailSergeenkov/shortener/internal/app/data"
	"github.com/MihailSergeenkov/shortener/internal/app/logger"
	"github.com/MihailSergeenkov/shortener/internal/app/routes"
	"github.com/MihailSergeenkov/shortener/internal/app/services"
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

	l, err := logger.NewLogger(config.Params.LogLevel)
	if err != nil {
		return fmt.Errorf("logger error: %w", err)
	}

	l.Info("Running server on", zap.String("addr", config.Params.RunAddr))

	ctx := context.Background()

	s, err := data.NewStorage(ctx, l, &config.Params)
	if err != nil {
		return fmt.Errorf("storage error: %w", err)
	}

	defer func() {
		err := s.Close()

		if err != nil {
			l.Error("failed to close db connection", zap.Error(err))
		}
	}()

	r := routes.NewRouter(l, s)

	go services.BackgroundJob(ctx, l, s)

	if err := http.ListenAndServe(config.Params.RunAddr, r); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("HTTP server has encoutenred an error: %w", err)
		}
	}

	return nil
}

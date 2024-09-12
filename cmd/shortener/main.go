// Пакет main главный пакет сервиса.
package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"

	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"

	"github.com/MihailSergeenkov/shortener/internal/app/config"
	"github.com/MihailSergeenkov/shortener/internal/app/data"
	"github.com/MihailSergeenkov/shortener/internal/app/logger"
	"github.com/MihailSergeenkov/shortener/internal/app/routes"
	"github.com/MihailSergeenkov/shortener/internal/app/services"
	"github.com/go-chi/chi/v5"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	log.Printf("Build version: %s", buildVersion)
	log.Printf("Build date: %s", buildDate)
	log.Printf("Build commit: %s", buildCommit)

	if err := config.Setup(); err != nil {
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

	go services.BackgroundJob(ctx, l, s, config.Params.DropURLsPeriod)

	if config.Params.EnableHTTPS {
		err = trustedServer(l, r)
	} else {
		err = defaultServer(r)
	}
	if err != nil {
		return fmt.Errorf("start server error: %w", err)
	}

	return nil
}

func defaultServer(r chi.Router) error {
	if err := http.ListenAndServe(config.Params.RunAddr, r); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("HTTP server has encoutenred an error: %w", err)
		}
	}

	return nil
}

func trustedServer(l *zap.Logger, r chi.Router) error {
	certManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache("/tmp/certs"),
	}
	server := &http.Server{
		Addr:    ":443",
		Handler: r,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
			MinVersion:     tls.VersionTLS13,
		},
	}
	go func() {
		err := http.ListenAndServe(config.Params.RunAddr, certManager.HTTPHandler(nil))

		if err != nil {
			l.Error("failed to listen default server", zap.Error(err))
		}
	}()

	if err := server.ListenAndServeTLS("", ""); err != nil {
		return fmt.Errorf("HTTPS server has encoutenred an error: %w", err)
	}

	return nil
}

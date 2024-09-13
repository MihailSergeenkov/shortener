// Пакет main главный пакет сервиса.
package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/sync/errgroup"

	"github.com/MihailSergeenkov/shortener/internal/app/config"
	"github.com/MihailSergeenkov/shortener/internal/app/data"
	"github.com/MihailSergeenkov/shortener/internal/app/logger"
	"github.com/MihailSergeenkov/shortener/internal/app/routes"
	"github.com/MihailSergeenkov/shortener/internal/app/services"
	"github.com/go-chi/chi/v5"
)

const (
	timeoutServerShutdown = time.Second * 5
	timeoutShutdown       = time.Second * 10
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

	ctx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancelCtx()

	g, ctx := errgroup.WithContext(ctx)

	context.AfterFunc(ctx, func() {
		ctx, cancelCtx := context.WithTimeout(context.Background(), timeoutShutdown)
		defer cancelCtx()

		<-ctx.Done()
		log.Fatal("failed to gracefully shutdown the service")
	})

	s, err := data.NewStorage(ctx, l, &config.Params)
	if err != nil {
		return fmt.Errorf("storage error: %w", err)
	}

	g.Go(func() error {
		defer log.Print("closed DB")

		<-ctx.Done()

		if err := s.Close(); err != nil {
			l.Error("failed to close db connection", zap.Error(err))
		}
		return nil
	})

	r := routes.NewRouter(l, s)

	go services.BackgroundJob(ctx, l, s, config.Params.DropURLsPeriod)

	srv := configureServer(r)

	g.Go(func() error {
		defer func() {
			errRec := recover()
			if errRec != nil {
				err = fmt.Errorf("a panic occurred: %v", errRec)
				l.Error("failed", zap.Error(err))
			}
		}()
		if err := runServer(srv); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				return fmt.Errorf("HTTP server has encoutenred an error: %w", err)
			}
		}
		return nil
	})

	g.Go(func() error {
		defer log.Print("server has been shutdown")
		<-ctx.Done()

		shutdownTimeoutCtx, cancelShutdownTimeoutCtx := context.WithTimeout(context.Background(), timeoutServerShutdown)
		defer cancelShutdownTimeoutCtx()
		if err := srv.Shutdown(shutdownTimeoutCtx); err != nil {
			log.Printf("an error occurred during server shutdown: %v", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("some errorgroup error: %w", err)
	}

	return nil
}

func configureServer(r chi.Router) *http.Server {
	if config.Params.EnableHTTPS {
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache("/tmp/certs"),
			HostPolicy: autocert.HostWhitelist("mynetwork.keenetic.link"),
		}
		server := &http.Server{
			Addr:    config.Params.RunAddr,
			Handler: r,
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
				MinVersion:     tls.VersionTLS13,
			},
		}

		return server
	}

	return &http.Server{
		Addr:    config.Params.RunAddr,
		Handler: r,
	}
}

func runServer(srv *http.Server) error {
	var err error

	if config.Params.EnableHTTPS {
		err = srv.ListenAndServeTLS("", "")
	} else {
		err = srv.ListenAndServe()
	}
	if err != nil {
		return fmt.Errorf("listen and server has failed:: %w", err)
	}

	return nil
}

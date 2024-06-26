package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
	"go.uber.org/zap/zapcore"
)

type Settings struct {
	RunAddr         string        `env:"SERVER_ADDRESS"`
	BaseURL         string        `env:"BASE_URL"`
	FileStoragePath string        `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string        `env:"DATABASE_DSN"`
	SecretKey       string        `env:"SECRET_KEY"`
	LogLevel        zapcore.Level `env:"LOG_LEVEL"`
}

var Params Settings = Settings{LogLevel: zapcore.ErrorLevel}

func ParseFlags() error {
	flag.StringVar(&Params.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&Params.BaseURL, "b", "http://localhost:8080", "address and port to urls")
	flag.Func("l", `level for logger (default "ERROR")`, func(s string) error {
		lev, err := zapcore.ParseLevel(s)

		if err != nil {
			return fmt.Errorf("parse log level env error: %w", err)
		}

		Params.LogLevel = lev
		return nil
	})

	flag.StringVar(&Params.FileStoragePath, "f", "/tmp/short-url-db.json", "file storage path")
	flag.StringVar(&Params.DatabaseDSN, "d", "", "database DSN")
	flag.StringVar(&Params.SecretKey, "s", "1234567890", "secret key for generate cookie token")

	flag.Parse()

	err := env.Parse(&Params)

	if err != nil {
		return fmt.Errorf("env error: %w", err)
	}

	return nil
}

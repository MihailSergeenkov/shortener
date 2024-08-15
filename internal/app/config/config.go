package config

import (
	"flag"
	"fmt"
	"net/url"

	"github.com/caarlos0/env/v11"
	"go.uber.org/zap/zapcore"
)

type Settings struct {
	RunAddr         string        `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         url.URL       `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string        `env:"FILE_STORAGE_PATH" envDefault:"/tmp/short-url-db.json"`
	DatabaseDSN     string        `env:"DATABASE_DSN" envDefault:""`
	SecretKey       string        `env:"SECRET_KEY" envDefault:"1234567890"`
	LogLevel        zapcore.Level `env:"LOG_LEVEL" envDefault:"ERROR"`
}

var Params Settings

func init() {
	Params = Settings{
		LogLevel: zapcore.ErrorLevel,
	}
}

func ParseFlags() error {
	err := env.Parse(&Params)

	if err != nil {
		return fmt.Errorf("env error: %w", err)
	}

	flag.StringVar(&Params.RunAddr, "a", Params.RunAddr, "address and port to run server")
	flag.Func("b", `address and port to urls (default "http://localhost:8080")`, func(s string) error {
		parsedBaseURL, err := url.Parse(s)
		if err != nil {
			return fmt.Errorf("parse user base url env error: %w", err)
		}

		Params.BaseURL = *parsedBaseURL
		return nil
	})
	flag.Func("l", `level for logger (default "ERROR")`, func(s string) error {
		lev, err := zapcore.ParseLevel(s)

		if err != nil {
			return fmt.Errorf("parse log level env error: %w", err)
		}

		Params.LogLevel = lev
		return nil
	})

	flag.StringVar(&Params.FileStoragePath, "f", Params.FileStoragePath, "file storage path")
	flag.StringVar(&Params.DatabaseDSN, "d", Params.DatabaseDSN, "database DSN")
	flag.StringVar(&Params.SecretKey, "s", Params.SecretKey, "secret key for generate cookie token")

	flag.Parse()

	return nil
}

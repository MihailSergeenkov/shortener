// Пакет config предназначен для конфигурирования сервиса.
package config

import (
	"flag"
	"fmt"
	"net/url"
	"time"

	"github.com/caarlos0/env/v11"
	"go.uber.org/zap/zapcore"
)

// Settings структура для конфигурирования сервиса.
type Settings struct {
	BaseURL         url.URL       `env:"BASE_URL" envDefault:"http://localhost:8080"`
	RunAddr         string        `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	FileStoragePath string        `env:"FILE_STORAGE_PATH" envDefault:"/tmp/short-url-db.json"`
	DatabaseDSN     string        `env:"DATABASE_DSN" envDefault:""`
	SecretKey       string        `env:"SECRET_KEY" envDefault:"1234567890"`
	DropURLsPeriod  time.Duration `env:"DROP_URLS_PERIOD" envDefault:"1m"`
	LogLevel        zapcore.Level `env:"LOG_LEVEL" envDefault:"ERROR"`
	EnableHTTPS     bool          `env:"ENABLE_HTTPS" envDefault:"false"`
}

// Params глобальная переменная типа Settings, инициализируется в момент старта сервиса.
var Params Settings

func init() {
	Params = Settings{
		LogLevel: zapcore.ErrorLevel,
	}
}

// Setup функция считывания и применения пользовательских настроек сервиса.
func Setup() error {
	if err := Params.parseEnv(); err != nil {
		return fmt.Errorf("failed to parse envs: %w", err)
	}

	Params.parseFlags()

	return nil
}

func (s *Settings) parseEnv() error {
	err := env.Parse(&Params)
	if err != nil {
		return fmt.Errorf("env error: %w", err)
	}

	return nil
}

func (s *Settings) parseFlags() {
	flag.StringVar(&s.RunAddr, "a", s.RunAddr, "address and port to run server")
	flag.Func("b", `address and port to urls (default "http://localhost:8080")`, func(v string) error {
		parsedBaseURL, err := url.Parse(v)
		if err != nil {
			return fmt.Errorf("parse user base url env error: %w", err)
		}

		s.BaseURL = *parsedBaseURL
		return nil
	})
	flag.Func("l", `level for logger (default "ERROR")`, func(v string) error {
		lev, err := zapcore.ParseLevel(v)

		if err != nil {
			return fmt.Errorf("parse log level env error: %w", err)
		}

		s.LogLevel = lev
		return nil
	})

	flag.StringVar(&s.FileStoragePath, "f", s.FileStoragePath, "file storage path")
	flag.StringVar(&s.DatabaseDSN, "d", s.DatabaseDSN, "database DSN")
	flag.StringVar(&s.SecretKey, "sk", s.SecretKey, "secret key for generate cookie token")
	flag.DurationVar(&s.DropURLsPeriod, "dp", s.DropURLsPeriod, "drop urls period")
	flag.BoolVar(&s.EnableHTTPS, "s", s.EnableHTTPS, "enable HTTPS")

	flag.Parse()
}

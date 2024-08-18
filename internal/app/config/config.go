// Пакет config предназначен для конфигурирования сервиса.
package config

import (
	"flag"
	"fmt"
	"net/url"

	"github.com/caarlos0/env/v11"
	"go.uber.org/zap/zapcore"
)

// Settings структура для конфигурирования сервиса.
type Settings struct {
	RunAddr         string        `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`            // адрес и порт сервиса
	BaseURL         url.URL       `env:"BASE_URL" envDefault:"http://localhost:8080"`           // URL для коротких ссылок
	FileStoragePath string        `env:"FILE_STORAGE_PATH" envDefault:"/tmp/short-url-db.json"` // путь до файловой БД
	DatabaseDSN     string        `env:"DATABASE_DSN" envDefault:""`                            // адрес БД
	SecretKey       string        `env:"SECRET_KEY" envDefault:"1234567890"`                    // секретный ключ
	LogLevel        zapcore.Level `env:"LOG_LEVEL" envDefault:"ERROR"`                          // уровень логирования
}

// Params глобальная переменная типа Settings, инициализируется в момент старта сервиса.
var Params Settings

func init() {
	Params = Settings{
		LogLevel: zapcore.ErrorLevel,
	}
}

// ParseFlags функция считывания и применения пользовательских настроек сервиса.
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

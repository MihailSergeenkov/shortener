// Пакет config предназначен для конфигурирования сервиса.
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"time"

	"github.com/caarlos0/env/v11"
	"go.uber.org/zap/zapcore"
)

// Settings структура для конфигурирования сервиса.
type Settings struct {
	BaseURL         url.URL       `json:"base_url" env:"BASE_URL" envDefault:"http://localhost:8080"`
	RunAddr         string        `json:"server_address" env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	FileStoragePath string        `json:"file_storage_path" env:"FILE_STORAGE_PATH" envDefault:"/tmp/url-db.json"`
	DatabaseDSN     string        `json:"database_dsn" env:"DATABASE_DSN" envDefault:""`
	SecretKey       string        `json:"secret_key" env:"SECRET_KEY" envDefault:"1234567890"`
	DropURLsPeriod  time.Duration `json:"drop_urls_period" env:"DROP_URLS_PERIOD" envDefault:"1m"`
	LogLevel        zapcore.Level `json:"log_level" env:"LOG_LEVEL" envDefault:"ERROR"`
	EnableHTTPS     bool          `json:"enable_https" env:"ENABLE_HTTPS" envDefault:"false"`
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
	configData, presentData, err := getConfigData()
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	if presentData {
		if err := parseConfigData(configData); err != nil {
			return fmt.Errorf("failed to parse config data: %w", err)
		}
	}

	if err := Params.parseEnv(); err != nil {
		return fmt.Errorf("failed to parse envs: %w", err)
	}

	Params.parseFlags()

	return nil
}

func getConfigData() ([]byte, bool, error) {
	configFile := os.Getenv("CONFIG")

	for i, arg := range os.Args {
		if arg == "-c" || arg == "-config" {
			configFile = os.Args[i+1]
			break
		}
	}

	if configFile == "" {
		return []byte{}, false, nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return []byte{}, false, fmt.Errorf("failed to read config file: %w", err)
	}

	return data, true, nil
}

func parseConfigData(data []byte) error {
	config := struct {
		BaseURL         string `json:"base_url" env:"BASE_URL"`
		RunAddr         string `json:"server_address" env:"SERVER_ADDRESS"`
		FileStoragePath string `json:"file_storage_path" env:"FILE_STORAGE_PATH"`
		DatabaseDSN     string `json:"database_dsn" env:"DATABASE_DSN"`
		SecretKey       string `json:"secret_key" env:"SECRET_KEY"`
		DropURLsPeriod  string `json:"drop_urls_period" env:"DROP_URLS_PERIOD"`
		LogLevel        string `json:"log_level" env:"LOG_LEVEL"`
		EnableHTTPS     string `json:"enable_https" env:"ENABLE_HTTPS"`
	}{}

	err := json.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	vc := reflect.ValueOf(config)
	tc := vc.Type()
	for i := range tc.NumField() {
		field := tc.Field(i)
		envName := field.Tag.Get("env")

		if _, envPresent := os.LookupEnv(envName); envPresent {
			continue
		}

		if err := os.Setenv(envName, vc.Field(i).String()); err != nil {
			return fmt.Errorf("failed to set env from config: %w", err)
		}
	}

	return nil
}

func (s *Settings) parseEnv() error {
	err := env.Parse(s)
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

	flag.String("c", "", "config file path (shorthand)")
	flag.String("config", "", "config file path")

	flag.Parse()
}

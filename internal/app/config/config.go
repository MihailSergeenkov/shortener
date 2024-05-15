package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Settings struct {
	RunAddr string `env:"SERVER_ADDRESS"`
	BaseURL string `env:"BASE_URL"`
}

var Params Settings

func ParseFlags() error {
	flag.StringVar(&Params.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&Params.BaseURL, "b", "http://localhost:8080", "address and port to urls")

	flag.Parse()

	err := env.Parse(&Params)

	if err != nil {
		return fmt.Errorf("env error: %w", err)
	}

	return nil
}

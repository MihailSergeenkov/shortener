package config

import (
	"flag"
	"os"
)

type Settings struct {
	RunAddr string
	BaseURL string
}

var Params Settings

func ParseFlags() {
	flag.StringVar(&Params.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&Params.BaseURL, "b", "http://localhost:8080", "address and port to urls")

	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		Params.RunAddr = envRunAddr
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		Params.BaseURL = envBaseURL
	}
}

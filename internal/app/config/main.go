package config

import (
	"flag"
)

type Settings struct {
	SAddr string
	UAddr string
}

var Params Settings

func ParseFlags() {
	flag.StringVar(&Params.SAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&Params.UAddr, "b", "http://localhost:8080", "address and port to urls")

	flag.Parse()
}

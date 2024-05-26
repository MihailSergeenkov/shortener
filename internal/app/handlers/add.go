package handlers

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/MihailSergeenkov/shortener/internal/app/config"
	"github.com/MihailSergeenkov/shortener/internal/app/data"
)

func AddHandler(s data.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to read request body: %v", err)
			return
		}

		u, err := s.AddURL(string(body))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to add URL to storage: %v", err)
			return
		}

		result, err := url.JoinPath(config.Params.BaseURL, u.ShortURL)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to construct URL: %v", err)
			return
		}

		w.WriteHeader(http.StatusCreated)

		_, err = w.Write([]byte(result))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to write response body: %v", err)
			return
		}
	}
}

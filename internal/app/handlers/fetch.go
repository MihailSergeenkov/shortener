package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/MihailSergeenkov/shortener/internal/app/data"
)

func FetchHandler(s data.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		short_url := strings.TrimLeft(r.URL.Path, "/")
		u, err := s.FetchURL(short_url)

		if err != nil {
			if errors.Is(err, data.ErrURLNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to fetch URL from storage: %v", err)
			return
		}

		w.Header().Set("Location", u.OriginalUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

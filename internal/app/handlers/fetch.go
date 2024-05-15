package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/MihailSergeenkov/shortener/internal/app/storage"
)

func FetchHandler(urls storage.Urls) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimLeft(r.URL.Path, "/")
		u, err := urls.FetchURL(id)

		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to fetch URL from storage: %v", err)
			return
		}

		w.Header().Set("Location", u)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

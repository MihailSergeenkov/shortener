package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/MihailSergeenkov/shortener/internal/app/data"
	"go.uber.org/zap"
)

func FetchHandler(l *zap.Logger, s data.Storager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := strings.TrimLeft(r.URL.Path, "/")
		u, err := s.FetchURL(shortURL)

		if err != nil {
			if errors.Is(err, data.ErrURLNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			l.Error("failed to fetch URL from storage", zap.Error(err))
			return
		}

		w.Header().Set("Location", u.OriginalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

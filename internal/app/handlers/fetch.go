package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/MihailSergeenkov/shortener/internal/app/common"
	"github.com/MihailSergeenkov/shortener/internal/app/data"
	"github.com/MihailSergeenkov/shortener/internal/app/services"
)

// FetchHandler обработчик получения оригинальной ссылки по короткой.
func FetchHandler(l *zap.Logger, s data.Storager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := strings.TrimLeft(r.URL.Path, "/")
		u, err := s.GetURL(r.Context(), shortURL)

		if err != nil {
			if errors.Is(err, data.ErrURLNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			l.Error("failed to fetch URL from storage", zap.Error(err))
			return
		}

		if u.DeletedFlag {
			w.WriteHeader(http.StatusGone)
			return
		}

		w.Header().Set("Location", u.OriginalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

// APIFetchUserURLsHandler обработчик получения всех сохраненных ссылок пользователя для API.
func APIFetchUserURLsHandler(l *zap.Logger, s data.Storager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := services.FetchUserURLs(r.Context(), s)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			l.Error("failed to fetch URLs from storage", zap.Error(err))
			return
		}

		if len(resp) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set(common.ContentTypeHeader, common.JSONContentType)
		w.WriteHeader(http.StatusOK)

		enc := json.NewEncoder(w)
		if err := enc.Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			l.Error(common.EncRespErrStr, zap.Error(err))
			return
		}
	}
}

package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/MihailSergeenkov/shortener/internal/app/common"
	"github.com/MihailSergeenkov/shortener/internal/app/data"
	"github.com/MihailSergeenkov/shortener/internal/app/services"
	"go.uber.org/zap"
)

func FetchHandler(l *zap.Logger, s data.Storager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := strings.TrimLeft(r.URL.Path, "/")
		u, err := s.GetOriginalURL(r.Context(), shortURL)

		if err != nil {
			if errors.Is(err, data.ErrURLNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			l.Error("failed to fetch URL from storage", zap.Error(err))
			return
		}

		w.Header().Set("Location", u)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

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

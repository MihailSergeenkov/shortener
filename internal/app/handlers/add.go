package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/MihailSergeenkov/shortener/internal/app/config"
	"github.com/MihailSergeenkov/shortener/internal/app/data"
	"github.com/MihailSergeenkov/shortener/internal/app/models"
	"github.com/MihailSergeenkov/shortener/internal/app/services"
	"go.uber.org/zap"
)

func AddHandler(l *zap.Logger, s data.Storager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			l.Error("failed to read request body", zap.Error(err))
			return
		}

		shortURL, err := services.AddShortURL(s, string(body))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			l.Error("failed to add URL to storage", zap.Error(err))
			return
		}

		result, err := url.JoinPath(config.Params.BaseURL, shortURL)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			l.Error("failed to construct URL", zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusCreated)

		_, err = w.Write([]byte(result))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			l.Error("failed to write response body", zap.Error(err))
			return
		}
	}
}

func APIAddHandler(l *zap.Logger, s data.Storager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.Request
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			l.Error("failed to read request body", zap.Error(err))
			return
		}

		shortURL, err := services.AddShortURL(s, req.URL)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			l.Error("failed to add URL to storage", zap.Error(err))
			return
		}

		result, err := url.JoinPath(config.Params.BaseURL, shortURL)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			l.Error("failed to construct URL", zap.Error(err))
			return
		}

		resp := models.Response{Result: result}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		enc := json.NewEncoder(w)
		if err := enc.Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			l.Error("error encoding response", zap.Error(err))
			return
		}
	}
}

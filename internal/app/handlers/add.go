package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path"

	"github.com/MihailSergeenkov/shortener/internal/app/common"
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
			l.Error(common.ReadReqErrStr, zap.Error(err))
			return
		}

		shortURL, err := services.AddShortURL(r.Context(), s, string(body))

		if err != nil {
			var origErr *data.OriginalURLAlreadyExistError
			if errors.As(err, &origErr) {
				w.WriteHeader(http.StatusConflict)
				result := path.Join(config.Params.BaseURL.String(), origErr.ShortURL)
				_, err = w.Write([]byte(result))

				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					l.Error("failed to write response body", zap.Error(err))
					return
				}
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			l.Error("failed to add URL to storage", zap.Error(err))
			return
		}

		result := path.Join(config.Params.BaseURL.String(), shortURL)
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
			l.Error(common.ReadReqErrStr, zap.Error(err))
			return
		}

		shortURL, err := services.AddShortURL(r.Context(), s, req.URL)

		if err != nil {
			var origErr *data.OriginalURLAlreadyExistError
			if errors.As(err, &origErr) {
				w.Header().Set(common.ContentTypeHeader, common.JSONContentType)
				w.WriteHeader(http.StatusConflict)
				result := path.Join(config.Params.BaseURL.String(), origErr.ShortURL)
				resp := models.Response{Result: result}

				enc := json.NewEncoder(w)
				if err := enc.Encode(resp); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					l.Error(common.EncRespErrStr, zap.Error(err))
					return
				}
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			l.Error("failed to add URL to storage", zap.Error(err))
			return
		}

		result := path.Join(config.Params.BaseURL.String(), shortURL)
		resp := models.Response{Result: result}

		w.Header().Set(common.ContentTypeHeader, common.JSONContentType)
		w.WriteHeader(http.StatusCreated)

		enc := json.NewEncoder(w)
		if err := enc.Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			l.Error(common.EncRespErrStr, zap.Error(err))
			return
		}
	}
}

func APIAddBatchHandler(l *zap.Logger, s data.Storager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.BatchRequest

		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			l.Error(common.ReadReqErrStr, zap.Error(err))
			return
		}

		resp, err := services.AddBatchShortURL(r.Context(), s, req)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			l.Error("failed to add URLs to storage", zap.Error(err))
			return
		}

		w.Header().Set(common.ContentTypeHeader, common.JSONContentType)
		w.WriteHeader(http.StatusCreated)

		enc := json.NewEncoder(w)
		if err := enc.Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			l.Error(common.EncRespErrStr, zap.Error(err))
			return
		}
	}
}

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MihailSergeenkov/shortener/internal/app/common"
	"github.com/MihailSergeenkov/shortener/internal/app/data"
	"github.com/MihailSergeenkov/shortener/internal/app/services"
	"go.uber.org/zap"
)

// APIFetchStatsHandler обработчик получения статистических данных.
func APIFetchStatsHandler(l *zap.Logger, s data.Storager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := services.FetchStats(r.Context(), s)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			l.Error("failed to fetch stats from storage", zap.Error(err))
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

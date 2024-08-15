package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/MihailSergeenkov/shortener/internal/app/common"
	"github.com/MihailSergeenkov/shortener/internal/app/data"
	"github.com/MihailSergeenkov/shortener/internal/app/services"
)

func APIDeleteUserURLsHandler(l *zap.Logger, s data.Storager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req []string
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			l.Error(common.ReadReqErrStr, zap.Error(err))
			return
		}

		err := services.DeleteUserURLs(r.Context(), l, s, req)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			l.Error("failed to delete URLs from storage", zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

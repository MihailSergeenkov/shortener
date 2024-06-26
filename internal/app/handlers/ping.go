package handlers

import (
	"net/http"

	"github.com/MihailSergeenkov/shortener/internal/app/data"
	"go.uber.org/zap"
)

func PingHandler(l *zap.Logger, s data.Storager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.Ping(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			l.Error("failed to connect to DB", zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

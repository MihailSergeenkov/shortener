package services

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/MihailSergeenkov/shortener/internal/app/data"
)

const dropPeriod = 10 // in minutes

// BackgroundJob функция запуска отложенных задач сервиса (очистка из БД удаленных ссылкок).
func BackgroundJob(ctx context.Context, l *zap.Logger, s data.Storager) {
	ticker := time.NewTicker(dropPeriod * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := s.DropDeletedURLs(ctx)

			if err != nil {
				l.Error("failed to drop URLs from storage", zap.Error(err))
			}
		}
	}
}

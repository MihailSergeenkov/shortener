package services

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/MihailSergeenkov/shortener/internal/app/data"
)

// BackgroundJob функция запуска отложенных задач сервиса (очистка из БД удаленных ссылкок).
func BackgroundJob(ctx context.Context, l *zap.Logger, s data.Storager, dropPeriod time.Duration) {
	ticker := time.NewTicker(dropPeriod)

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

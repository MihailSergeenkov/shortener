package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MihailSergeenkov/shortener/internal/app/data/mock"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestBackgroundJob_Ok(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	logger := zap.NewNop()
	storage := mock.NewMockStorager(mockCtrl)
	dropPeriod := 100 * time.Millisecond

	t.Run("success run", func(t *testing.T) {
		storage.EXPECT().DropDeletedURLs(ctx).AnyTimes().Return(nil)

		BackgroundJob(ctx, logger, storage, dropPeriod)
	})
}

func TestBackgroundJob_Failed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	logger := zap.NewNop()
	storage := mock.NewMockStorager(mockCtrl)
	dropPeriod := 100 * time.Millisecond
	errSome := errors.New("some error")

	t.Run("failed run", func(t *testing.T) {
		storage.EXPECT().DropDeletedURLs(ctx).AnyTimes().Return(errSome)

		BackgroundJob(ctx, logger, storage, dropPeriod)
	})
}

func TestBackgroundJob_CtxDone(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	logger := zap.NewNop()
	storage := mock.NewMockStorager(mockCtrl)
	dropPeriod := 1 * time.Minute

	t.Run("ctx Done", func(t *testing.T) {
		storage.EXPECT().DropDeletedURLs(ctx).Times(0)
		cancel()
		BackgroundJob(ctx, logger, storage, dropPeriod)
	})
}

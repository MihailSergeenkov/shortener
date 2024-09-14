package data

import (
	"context"
	"testing"

	"github.com/MihailSergeenkov/shortener/internal/app/data/mock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestTraceQueryStart(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := mock.NewMockLogger(mockCtrl)
	tracer := queryTracer{logger: logger}
	ctx := context.Background()
	data := pgx.TraceQueryStartData{}

	logger.EXPECT().Info("Running query", zap.String("query", data.SQL), zap.Any("args", data.Args)).Times(1)

	t.Run("log start", func(t *testing.T) {
		returnCtx := tracer.TraceQueryStart(ctx, nil, data)
		assert.Equal(t, ctx, returnCtx)
	})
}

func TestTraceQueryEnd(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := mock.NewMockLogger(mockCtrl)
	tracer := queryTracer{logger: logger}
	ctx := context.Background()
	data := pgx.TraceQueryEndData{}

	logger.EXPECT().Info("End query", zap.Any("tag", data.CommandTag)).Times(1)

	t.Run("log end", func(t *testing.T) {
		tracer.TraceQueryEnd(ctx, nil, data)
	})
}

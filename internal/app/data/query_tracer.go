package data

import (
	"context"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger интерфейс к логгеру.
type Logger interface {
	Info(msg string, fields ...zapcore.Field)
}

type queryTracer struct {
	logger Logger
}

// TraceQueryStart метод начала трасировки запроса.
func (t *queryTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	t.logger.Info("Running query", zap.String("query", data.SQL), zap.Any("args", data.Args))
	return ctx
}

// TraceQueryEnd метод окончания трасировки запроса.
func (t *queryTracer) TraceQueryEnd(_ context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	t.logger.Info("End query", zap.Any("tag", data.CommandTag))
}

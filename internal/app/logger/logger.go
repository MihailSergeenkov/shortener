package logger

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger = zap.NewNop()

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size

	return size, err //nolint:wrapcheck // Нужно обернуть, но возврат должен остаться оригинальным
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func Init(level zapcore.Level) error {
	cfg := zap.NewProductionConfig()
	cfg.Level.SetLevel(level)

	logger, err := cfg.Build()
	if err != nil {
		return fmt.Errorf("logger build failed: %w", err)
	}
	defer func() {
		err := logger.Sync()

		if err != nil {
			log.Printf("failed to sync logger: %v", err)
		}
	}()

	Log = logger
	return nil
}

func WithRequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		Log.Info("got incoming HTTP request",
			zap.String("uri", uri),
			zap.String("method", method),
			zap.String("duration", duration.String()),
			zap.Int("status", responseData.status),
			zap.Int("size", responseData.size),
		)
	})
}

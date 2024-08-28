package data

import (
	"context"
	"testing"

	"github.com/MihailSergeenkov/shortener/internal/app/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewStorage(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Settings
	}{
		{
			name: "new base storage",
			config: &config.Settings{
				LogLevel:        zapcore.ErrorLevel,
				DatabaseDSN:     "",
				FileStoragePath: "",
			},
		},
		{
			name: "new file storage",
			config: &config.Settings{
				LogLevel:        zapcore.ErrorLevel,
				DatabaseDSN:     "",
				FileStoragePath: "/tmp/short-url-db.json",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			logger := zap.NewNop()
			storage, err := NewStorage(ctx, logger, test.config)

			require.NoError(t, err)
			assert.Implements(t, (*Storager)(nil), storage)
		})
	}
}

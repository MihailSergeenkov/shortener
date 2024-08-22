package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestNewLogger(t *testing.T) {
	t.Run("init logger", func(t *testing.T) {
		logLevel := zapcore.ErrorLevel

		logger, err := NewLogger(logLevel)

		require.NoError(t, err)
		assert.Equal(t, logLevel, logger.Level())
	})
}

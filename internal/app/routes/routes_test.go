package routes

import (
	"testing"

	"github.com/MihailSergeenkov/shortener/internal/app/data/mock"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewRouter(t *testing.T) {
	t.Run("init router", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		logger := zap.NewNop()
		storage := mock.NewMockStorager(mockCtrl)

		r := NewRouter(logger, storage)
		assert.Implements(t, (*chi.Router)(nil), r)
	})
}

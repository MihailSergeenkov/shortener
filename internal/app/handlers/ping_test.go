package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MihailSergeenkov/shortener/internal/app/data/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSuccessPing(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	l := zap.NewNop()
	s := mock.NewMockStorager(mockCtrl)

	_ = s.EXPECT().Ping(gomock.Any()).Times(1).Return(nil)

	t.Run("ping success", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/ping", http.NoBody)
		w := httptest.NewRecorder()
		PingHandler(l, s)(w, request)

		res := w.Result()
		defer closeBody(t, res)

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func TestFailedPing(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	l := zap.NewNop()
	s := mock.NewMockStorager(mockCtrl)

	errSome := errors.New("some error")
	_ = s.EXPECT().Ping(gomock.Any()).Times(1).Return(errSome)

	t.Run("ping failed", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/ping", http.NoBody)
		w := httptest.NewRecorder()
		PingHandler(l, s)(w, request)

		res := w.Result()
		defer closeBody(t, res)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}

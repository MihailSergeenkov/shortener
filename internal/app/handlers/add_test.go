package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MihailSergeenkov/shortener/internal/app/common"
	"github.com/MihailSergeenkov/shortener/internal/app/data"
	"github.com/MihailSergeenkov/shortener/internal/app/data/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAddHandler_Ok(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zap.NewNop()
	storage := mock.NewMockStorager(mockCtrl)

	type request struct {
		body string
	}
	type want struct {
		code int
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "success add url",
			request: request{
				body: "https://ya.ru/some",
			},
			want: want{
				code: http.StatusCreated,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.EXPECT().StoreShortURL(gomock.Any(), gomock.Any(), test.request.body).Times(1).Return(nil)

			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.request.body))
			w := httptest.NewRecorder()
			AddHandler(logger, storage)(w, request)

			res := w.Result()
			defer closeBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.NotEmpty(t, resBody)
		})
	}
}

func TestAddHandler_Failed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zap.NewNop()
	storage := mock.NewMockStorager(mockCtrl)
	shortURL := "some url"

	type request struct {
		body string
	}
	type want struct {
		err  error
		code int
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "failed add url",
			request: request{
				body: "https://ya.ru/some",
			},
			want: want{
				err:  errors.New("some error"),
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "url already exist",
			request: request{
				body: "https://ya.ru/some",
			},
			want: want{
				err:  &data.OriginalURLAlreadyExistError{ShortURL: shortURL},
				code: http.StatusConflict,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.EXPECT().StoreShortURL(gomock.Any(), gomock.Any(), test.request.body).Times(1).Return(test.want.err)

			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.request.body))
			w := httptest.NewRecorder()
			AddHandler(logger, storage)(w, request)

			res := w.Result()
			defer closeBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)

			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
		})
	}
}

func TestAPIAddHandler_Ok(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zap.NewNop()
	storage := mock.NewMockStorager(mockCtrl)
	originalURL := "https://practicum.yandex.ru"

	type request struct {
		body string
	}
	type want struct {
		code int
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "success add url",
			request: request{
				body: `{"url": "https://practicum.yandex.ru"}`,
			},
			want: want{
				code: http.StatusCreated,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.EXPECT().StoreShortURL(gomock.Any(), gomock.Any(), originalURL).Times(1).Return(nil)

			request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(test.request.body))
			w := httptest.NewRecorder()
			APIAddHandler(logger, storage)(w, request)

			res := w.Result()
			defer closeBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.NotEmpty(t, resBody)
		})
	}
}

func TestAPIAddHandler_Failed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zap.NewNop()
	storage := mock.NewMockStorager(mockCtrl)
	originalURL := "https://practicum.yandex.ru"
	shortURL := "some url"

	type request struct {
		body string
	}
	type want struct {
		err  error
		code int
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "failed add url",
			request: request{
				body: `{"url": "https://practicum.yandex.ru"}`,
			},
			want: want{
				err:  errors.New("some error"),
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "url already exist",
			request: request{
				body: `{"url": "https://practicum.yandex.ru"}`,
			},
			want: want{
				err:  &data.OriginalURLAlreadyExistError{ShortURL: shortURL},
				code: http.StatusConflict,
			},
		},
		{
			name: "bad request",
			request: request{
				body: `sdfsdfsdfsdf`,
			},
			want: want{
				err:  nil,
				code: http.StatusBadRequest,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.want.err != nil {
				storage.EXPECT().StoreShortURL(gomock.Any(), gomock.Any(), originalURL).Times(1).Return(test.want.err)
			}

			request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(test.request.body))
			w := httptest.NewRecorder()
			APIAddHandler(logger, storage)(w, request)

			res := w.Result()
			defer closeBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			if test.want.code == http.StatusConflict {
				assert.NotEmpty(t, resBody)
			}
		})
	}
}

func TestAPIAddBatchHandler_Ok(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zap.NewNop()
	storage := mock.NewMockStorager(mockCtrl)

	type request struct {
		body string
	}
	type want struct {
		code int
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "success add url",
			request: request{
				body: `[{"correlation_id":"some_id","original_url":"https://practicum.yandex.ru"}]`,
			},
			want: want{
				code: http.StatusCreated,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.EXPECT().StoreShortURLs(gomock.Any(), gomock.Any()).Times(1).Return(nil)

			request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(test.request.body))
			newContext := context.WithValue(request.Context(), common.KeyUserID, "user_1")

			w := httptest.NewRecorder()
			APIAddBatchHandler(logger, storage)(w, request.WithContext(newContext))

			res := w.Result()
			defer closeBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.NotEmpty(t, resBody)
		})
	}
}

func TestAPIAddBatchHandler_Failed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zap.NewNop()
	storage := mock.NewMockStorager(mockCtrl)

	type request struct {
		body string
	}
	type want struct {
		err  error
		code int
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "failed add url",
			request: request{
				body: `[{"correlation_id":"some_id","original_url":"https://practicum.yandex.ru"}]`,
			},
			want: want{
				err:  errors.New("some error"),
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "bad request",
			request: request{
				body: `sdfsdfsdfsdf`,
			},
			want: want{
				err:  nil,
				code: http.StatusBadRequest,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.want.err != nil {
				storage.EXPECT().StoreShortURLs(gomock.Any(), gomock.Any()).Times(1).Return(test.want.err)
			}

			request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(test.request.body))
			newContext := context.WithValue(request.Context(), common.KeyUserID, "user_1")

			w := httptest.NewRecorder()
			APIAddBatchHandler(logger, storage)(w, request.WithContext(newContext))

			res := w.Result()
			defer closeBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			if test.want.code == http.StatusConflict {
				assert.NotEmpty(t, resBody)
			}
		})
	}
}

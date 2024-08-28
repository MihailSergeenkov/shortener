package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/MihailSergeenkov/shortener/internal/app/data"
	"github.com/MihailSergeenkov/shortener/internal/app/data/mock"
	"github.com/MihailSergeenkov/shortener/internal/app/models"
)

func TestFetchHandler_Ok(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zap.NewNop()
	storage := mock.NewMockStorager(mockCtrl)
	path := "/123"
	originalURL := "https://ya.ru/some"

	type want struct {
		code int
		url  string
	}
	tests := []struct {
		name string
		url  models.URL
		want want
	}{
		{
			name: "success fetch url",
			url: models.URL{
				ID:          1,
				ShortURL:    "123",
				OriginalURL: originalURL,
				DeletedFlag: false,
			},
			want: want{
				code: http.StatusTemporaryRedirect,
				url:  originalURL,
			},
		},
		{
			name: "when url deleted",
			url: models.URL{
				ID:          1,
				ShortURL:    "123",
				OriginalURL: originalURL,
				DeletedFlag: true,
			},
			want: want{
				code: http.StatusGone,
				url:  "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.EXPECT().GetURL(gomock.Any(), "123").Times(1).Return(test.url, nil)

			request := httptest.NewRequest(http.MethodGet, path, http.NoBody)
			w := httptest.NewRecorder()
			FetchHandler(logger, storage)(w, request)

			res := w.Result()
			defer closeBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)

			if test.want.code == http.StatusTemporaryRedirect {
				assert.Equal(t, test.want.url, res.Header.Get("Location"))
			}
		})
	}
}

func TestFetchHandler_Failed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zap.NewNop()
	storage := mock.NewMockStorager(mockCtrl)
	path := "/123"

	type want struct {
		code int
	}
	tests := []struct {
		name string
		err  error
		want want
	}{
		{
			name: "failed fetch url",
			err:  errors.New("some error"),
			want: want{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "when url not found",
			err:  data.ErrURLNotFound,
			want: want{
				code: http.StatusNotFound,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.EXPECT().GetURL(gomock.Any(), "123").Times(1).Return(models.URL{}, test.err)

			request := httptest.NewRequest(http.MethodGet, path, http.NoBody)
			w := httptest.NewRecorder()
			FetchHandler(logger, storage)(w, request)

			res := w.Result()
			defer closeBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}

func TestAPIFetchUserURLsHandler_Ok(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zap.NewNop()
	storage := mock.NewMockStorager(mockCtrl)

	type want struct {
		code int
	}
	tests := []struct {
		name string
		urls []models.URL
		want want
	}{
		{
			name: "success fetch urls",
			urls: []models.URL{
				{
					ShortURL:    "some_url",
					OriginalURL: "some_url",
					UserID:      "some_id",
					ID:          1,
					DeletedFlag: false,
				},
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "when urls not found",
			urls: []models.URL{},
			want: want{
				code: http.StatusNoContent,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.EXPECT().FetchUserURLs(gomock.Any()).Times(1).Return(test.urls, nil)

			request := httptest.NewRequest(http.MethodGet, "/api/user/urls", http.NoBody)
			w := httptest.NewRecorder()
			APIFetchUserURLsHandler(logger, storage)(w, request)

			res := w.Result()
			defer closeBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)

			resBody, err := io.ReadAll(res.Body)

			if test.want.code == http.StatusOK {
				require.NoError(t, err)
				assert.NotEmpty(t, resBody)
			}
		})
	}
}

func TestAPIFetchUserURLsHandler_Failed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zap.NewNop()
	storage := mock.NewMockStorager(mockCtrl)
	t.Run("when fetch failed", func(t *testing.T) {
		errSome := errors.New("some error")
		storage.EXPECT().FetchUserURLs(gomock.Any()).Times(1).Return([]models.URL{}, errSome)

		request := httptest.NewRequest(http.MethodGet, "/api/user/urls", http.NoBody)
		w := httptest.NewRecorder()
		APIFetchUserURLsHandler(logger, storage)(w, request)

		res := w.Result()
		defer closeBody(t, res)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}

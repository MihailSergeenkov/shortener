package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MihailSergeenkov/shortener/internal/app/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddHandler(t *testing.T) {
	type request struct {
		method string
		body   string
	}

	type want struct {
		code int
	}
	tests := []struct {
		name    string
		storage data.Storage
		request request
		want    want
	}{
		{
			name: "success add url",
			storage: data.Storage{
				FileStoragePath: "some/path",
				URLs: map[string]data.URL{
					"123": {
						ID:          1,
						ShortURL:    "123",
						OriginalURL: "https://ya.ru/main",
					},
				},
			},
			request: request{
				method: http.MethodPost,
				body:   "https://ya.ru/some",
			},
			want: want{
				code: http.StatusCreated,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.request.method, "/", http.NoBody)
			w := httptest.NewRecorder()
			AddHandler(test.storage)(w, request)

			res := w.Result()
			defer closeBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)

			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.NotEmpty(t, resBody)
		})
	}
}

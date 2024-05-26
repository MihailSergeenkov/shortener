package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MihailSergeenkov/shortener/internal/app/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIAddHandler(t *testing.T) {
	storage := data.Storage{
		FileStoragePath: "some/path",
		URLs: map[string]data.URL{
			"123": {
				ID:          1,
				ShortURL:    "123",
				OriginalURL: "https://ya.ru/main",
			},
		},
	}

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
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(test.request.body))
			w := httptest.NewRecorder()
			APIAddHandler(storage)(w, request)

			res := w.Result()
			defer func() {
				err := res.Body.Close()

				if err != nil {
					t.Log(err)
				}
			}()

			assert.Equal(t, test.want.code, res.StatusCode)

			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.NotEmpty(t, resBody)
		})
	}
}

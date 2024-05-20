package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MihailSergeenkov/shortener/internal/app/storage"
	"github.com/MihailSergeenkov/shortener/internal/app/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApiAddHandler(t *testing.T) {
	urls := storage.Urls{"123": "https://ya.ru/main"}

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
			ApiAddHandler(urls)(w, request)

			res := w.Result()
			defer test_helpers.CloseBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)

			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.NotEmpty(t, resBody)
		})
	}
}

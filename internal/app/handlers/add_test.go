package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MihailSergeenkov/shortener/internal/app/storage"
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
		urls    storage.Urls
		request request
		want    want
	}{
		{
			name: "success add url",
			urls: storage.Urls{
				"123": "https://ya.ru/main",
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
			request := httptest.NewRequest(test.request.method, "/", nil)
			w := httptest.NewRecorder()
			AddHandler(test.urls)(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)

			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.NotEmpty(t, resBody)
		})
	}
}

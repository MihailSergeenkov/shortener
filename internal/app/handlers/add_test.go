package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAddHandler(t *testing.T) {
	logger := zap.NewNop()
	storage := MockStorage{}

	type request struct {
		method string
		body   string
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
			AddHandler(logger, &storage)(w, request)

			res := w.Result()
			defer closeBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)

			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.NotEmpty(t, resBody)
		})
	}
}

func TestAPIAddHandler(t *testing.T) {
	logger := zap.NewNop()
	storage := MockStorage{}

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
			APIAddHandler(logger, &storage)(w, request)

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

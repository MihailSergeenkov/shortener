package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/MihailSergeenkov/shortener/internal/app/models"
)

func TestFetchHandler(t *testing.T) {
	logger := zap.NewNop()
	storage := MockStorage{
		urls: map[string]models.URL{
			"123": {
				ID:          1,
				ShortURL:    "123",
				OriginalURL: originalURL,
			},
		},
	}

	type request struct {
		method string
		path   string
	}

	type want struct {
		code int
		url  string
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "success fetch url",
			request: request{
				method: http.MethodGet,
				path:   "/123",
			},
			want: want{
				code: http.StatusTemporaryRedirect,
				url:  originalURL,
			},
		},
		{
			name: "when url not found",
			request: request{
				method: http.MethodGet,
				path:   "/12345678",
			},
			want: want{
				code: http.StatusNotFound,
				url:  "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.request.method, test.request.path, http.NoBody)
			w := httptest.NewRecorder()
			FetchHandler(logger, &storage)(w, request)

			res := w.Result()
			defer closeBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)

			if test.want.code == http.StatusTemporaryRedirect {
				assert.Equal(t, test.want.url, res.Header.Get("Location"))
			}
		})
	}
}

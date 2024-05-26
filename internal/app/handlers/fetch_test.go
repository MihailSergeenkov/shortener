package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MihailSergeenkov/shortener/internal/app/data"
	"github.com/stretchr/testify/assert"
)

func TestFetchHandler(t *testing.T) {
	url := "https://ya.ru/some"

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
		storage data.Storage
		request request
		want    want
	}{
		{
			name: "success fetch url",
			storage: data.Storage{
				FileStoragePath: "some/path",
				Urls: map[string]data.Url{
					"123": {
						ID:          1,
						ShortUrl:    "123",
						OriginalUrl: url,
					},
				},
			},
			request: request{
				method: http.MethodGet,
				path:   "/123",
			},
			want: want{
				code: http.StatusTemporaryRedirect,
				url:  url,
			},
		},
		{
			name: "when url not found",
			storage: data.Storage{
				FileStoragePath: "some/path",
				Urls: map[string]data.Url{
					"123": {
						ID:          1,
						ShortUrl:    "123",
						OriginalUrl: url,
					},
				},
			},
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
			FetchHandler(test.storage)(w, request)

			res := w.Result()
			defer closeBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)

			if test.want.code == http.StatusTemporaryRedirect {
				assert.Equal(t, test.want.url, res.Header.Get("Location"))
			}
		})
	}
}

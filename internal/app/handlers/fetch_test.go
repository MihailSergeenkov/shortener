package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MihailSergeenkov/shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
)

func TestFetchHandler(t *testing.T) {
	url := "https://ya.ru/some"
	h, _ := storage.AddUrl(url)

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
				path:   "/" + h,
			},
			want: want{
				code: 307,
				url:  url,
			},
		},
		{
			name: "when does not Get method",
			request: request{
				method: http.MethodPost,
				path:   "/" + h,
			},
			want: want{
				code: 400,
				url:  "",
			},
		},
		{
			name: "when url not found",
			request: request{
				method: http.MethodGet,
				path:   "/12345678",
			},
			want: want{
				code: 404,
				url:  "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.request.method, test.request.path, nil)
			w := httptest.NewRecorder()
			FetchHandler(w, request)

			res := w.Result()

			assert.Equal(t, test.want.code, res.StatusCode)

			if test.want.code == 307 {
				assert.Equal(t, test.want.url, res.Header.Get("Location"))
			}
		})
	}
}

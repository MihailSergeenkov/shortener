package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/MihailSergeenkov/shortener/internal/app/models"
	"go.uber.org/zap"
)

func ExamplePingHandler() {
	logger := zap.NewNop()
	storage := MockStorage{}

	request := httptest.NewRequest(http.MethodGet, "/ping", http.NoBody)
	w := httptest.NewRecorder()
	PingHandler(logger, &storage)(w, request)

	res := w.Result()
	closeExampleBody(res)

	fmt.Println(res.StatusCode)

	// Output:
	// 200
}

func ExampleAddHandler() {
	logger := zap.NewNop()
	storage := MockStorage{}

	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://ya.ru/some"))
	w := httptest.NewRecorder()
	AddHandler(logger, &storage)(w, request)

	res := w.Result()
	defer closeExampleBody(res)

	fmt.Println(res.StatusCode)

	// Output:
	// 201
}

func ExampleAPIAddHandler() {
	logger := zap.NewNop()
	storage := MockStorage{}

	request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader("https://ya.ru/some"))
	w := httptest.NewRecorder()
	AddHandler(logger, &storage)(w, request)

	res := w.Result()
	defer closeExampleBody(res)

	fmt.Println(res.StatusCode)

	// Output:
	// 201
}

func ExampleAPIDeleteUserURLsHandler() {
	logger := zap.NewNop()
	storage := MockStorage{}
	body := `["6qxTVvsy", "RTfd56hn", "Jlfd67ds"]`
	request := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(body))
	w := httptest.NewRecorder()
	APIDeleteUserURLsHandler(logger, &storage)(w, request)

	res := w.Result()
	defer closeExampleBody(res)

	fmt.Println(res.StatusCode)

	// Output:
	// 202
}

func ExampleFetchHandler() {
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

	request := httptest.NewRequest(http.MethodGet, "/123", http.NoBody)
	w := httptest.NewRecorder()
	FetchHandler(logger, &storage)(w, request)

	res := w.Result()
	defer closeExampleBody(res)

	fmt.Println(res.StatusCode)

	// Output:
	// 307
}

func closeExampleBody(r *http.Response) {
	if err := r.Body.Close(); err != nil {
		fmt.Print(err)
	}
}

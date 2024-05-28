package handlers

import (
	"net/http"
	"testing"

	"github.com/MihailSergeenkov/shortener/internal/app/data"
	"github.com/MihailSergeenkov/shortener/internal/app/models"
)

type MockStorage struct {
	urls map[string]models.URL
}

var shortURL = "123"
var originalURL = "https://ya.ru/some"

func (s *MockStorage) FetchURL(shortURL string) (models.URL, error) {
	u, present := s.urls[shortURL]

	if !present {
		return models.URL{}, data.ErrURLNotFound
	}

	return u, nil
}

func (s *MockStorage) AddURL(originalURL string) (models.URL, error) {
	return models.URL{
		ID:          1,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}, nil
}

func closeBody(t *testing.T, r *http.Response) {
	t.Helper()
	err := r.Body.Close()

	if err != nil {
		t.Log(err)
	}
}

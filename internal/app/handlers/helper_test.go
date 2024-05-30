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

var originalURL = "https://ya.ru/some"

func (s *MockStorage) GetOriginalURL(shortURL string) (string, error) {
	u, present := s.urls[shortURL]

	if !present {
		return "", data.ErrURLNotFound
	}

	return u.OriginalURL, nil
}

func (s *MockStorage) StoreShortURL(shortURL string, originalURL string) error {
	return nil
}

func closeBody(t *testing.T, r *http.Response) {
	t.Helper()
	err := r.Body.Close()

	if err != nil {
		t.Log(err)
	}
}

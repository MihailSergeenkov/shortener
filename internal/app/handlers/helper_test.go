package handlers

import (
	"context"
	"net/http"
	"testing"

	"github.com/MihailSergeenkov/shortener/internal/app/data"
	"github.com/MihailSergeenkov/shortener/internal/app/models"
)

type MockStorage struct {
	urls map[string]models.URL
}

func (s *MockStorage) GetURL(_ context.Context, shortURL string) (models.URL, error) {
	u, present := s.urls[shortURL]

	if !present {
		return models.URL{}, data.ErrURLNotFound
	}

	return u, nil
}

func (s *MockStorage) StoreShortURL(_ context.Context, shortURL string, originalURL string) error {
	return nil
}

func (s *MockStorage) StoreShortURLs(_ context.Context, urls []models.URL) error {
	return nil
}

func (s *MockStorage) FetchUserURLs(_ context.Context) ([]models.URL, error) {
	return []models.URL{}, nil
}

func (s *MockStorage) DeleteShortURLs(ctx context.Context, urls []string) error {
	return nil
}

func (s *MockStorage) DropDeletedURLs(_ context.Context) error {
	return nil
}

func (s *MockStorage) Ping(_ context.Context) error {
	return nil
}

func (s *MockStorage) Close() error {
	return nil
}

func closeBody(t *testing.T, r *http.Response) {
	t.Helper()
	err := r.Body.Close()

	if err != nil {
		t.Log(err)
	}
}

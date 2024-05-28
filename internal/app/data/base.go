package data

import (
	"errors"
	"fmt"

	"github.com/MihailSergeenkov/shortener/internal/app/models"
	"github.com/MihailSergeenkov/shortener/internal/app/services"
)

const (
	initSize int = 100
	maxRetry int = 5
)

var (
	ErrURLNotFound = errors.New("url not found")
	ErrMaxRetry    = errors.New("generation attempts exceeded")
)

type BaseStorage struct {
	urls map[string]models.URL
}

func NewBaseStorage() *BaseStorage {
	return &BaseStorage{
		urls: make(map[string]models.URL, initSize),
	}
}

func (s *BaseStorage) AddURL(originalURL string) (models.URL, error) {
	for range maxRetry {
		shortURL, err := services.GenerateShortURL()
		if err != nil {
			return models.URL{}, fmt.Errorf("failed to generate short URL: %w", err)
		}

		if _, ok := s.urls[shortURL]; ok {
			continue
		}

		url := models.URL{
			ID:          uint(len(s.urls) + 1),
			ShortURL:    shortURL,
			OriginalURL: originalURL,
		}

		s.urls[shortURL] = url

		return url, nil
	}

	return models.URL{}, fmt.Errorf("%w for original URL %s", ErrMaxRetry, originalURL)
}

func (s *BaseStorage) FetchURL(shortURL string) (models.URL, error) {
	u, ok := s.urls[shortURL]

	if !ok {
		return models.URL{}, fmt.Errorf("%w for short URL %s", ErrURLNotFound, shortURL)
	}

	return u, nil
}

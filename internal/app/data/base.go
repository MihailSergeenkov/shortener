package data

import (
	"context"
	"fmt"

	"github.com/MihailSergeenkov/shortener/internal/app/models"
)

const initSize int = 100

type BaseStorage struct {
	urls map[string]models.URL
}

func NewBaseStorage() *BaseStorage {
	return &BaseStorage{
		urls: make(map[string]models.URL, initSize),
	}
}

func (s *BaseStorage) StoreShortURL(_ context.Context, shortURL string, originalURL string) error {
	if _, ok := s.urls[shortURL]; ok {
		return ErrShortURLAlreadyExist
	}

	url := models.URL{
		ID:          uint(len(s.urls) + 1),
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}

	s.urls[shortURL] = url

	return nil
}

func (s *BaseStorage) StoreShortURLs(_ context.Context, urls []models.URL) error {
	for _, url := range urls {
		if _, ok := s.urls[url.ShortURL]; ok {
			return ErrShortURLAlreadyExist
		}
	}

	lastID := len(s.urls)

	for i, url := range urls {
		url.ID = uint(lastID + i)
		s.urls[url.ShortURL] = url
	}

	return nil
}

func (s *BaseStorage) GetOriginalURL(_ context.Context, shortURL string) (string, error) {
	u, ok := s.urls[shortURL]

	if !ok {
		return "", fmt.Errorf("%w for short URL %s", ErrURLNotFound, shortURL)
	}

	return u.OriginalURL, nil
}

func (s *BaseStorage) Ping(_ context.Context) error {
	return nil
}

func (s *BaseStorage) Close() error {
	return nil
}

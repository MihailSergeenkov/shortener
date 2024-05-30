package data

import (
	"errors"
	"fmt"

	"github.com/MihailSergeenkov/shortener/internal/app/models"
)

const initSize int = 100

var (
	ErrURLNotFound          = errors.New("url not found")
	ErrShortURLAlreadyExist = errors.New("short url already exist")
)

type BaseStorage struct {
	urls map[string]models.URL
}

func NewBaseStorage() *BaseStorage {
	return &BaseStorage{
		urls: make(map[string]models.URL, initSize),
	}
}

func (s *BaseStorage) StoreShortURL(shortURL string, originalURL string) error {
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

func (s *BaseStorage) GetOriginalURL(shortURL string) (string, error) {
	u, ok := s.urls[shortURL]

	if !ok {
		return "", fmt.Errorf("%w for short URL %s", ErrURLNotFound, shortURL)
	}

	return u.OriginalURL, nil
}

package data

import (
	"context"
	"fmt"

	"github.com/MihailSergeenkov/shortener/internal/app/constants"
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

func (s *BaseStorage) StoreShortURL(ctx context.Context, shortURL string, originalURL string) error {
	if _, ok := s.urls[shortURL]; ok {
		return ErrShortURLAlreadyExist
	}

	userID, ok := ctx.Value(constants.KeyUserID).(string)

	if !ok {
		return fmt.Errorf("failed to fetch user id from context")
	}

	url := models.URL{
		ID:          uint(len(s.urls) + 1),
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		UserID:      userID,
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

func (s *BaseStorage) FetchUserURLs(ctx context.Context) ([]models.URL, error) {
	urls := []models.URL{}
	userID := ctx.Value(constants.KeyUserID)

	for _, u := range s.urls {
		if u.UserID == userID {
			urls = append(urls, u)
		}
	}

	return urls, nil
}

func (s *BaseStorage) Ping(_ context.Context) error {
	return nil
}

func (s *BaseStorage) Close() error {
	return nil
}

package data

import (
	"context"
	"fmt"

	"github.com/MihailSergeenkov/shortener/internal/app/common"
	"github.com/MihailSergeenkov/shortener/internal/app/models"
)

const initSize int = 100

// BaseStorage структура in-memory БД.
type BaseStorage struct {
	urls map[string]models.URL
}

// NewBaseStorage инициализирует in-memory БД.
func NewBaseStorage() *BaseStorage {
	return &BaseStorage{
		urls: make(map[string]models.URL, initSize),
	}
}

// StoreShortURL сохраняет короткую ссылку.
func (s *BaseStorage) StoreShortURL(ctx context.Context, shortURL string, originalURL string) error {
	if _, ok := s.urls[shortURL]; ok {
		return ErrShortURLAlreadyExist
	}

	userID, ok := ctx.Value(common.KeyUserID).(string)

	if !ok {
		return common.ErrFetchUserIDFromContext
	}

	url := models.URL{
		ID:          uint(len(s.urls) + 1),
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		UserID:      userID,
		DeletedFlag: false,
	}

	s.urls[shortURL] = url

	return nil
}

// StoreShortURLs сохраняет несколько коротких ссылок.
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

// GetURL получает оригинальную ссылку по короткой.
func (s *BaseStorage) GetURL(_ context.Context, shortURL string) (models.URL, error) {
	u, ok := s.urls[shortURL]

	if !ok {
		return models.URL{}, fmt.Errorf("%w for short URL %s", ErrURLNotFound, shortURL)
	}

	return u, nil
}

// FetchUserURLs получает все пользовательские ссылки.
func (s *BaseStorage) FetchUserURLs(ctx context.Context) ([]models.URL, error) {
	urls := []models.URL{}
	userID := ctx.Value(common.KeyUserID)

	for _, u := range s.urls {
		if u.UserID == userID {
			urls = append(urls, u)
		}
	}

	return urls, nil
}

// DeleteShortURLs мягко удаляет ссылки.
func (s *BaseStorage) DeleteShortURLs(ctx context.Context, urls []string) error {
	for _, url := range urls {
		u, ok := s.urls[url]
		if !ok {
			return fmt.Errorf("%w for short URL %s", ErrURLNotFound, url)
		}

		u.DeletedFlag = true
		s.urls[url] = u
	}

	return nil
}

// DropDeletedURLs очищает из БД удаленные ссылки (не используется для in-memory БД).
func (s *BaseStorage) DropDeletedURLs(_ context.Context) error {
	return nil
}

// FetchStats получает статистические данные.
func (s *BaseStorage) FetchStats(ctx context.Context) (int, int, error) {
	users := map[string]struct{}{}

	for _, u := range s.urls {
		users[u.UserID] = struct{}{}
	}

	return len(s.urls), len(users), nil
}

// Ping проверяет работоспособность БД (не используется для in-memory БД).
func (s *BaseStorage) Ping(_ context.Context) error {
	return nil
}

// Close закрывает соединение с БД (не используется для in-memory БД).
func (s *BaseStorage) Close() error {
	return nil
}

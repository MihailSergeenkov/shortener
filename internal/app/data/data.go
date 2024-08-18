// Пакет data предназначен для подключения БД к сервису.
package data

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/MihailSergeenkov/shortener/internal/app/config"
	"github.com/MihailSergeenkov/shortener/internal/app/models"
)

// Ошибки БД
var (
	ErrURLNotFound          = errors.New("url not found")           // короткая ссылка не найдена
	ErrShortURLAlreadyExist = errors.New("short url already exist") // короткая ссылка уже существует в сервисе
)

// OriginalURLAlreadyExistError структура ошибки, когда оригинальная ссылка уже существует в сервисе.
type OriginalURLAlreadyExistError struct {
	ShortURL string
}

// Error возвращает описание ошибки.
func (e *OriginalURLAlreadyExistError) Error() string {
	return "original url already exist"
}

func newOriginalURLAlreadyExistError(url string) error {
	return &OriginalURLAlreadyExistError{
		ShortURL: url,
	}
}

// Storager интерфейс к БД
type Storager interface {
	StoreShortURL(ctx context.Context, shortURL string, originalURL string) error // сохранение короткой ссылки
	StoreShortURLs(ctx context.Context, urls []models.URL) error                  // сохранение нескольких коротких ссылок
	GetURL(ctx context.Context, shortURL string) (models.URL, error)              // получение оригинальной ссылки по короткой
	FetchUserURLs(ctx context.Context) ([]models.URL, error)                      // получить все сохраненные ссылки пользователя
	DeleteShortURLs(ctx context.Context, urls []string) error                     // мягко удалить ссылки
	DropDeletedURLs(ctx context.Context) error                                    // очистить из БД удаленные ссылки
	Ping(ctx context.Context) error                                               // проверка работоспособности БД
	Close() error                                                                 // закрыть соединение с БД
}

// NewStorage инициализирует БД.
func NewStorage(ctx context.Context, logger *zap.Logger, params *config.Settings) (Storager, error) {
	dbDSN := params.DatabaseDSN
	fsp := params.FileStoragePath

	if dbDSN != "" {
		return NewDBStorage(ctx, logger, dbDSN)
	}

	if fsp == "" {
		return NewBaseStorage(), nil
	}

	return NewFileStorage(logger, fsp)
}

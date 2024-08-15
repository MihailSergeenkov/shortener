package data

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/MihailSergeenkov/shortener/internal/app/config"
	"github.com/MihailSergeenkov/shortener/internal/app/models"
)

var (
	ErrURLNotFound          = errors.New("url not found")
	ErrShortURLAlreadyExist = errors.New("short url already exist")
)

type OriginalURLAlreadyExistError struct {
	ShortURL string
}

func (e *OriginalURLAlreadyExistError) Error() string {
	return "original url already exist"
}

func newOriginalURLAlreadyExistError(url string) error {
	return &OriginalURLAlreadyExistError{
		ShortURL: url,
	}
}

type Storager interface {
	StoreShortURL(ctx context.Context, shortURL string, originalURL string) error
	StoreShortURLs(ctx context.Context, urls []models.URL) error
	GetURL(ctx context.Context, shortURL string) (models.URL, error)
	FetchUserURLs(ctx context.Context) ([]models.URL, error)
	DeleteShortURLs(ctx context.Context, urls []string) error
	DropDeletedURLs(ctx context.Context) error
	Ping(ctx context.Context) error
	Close() error
}

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

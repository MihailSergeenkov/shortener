package data

import (
	"context"
	"errors"
	"fmt"

	"github.com/MihailSergeenkov/shortener/internal/app/config"
	"github.com/MihailSergeenkov/shortener/internal/app/models"
	"go.uber.org/zap"
)

var (
	ErrURLNotFound          = errors.New("url not found")
	ErrShortURLAlreadyExist = errors.New("short url already exist")
)

type ErrOriginalURLAlreadyExist struct {
	ShortURL string
}

func (e *ErrOriginalURLAlreadyExist) Error() string {
	return fmt.Sprintf("original url already exist")
}

func newErrOriginalURLAlreadyExist(url string) error {
	return &ErrOriginalURLAlreadyExist{
		ShortURL: url,
	}
}

type Storager interface {
	StoreShortURL(ctx context.Context, shortURL string, originalURL string) error
	StoreShortURLs(ctx context.Context, urls []models.URL) error
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)
	Close() error
}

func NewStorage(logger *zap.Logger, params config.Settings) (Storager, error) {
	dbDSN := params.DatabaseDSN
	fsp := params.FileStoragePath

	if dbDSN != "" {
		return NewDBStorage(logger, dbDSN)
	}

	if fsp == "" {
		return NewBaseStorage(), nil
	}

	return NewFileStorage(logger, fsp)
}

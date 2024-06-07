package data

import (
	"context"
	"errors"

	"github.com/MihailSergeenkov/shortener/internal/app/config"
	"github.com/MihailSergeenkov/shortener/internal/app/models"
	"go.uber.org/zap"
)

var (
	ErrURLNotFound          = errors.New("url not found")
	ErrShortURLAlreadyExist = errors.New("short url already exist")
)

type Storager interface {
	StoreShortURL(ctx context.Context, shortURL string, originalURL string) error
	StoreShortURLs(ctx context.Context, URLs []models.URL) error
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

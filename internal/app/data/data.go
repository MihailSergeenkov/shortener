package data

import (
	"go.uber.org/zap"
)

type Storager interface {
	StoreShortURL(shortURL string, originalURL string) error
	GetOriginalURL(shortURL string) (string, error)
	Close() error
}

func NewStorage(logger *zap.Logger, fsp, dbDSN string) (Storager, error) {
	if dbDSN != "" {
		return NewDbStorage(logger, dbDSN)
	}

	if fsp == "" {
		return NewBaseStorage(), nil
	}

	return NewFileStorage(logger, fsp)
}

package data

import (
	"go.uber.org/zap"
)

type Storager interface {
	StoreShortURL(shortURL string, originalURL string) error
	GetOriginalURL(shortURL string) (string, error)
}

func NewStorage(logger *zap.Logger, fsp string) (Storager, error) {
	if fsp == "" {
		return NewBaseStorage(), nil
	}

	return NewFileStorage(logger, fsp)
}

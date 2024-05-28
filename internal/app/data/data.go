package data

import (
	"github.com/MihailSergeenkov/shortener/internal/app/models"
	"go.uber.org/zap"
)

type Storager interface {
	AddURL(originalURL string) (models.URL, error)
	FetchURL(shortURL string) (models.URL, error)
}

func NewStorage(logger *zap.Logger, fsp string) (Storager, error) {
	if fsp == "" {
		return NewBaseStorage(), nil
	}

	return NewFileStorage(logger, fsp)
}

package data

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type DBStorage struct {
	*sql.DB
	logger *zap.Logger
}

func NewDBStorage(logger *zap.Logger, dbDSN string) (*DBStorage, error) {
	db, err := sql.Open("pgx", dbDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	return &DBStorage{
		logger: logger,
		DB:     db,
	}, nil
}

func (s *DBStorage) StoreShortURL(shortURL string, originalURL string) error {
	return nil
}

func (s *DBStorage) GetOriginalURL(shortURL string) (string, error) {
	return "", nil
}

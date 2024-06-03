package data

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type DbStorage struct {
	*sql.DB
	logger *zap.Logger
}

func NewDbStorage(logger *zap.Logger, dbDSN string) (*DbStorage, error) {
	db, err := sql.Open("pgx", dbDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}
	// defer db.Close()

	return &DbStorage{
		logger: logger,
		DB:     db,
	}, nil
}

func (s *DbStorage) StoreShortURL(shortURL string, originalURL string) error {
	return nil
}

func (s *DbStorage) GetOriginalURL(shortURL string) (string, error) {
	return "", nil
}

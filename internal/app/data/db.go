package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/MihailSergeenkov/shortener/internal/app/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type DBStorage struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

const stmt = `INSERT INTO urls (short_url, original_url) VALUES ($1, $2)`

func NewDBStorage(ctx context.Context, logger *zap.Logger, dbDSN string) (*DBStorage, error) {
	pool, err := initPool(ctx, logger, dbDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a connection pool: %w", err)
	}

	s := &DBStorage{
		logger: logger,
		pool:   pool,
	}

	if err := s.initDB(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize a DB: %w", err)
	}

	return s, nil
}

func (s *DBStorage) StoreShortURL(ctx context.Context, shortURL string, originalURL string) error {
	tx, txErr := s.pool.Begin(ctx)
	if txErr != nil {
		return fmt.Errorf("failed to start transaction: %w", txErr)
	}

	url, exist, getErr := getShortURLForOriginalURL(ctx, tx, originalURL)
	if getErr != nil {
		rollbackTx(ctx, tx, s.logger)
		return fmt.Errorf("failed to check present original URL: %w", getErr)
	}
	if exist {
		rollbackTx(ctx, tx, s.logger)
		return newOriginalURLAlreadyExistError(url)
	}

	_, err := tx.Exec(ctx, stmt, shortURL, originalURL)
	if err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}

	cErr := tx.Commit(ctx)
	if cErr != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func getShortURLForOriginalURL(ctx context.Context, tx pgx.Tx, originalURL string) (string, bool, error) {
	const queryStmt = `SELECT id, short_url, original_url
		FROM urls
		WHERE original_url = $1
		LIMIT 1`

	row := tx.QueryRow(ctx, queryStmt, originalURL)

	var u models.URL
	err := row.Scan(&u.ID, &u.ShortURL, &u.OriginalURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", false, nil
		}

		return "", false, fmt.Errorf("failed to scan a response row: %w", err)
	}

	return u.ShortURL, true, nil
}

func (s *DBStorage) StoreShortURLs(ctx context.Context, urls []models.URL) error {
	batch := &pgx.Batch{}

	for _, url := range urls {
		batch.Queue(stmt, url.ShortURL, url.OriginalURL)
	}

	result := s.pool.SendBatch(ctx, batch)
	defer func() {
		if err := result.Close(); err != nil {
			s.logger.Error("failed to close batch result", zap.Error(err))
		}
	}()

	_, err := result.Exec()
	if err != nil {
		return fmt.Errorf("unable to insert batch: %w", err)
	}

	return result.Close()
}

func (s *DBStorage) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	const queryStmt = `SELECT id, short_url, original_url
		FROM urls
		WHERE short_url = $1
		LIMIT 1`

	row := s.pool.QueryRow(ctx, queryStmt, shortURL)

	var u models.URL
	err := row.Scan(&u.ID, &u.ShortURL, &u.OriginalURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%w for short URL %s", ErrURLNotFound, shortURL)
		}

		return "", fmt.Errorf("failed to scan a response row: %w", err)
	}

	return u.OriginalURL, nil
}

func (s *DBStorage) Ping(ctx context.Context) error {
	return s.pool.Ping(ctx)
}

func (s *DBStorage) Close() error {
	s.pool.Close()
	return nil
}

func (s *DBStorage) initDB(ctx context.Context) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	if err := s.createSchema(ctx, tx); err != nil {
		return fmt.Errorf("failed to create the DB schema: %w", err)
	}

	cErr := tx.Commit(ctx)
	if cErr != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *DBStorage) createSchema(ctx context.Context, tx pgx.Tx) error {
	createSchemaStmts := []string{
		`CREATE TABLE IF NOT EXISTS urls(
			id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			short_url VARCHAR(200) NOT NULL,
			original_url VARCHAR(300) NOT NULL
		)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS short_url_index ON urls(short_url)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS original_url_index ON urls(original_url)`,
	}

	for _, stmt := range createSchemaStmts {
		_, err := tx.Exec(ctx, stmt)

		if err != nil {
			rollbackTx(ctx, tx, s.logger)
			return fmt.Errorf("failed to exec transaction: %w", err)
		}
	}

	return nil
}

func initPool(ctx context.Context, logger *zap.Logger, dbDSN string) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(dbDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the DSN: %w", err)
	}

	poolCfg.ConnConfig.Tracer = &queryTracer{logger: logger}
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping the DB: %w", err)
	}

	return pool, nil
}

func rollbackTx(ctx context.Context, tx pgx.Tx, logger *zap.Logger) {
	if rErr := tx.Rollback(ctx); rErr != nil {
		if !errors.Is(rErr, sql.ErrTxDone) {
			logger.Error("failed to rollback the transaction", zap.Error(rErr))
		}
	}
}

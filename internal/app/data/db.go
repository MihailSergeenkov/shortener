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
	Pool   *pgxpool.Pool
	logger *zap.Logger
}

const stmt = `INSERT INTO urls (short_url, original_url) VALUES ($1, $2)`

func NewDBStorage(logger *zap.Logger, dbDSN string) (*DBStorage, error) {
	ctx := context.Background()

	if err := initDB(ctx, logger, dbDSN); err != nil {
		return nil, fmt.Errorf("failed to initialize a DB: %w", err)
	}

	pool, err := initPool(ctx, logger, dbDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a connection pool: %w", err)
	}

	return &DBStorage{
		logger: logger,
		Pool:   pool,
	}, nil
}

func (s *DBStorage) StoreShortURL(ctx context.Context, shortURL string, originalURL string) error {
	_, err := s.Pool.Exec(ctx, stmt, shortURL, originalURL)
	if err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}

	return nil
}

func (s *DBStorage) StoreShortURLs(ctx context.Context, urls []models.URL) error {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	for _, url := range urls {
		_, err := tx.Exec(ctx, stmt, url.ShortURL, url.OriginalURL)

		if err != nil {
			rErr := tx.Rollback(ctx)
			if rErr != nil {
				s.logger.Error("failed to rollback the transaction", zap.Error(err))
			}
			return fmt.Errorf("failed to exec transaction: %w", err)
		}
	}

	cErr := tx.Commit(ctx)
	if cErr != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *DBStorage) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	const queryStmt = `SELECT id, short_url, original_url
		FROM urls
		WHERE short_url = $1
		LIMIT 1`

	row := s.Pool.QueryRow(ctx, queryStmt, shortURL)

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

func (s *DBStorage) Close() error {
	s.Pool.Close()
	return nil
}

func initDB(ctx context.Context, logger *zap.Logger, dbDSN string) error {
	db, err := sql.Open("pgx", dbDSN)
	if err != nil {
		return fmt.Errorf("failed to open db connection: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("failed to properly close the DB connection", zap.Error(err))
		}
	}()

	if err := createSchema(ctx, logger, db); err != nil {
		return fmt.Errorf("failed to create the DB schema: %w", err)
	}

	return nil
}

func createSchema(ctx context.Context, logger *zap.Logger, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start a transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			if !errors.Is(err, sql.ErrTxDone) {
				logger.Error("failed to rollback the transaction", zap.Error(err))
			}
		}
	}()

	createSchemaStmts := []string{
		`CREATE TABLE IF NOT EXISTS urls(
			id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			short_url VARCHAR(200) NOT NULL,
			original_url VARCHAR(300) NOT NULL
		)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS short_url_index ON urls(short_url)`,
	}

	for _, stmt := range createSchemaStmts {
		if _, err := tx.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("failed to execute statement `%s`: %w", stmt, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit the transaction: %w", err)
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

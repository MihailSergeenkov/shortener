package data

import (
	"context"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/MihailSergeenkov/shortener/internal/app/common"
	"github.com/MihailSergeenkov/shortener/internal/app/models"
)

const stmt = `
	WITH new_url AS (
		INSERT INTO urls (short_url, original_url, user_id) 
		VALUES ($1, $2, $3)
		ON CONFLICT (original_url) WHERE is_deleted = false DO NOTHING
		RETURNING short_url
	)
	SELECT short_url, true as is_new FROM new_url
	UNION
	SELECT short_url, false as is_new FROM urls WHERE original_url = $2 AND is_deleted = false
`

// DBPooler интерфейс к пулу БД.
type DBPooler interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Ping(ctx context.Context) error
	Close()
}

// DBStorage структура postgresql БД.
type DBStorage struct {
	pool   DBPooler
	logger *zap.Logger
}

// NewDBStorage инициализирует postgresql БД.
func NewDBStorage(ctx context.Context, logger *zap.Logger, dbDSN string) (*DBStorage, error) {
	if err := runMigrations(dbDSN); err != nil {
		return nil, fmt.Errorf("failed to run DB migrations: %w", err)
	}

	pool, err := initPool(ctx, logger, dbDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a connection pool: %w", err)
	}

	s := &DBStorage{
		logger: logger,
		pool:   pool,
	}

	return s, nil
}

//go:embed migrations/*.sql
var migrationsDir embed.FS

func runMigrations(dsn string) error {
	d, err := iofs.New(migrationsDir, "migrations")
	if err != nil {
		return fmt.Errorf("failed to return an iofs driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return fmt.Errorf("failed to get a new migrate instance: %w", err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations to the DB: %w", err)
		}
	}
	return nil
}

// StoreShortURL сохраняет короткую ссылку.
func (s *DBStorage) StoreShortURL(ctx context.Context, shortURL string, originalURL string) error {
	row := s.pool.QueryRow(ctx, stmt, shortURL, originalURL, ctx.Value(common.KeyUserID))

	var url string
	var isNewURL bool

	err := row.Scan(&url, &isNewURL)
	if err != nil {
		return fmt.Errorf("failed to scan a response row: %w", err)
	}

	if !isNewURL {
		return newOriginalURLAlreadyExistError(url)
	}

	return nil
}

// StoreShortURLs сохраняет несколько коротких ссылок.
func (s *DBStorage) StoreShortURLs(ctx context.Context, urls []models.URL) error {
	batch := &pgx.Batch{}

	for _, url := range urls {
		batch.Queue(stmt, url.ShortURL, url.OriginalURL, url.UserID)
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

	return nil
}

// DeleteShortURLs мягко удаляет ссылки.
func (s *DBStorage) DeleteShortURLs(ctx context.Context, urls []string) error {
	const stmt = `UPDATE urls SET is_deleted = true WHERE short_url = $1`

	batch := &pgx.Batch{}

	for _, url := range urls {
		batch.Queue(stmt, url)
	}

	result := s.pool.SendBatch(ctx, batch)
	defer func() {
		if err := result.Close(); err != nil {
			s.logger.Error("failed to close batch result", zap.Error(err))
		}
	}()

	_, err := result.Exec()
	if err != nil {
		return fmt.Errorf("unable to update batch: %w", err)
	}

	return nil
}

// GetURL получает оригинальную ссылку по короткой.
func (s *DBStorage) GetURL(ctx context.Context, shortURL string) (models.URL, error) {
	const queryStmt = `SELECT id, short_url, original_url, is_deleted, user_id
		FROM urls
		WHERE short_url = $1
		LIMIT 1`

	row := s.pool.QueryRow(ctx, queryStmt, shortURL)

	var u models.URL
	err := row.Scan(&u.ID, &u.ShortURL, &u.OriginalURL, &u.DeletedFlag, &u.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.URL{}, fmt.Errorf("%w for short URL %s", ErrURLNotFound, shortURL)
		}

		return models.URL{}, fmt.Errorf("failed to scan a response row: %w", err)
	}

	return u, nil
}

// FetchUserURLs получает все пользовательские ссылки.
func (s *DBStorage) FetchUserURLs(ctx context.Context) ([]models.URL, error) {
	const queryStmt = `SELECT id, short_url, original_url, user_id
		FROM urls
		WHERE user_id = $1`

	urls := []models.URL{}

	rows, err := s.pool.Query(ctx, queryStmt, ctx.Value(common.KeyUserID))
	if err != nil {
		return []models.URL{}, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u models.URL
		err = rows.Scan(&u.ID, &u.ShortURL, &u.OriginalURL, &u.UserID)
		if err != nil {
			return []models.URL{}, fmt.Errorf("failed to scan query: %w", err)
		}

		urls = append(urls, u)
	}

	rowsErr := rows.Err()
	if rowsErr != nil {
		return []models.URL{}, fmt.Errorf("failed to read query: %w", err)
	}

	return urls, nil
}

// DropDeletedURLs очищает из БД удаленные ссылки.
func (s *DBStorage) DropDeletedURLs(ctx context.Context) error {
	const stmt = `DELETE FROM urls WHERE is_deleted = true`

	_, err := s.pool.Exec(ctx, stmt)
	if err != nil {
		return fmt.Errorf("failed to execute drop query: %w", err)
	}

	return nil
}

// FetchStats получает статистические данные.
func (s *DBStorage) FetchStats(ctx context.Context) (int, int, error) {
	const queryStmt = `SELECT count(*), count(DISTINCT user_id) FROM urls`

	row := s.pool.QueryRow(ctx, queryStmt)

	var urlsCount int
	var usersCount int

	err := row.Scan(&urlsCount, &usersCount)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to scan a response row: %w", err)
	}

	return urlsCount, usersCount, nil
}

// Ping проверяет работоспособность БД.
func (s *DBStorage) Ping(ctx context.Context) error {
	if err := s.pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping DB: %w", err)
	}

	return nil
}

// Close закрывает соединение с БД.
func (s *DBStorage) Close() error {
	s.pool.Close()
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

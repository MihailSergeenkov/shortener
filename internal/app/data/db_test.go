package data

import (
	"context"
	"errors"
	"testing"

	"github.com/MihailSergeenkov/shortener/internal/app/common"
	"github.com/MihailSergeenkov/shortener/internal/app/data/mock"
	"github.com/MihailSergeenkov/shortener/internal/app/models"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestDBStoreShortURL(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pool := mock.NewMockDBPooler(mockCtrl)
	logger := zap.NewNop()
	storage := DBStorage{
		pool:   pool,
		logger: logger,
	}
	currentUserID := "some_id"
	ctx := context.WithValue(context.Background(), common.KeyUserID, currentUserID)
	row := mock.NewMockRow(mockCtrl)
	shortURL := "short_url"
	originalURL := "some_url"

	tests := []struct {
		name    string
		errText string
		rowErr  error
	}{
		{
			name:    "url already exist",
			errText: "original url already exist",
			rowErr:  nil,
		},
		{
			name:    "failed read row",
			errText: "failed to scan a response row",
			rowErr:  errors.New("some error"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pool.EXPECT().QueryRow(ctx, stmt, shortURL, originalURL, currentUserID).Times(1).Return(row)

			row.EXPECT().Scan(gomock.Any()).Times(1).Return(test.rowErr)

			err := storage.StoreShortURL(ctx, shortURL, originalURL)

			require.Error(t, err)
			require.ErrorContains(t, err, test.errText)
		})
	}
}

func TestDBStoreShortURLs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pool := mock.NewMockDBPooler(mockCtrl)
	logger := zap.NewNop()
	storage := DBStorage{
		pool:   pool,
		logger: logger,
	}
	ctx := context.Background()
	batchResults := mock.NewMockBatchResults(mockCtrl)
	urls := []models.URL{
		{
			ShortURL:    "short_url",
			OriginalURL: "some_url",
			UserID:      "some_id",
		},
	}
	tests := []struct {
		name      string
		wantErr   bool
		resultErr error
	}{
		{
			name:      "success store",
			wantErr:   false,
			resultErr: nil,
		},
		{
			name:      "failed exec batch",
			wantErr:   true,
			resultErr: errors.New("some error"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pool.EXPECT().SendBatch(ctx, gomock.Any()).Times(1).Return(batchResults)

			batchResults.EXPECT().Exec().Times(1).Return(pgconn.CommandTag{}, test.resultErr)
			batchResults.EXPECT().Close().Times(1)

			err := storage.StoreShortURLs(ctx, urls)

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, "unable to insert batch")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDBDeleteShortURLs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pool := mock.NewMockDBPooler(mockCtrl)
	logger := zap.NewNop()
	storage := DBStorage{
		pool:   pool,
		logger: logger,
	}
	ctx := context.Background()
	batchResults := mock.NewMockBatchResults(mockCtrl)
	urls := []string{"short_url"}

	tests := []struct {
		name      string
		wantErr   bool
		resultErr error
	}{
		{
			name:      "success delete",
			wantErr:   false,
			resultErr: nil,
		},
		{
			name:      "failed exec batch",
			wantErr:   true,
			resultErr: errors.New("some error"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pool.EXPECT().SendBatch(ctx, gomock.Any()).Times(1).Return(batchResults)

			batchResults.EXPECT().Exec().Times(1).Return(pgconn.CommandTag{}, test.resultErr)
			batchResults.EXPECT().Close().Times(1)

			err := storage.DeleteShortURLs(ctx, urls)

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, "unable to update batch")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDBGetURL(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pool := mock.NewMockDBPooler(mockCtrl)
	logger := zap.NewNop()
	storage := DBStorage{
		pool:   pool,
		logger: logger,
	}
	ctx := context.Background()
	stmt := `SELECT id, short_url, original_url, is_deleted, user_id
		FROM urls
		WHERE short_url = $1
		LIMIT 1`

	row := mock.NewMockRow(mockCtrl)
	shortURL := "short_url"

	tests := []struct {
		name    string
		wantErr bool
		errText string
		rowErr  error
	}{
		{
			name:    "success get",
			wantErr: false,
			errText: "",
			rowErr:  nil,
		},
		{
			name:    "failed read row",
			wantErr: true,
			errText: "failed to scan a response row",
			rowErr:  errors.New("some error"),
		},
		{
			name:    "not found row",
			wantErr: true,
			errText: "url not found",
			rowErr:  pgx.ErrNoRows,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pool.EXPECT().QueryRow(ctx, stmt, shortURL).Times(1).Return(row)

			row.EXPECT().Scan(gomock.Any()).Times(1).Return(test.rowErr)

			_, err := storage.GetURL(ctx, shortURL)

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, test.errText)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDBFetchUserURLs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pool := mock.NewMockDBPooler(mockCtrl)
	logger := zap.NewNop()
	storage := DBStorage{
		pool:   pool,
		logger: logger,
	}
	currentUserID := "some_id"
	ctx := context.WithValue(context.Background(), common.KeyUserID, currentUserID)
	stmt := `SELECT id, short_url, original_url, user_id
		FROM urls
		WHERE user_id = $1`

	rows := mock.NewMockRows(mockCtrl)

	tests := []struct {
		name    string
		wantErr bool
		errText string
		rowsErr error
	}{
		{
			name:    "success fetch",
			wantErr: false,
			errText: "",
			rowsErr: nil,
		},
		{
			name:    "failed read rows",
			wantErr: true,
			errText: "failed to read query",
			rowsErr: errors.New("some error"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pool.EXPECT().Query(ctx, stmt, currentUserID).Times(1).Return(rows, nil)

			rows.EXPECT().Close().Times(1)
			rows.EXPECT().Next().Times(1).Return(false)
			rows.EXPECT().Err().Times(1).Return(test.rowsErr)

			_, err := storage.FetchUserURLs(ctx)

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, test.errText)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDBFetchUserURLs_Failed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pool := mock.NewMockDBPooler(mockCtrl)
	logger := zap.NewNop()
	storage := DBStorage{
		pool:   pool,
		logger: logger,
	}
	currentUserID := "some_id"
	ctx := context.WithValue(context.Background(), common.KeyUserID, currentUserID)
	stmt := `SELECT id, short_url, original_url, user_id
		FROM urls
		WHERE user_id = $1`

	rows := mock.NewMockRows(mockCtrl)
	someErr := errors.New("some error")

	t.Run("failed fetch", func(t *testing.T) {
		pool.EXPECT().Query(ctx, stmt, currentUserID).Times(1).Return(rows, someErr)

		_, err := storage.FetchUserURLs(ctx)

		require.Error(t, err)
		require.ErrorContains(t, err, "failed to execute query")
	})
}

func TestDBDropDeletedURLs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pool := mock.NewMockDBPooler(mockCtrl)
	logger := zap.NewNop()
	storage := DBStorage{
		pool:   pool,
		logger: logger,
	}
	ctx := context.Background()
	stmt := `DELETE FROM urls WHERE is_deleted = true`

	tests := []struct {
		name    string
		wantErr bool
		err     error
	}{
		{
			name:    "success drop",
			wantErr: false,
			err:     nil,
		},
		{
			name:    "failed drop",
			wantErr: true,
			err:     errors.New("some error"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pool.EXPECT().Exec(ctx, stmt).Times(1).Return(pgconn.CommandTag{}, test.err)

			err := storage.DropDeletedURLs(ctx)

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, "failed to execute drop query")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDBPing(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pool := mock.NewMockDBPooler(mockCtrl)
	logger := zap.NewNop()
	storage := DBStorage{
		pool:   pool,
		logger: logger,
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		wantErr bool
		err     error
	}{
		{
			name:    "success ping",
			wantErr: false,
			err:     nil,
		},
		{
			name:    "failed ping",
			wantErr: true,
			err:     errors.New("some error"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pool.EXPECT().Ping(ctx).Times(1).Return(test.err)

			err := storage.Ping(ctx)

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, "failed to ping DB")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDBClose(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pool := mock.NewMockDBPooler(mockCtrl)
	logger := zap.NewNop()
	storage := DBStorage{
		pool:   pool,
		logger: logger,
	}

	t.Run("success close", func(t *testing.T) {
		pool.EXPECT().Close().Times(1)

		err := storage.Close()
		require.NoError(t, err)
	})
}

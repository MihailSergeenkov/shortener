package data

import (
	"context"
	"testing"

	"github.com/MihailSergeenkov/shortener/internal/app/common"
	"github.com/MihailSergeenkov/shortener/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBaseStorage(t *testing.T) {
	t.Run("create base storage", func(t *testing.T) {
		storage := NewBaseStorage()

		assert.IsType(t, (*BaseStorage)(nil), storage)
		assert.Implements(t, (*Storager)(nil), storage)
	})
}

func TestStoreShortURL(t *testing.T) {
	shortURL := "short_url"
	originalURL := "some_url"
	currentUserID := "some_id"
	ctx := context.WithValue(context.Background(), common.KeyUserID, currentUserID)

	tests := []struct {
		name    string
		storage *BaseStorage
		wantErr bool
		errText string
	}{
		{
			name:    "success store",
			storage: NewBaseStorage(),
			wantErr: false,
			errText: "",
		},
		{
			name: "short url already exist",
			storage: &BaseStorage{
				urls: map[string]models.URL{
					shortURL: {
						ShortURL:    shortURL,
						OriginalURL: originalURL,
					},
				},
			},
			wantErr: true,
			errText: "short url already exist",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.storage.StoreShortURL(ctx, shortURL, originalURL)

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, test.errText)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestStoreShortURL_FailedContext(t *testing.T) {
	shortURL := "short_url"
	originalURL := "some_url"

	ctx := context.Background()
	storage := NewBaseStorage()

	t.Run("context without user id", func(t *testing.T) {
		err := storage.StoreShortURL(ctx, shortURL, originalURL)

		require.Error(t, err)
		require.ErrorContains(t, err, "failed to fetch user id from context")
	})
}

func TestStoreShortURLs(t *testing.T) {
	ctx := context.Background()
	shortURL := "short_url"
	url := models.URL{
		ShortURL:    shortURL,
		OriginalURL: "some_url",
	}

	tests := []struct {
		name    string
		storage *BaseStorage
		wantErr bool
	}{
		{
			name:    "success store",
			storage: NewBaseStorage(),
			wantErr: false,
		},
		{
			name: "short url already exist",
			storage: &BaseStorage{
				urls: map[string]models.URL{
					shortURL: url,
				},
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.storage.StoreShortURLs(ctx, []models.URL{url})

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, "short url already exist")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetURL(t *testing.T) {
	ctx := context.Background()
	shortURL := "short_url"
	url := models.URL{
		ShortURL:    shortURL,
		OriginalURL: "some_url",
	}

	tests := []struct {
		name    string
		storage *BaseStorage
		wantErr bool
	}{
		{
			name: "success get",
			storage: &BaseStorage{
				urls: map[string]models.URL{
					shortURL: url,
				},
			},
			wantErr: false,
		},
		{
			name:    "short url not found",
			storage: NewBaseStorage(),
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			u, err := test.storage.GetURL(ctx, shortURL)

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, "url not found")
			} else {
				require.NoError(t, err)
				assert.Equal(t, url, u)
			}
		})
	}
}

func TestFetchUserURLs(t *testing.T) {
	currentUserID := "some_id"
	ctx := context.WithValue(context.Background(), common.KeyUserID, currentUserID)
	shortURL := "short_url"
	url := models.URL{
		ShortURL:    shortURL,
		OriginalURL: "some_url",
		UserID:      currentUserID,
	}

	tests := []struct {
		name        string
		storage     *BaseStorage
		fetchedURLs []models.URL
	}{
		{
			name: "success fetch",
			storage: &BaseStorage{
				urls: map[string]models.URL{
					shortURL: url,
				},
			},
			fetchedURLs: []models.URL{url},
		},
		{
			name:        "urls not found",
			storage:     NewBaseStorage(),
			fetchedURLs: []models.URL{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			u, err := test.storage.FetchUserURLs(ctx)

			require.NoError(t, err)
			assert.Equal(t, test.fetchedURLs, u)
		})
	}
}

func TestDeleteShortURLs(t *testing.T) {
	ctx := context.Background()
	shortURL := "short_url"
	urls := []string{shortURL}

	tests := []struct {
		name    string
		storage *BaseStorage
		wantErr bool
	}{
		{
			name: "success delete",
			storage: &BaseStorage{
				urls: map[string]models.URL{
					shortURL: {
						ShortURL:    shortURL,
						OriginalURL: "some_url",
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "short url not found",
			storage: NewBaseStorage(),
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.storage.DeleteShortURLs(ctx, urls)

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, "url not found")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDropDeletedURLs(t *testing.T) {
	storage := NewBaseStorage()
	ctx := context.Background()
	err := storage.DropDeletedURLs(ctx)

	require.NoError(t, err)
}

func TestPing(t *testing.T) {
	storage := NewBaseStorage()
	ctx := context.Background()
	err := storage.Ping(ctx)

	require.NoError(t, err)
}

func TestClose(t *testing.T) {
	storage := NewBaseStorage()
	err := storage.Close()

	require.NoError(t, err)
}

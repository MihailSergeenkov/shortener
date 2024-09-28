package data

import (
	"context"
	"testing"

	"github.com/MihailSergeenkov/shortener/internal/app/common"
	"github.com/MihailSergeenkov/shortener/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewFileStorage(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name            string
		fileStoragePath string
		wantErr         bool
	}{
		{
			name:            "success store",
			fileStoragePath: "/tmp/short-url-db.json",
			wantErr:         false,
		},
		{
			name:            "failed open file",
			fileStoragePath: "",
			wantErr:         true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage, err := NewFileStorage(logger, test.fileStoragePath)

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, "failed to open file storage")
			} else {
				require.NoError(t, err)
				assert.IsType(t, (*FileStorage)(nil), storage)
				assert.Implements(t, (*Storager)(nil), storage)
			}
		})
	}
}

func TestFileStoreShortURL(t *testing.T) {
	shortURL := "short_url"
	originalURL := "some_url"
	currentUserID := "some_id"
	ctx := context.WithValue(context.Background(), common.KeyUserID, currentUserID)

	tests := []struct {
		name    string
		storage *FileStorage
		wantErr bool
		errText string
	}{
		{
			name: "success store",
			storage: &FileStorage{
				baseStorage:     *NewBaseStorage(),
				fileStoragePath: "/tmp/short-url-db.json",
				logger:          zap.NewNop(),
			},
			wantErr: false,
		},
		{
			name: "failed open file",
			storage: &FileStorage{
				baseStorage:     *NewBaseStorage(),
				fileStoragePath: "",
				logger:          zap.NewNop(),
			},
			wantErr: true,
			errText: "failed to open file storage",
		},
		{
			name: "short url already exist",
			storage: &FileStorage{
				baseStorage: BaseStorage{
					urls: map[string]models.URL{
						shortURL: {
							ShortURL:    shortURL,
							OriginalURL: originalURL,
						},
					},
				},
				fileStoragePath: "/tmp/short-url-db.json",
				logger:          zap.NewNop(),
			},
			wantErr: true,
			errText: "failed to add url",
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

func TestFileStoreShortURLs(t *testing.T) {
	ctx := context.Background()
	shortURL := "short_url"
	url := models.URL{
		ShortURL:    shortURL,
		OriginalURL: "some_url",
	}

	tests := []struct {
		name    string
		storage *FileStorage
		wantErr bool
		errText string
	}{
		{
			name: "success store",
			storage: &FileStorage{
				baseStorage:     *NewBaseStorage(),
				fileStoragePath: "/tmp/short-url-db.json",
				logger:          zap.NewNop(),
			},
			wantErr: false,
		},
		{
			name: "failed open file",
			storage: &FileStorage{
				baseStorage:     *NewBaseStorage(),
				fileStoragePath: "",
				logger:          zap.NewNop(),
			},
			wantErr: true,
			errText: "failed to open file storage",
		},
		{
			name: "short url already exist",
			storage: &FileStorage{
				baseStorage: BaseStorage{
					urls: map[string]models.URL{
						shortURL: url,
					},
				},
				fileStoragePath: "/tmp/short-url-db.json",
				logger:          zap.NewNop(),
			},
			wantErr: true,
			errText: "failed to add urls",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.storage.StoreShortURLs(ctx, []models.URL{url})

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, test.errText)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestFileFetchUserURLs(t *testing.T) {
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
		storage     *FileStorage
		fetchedURLs []models.URL
	}{
		{
			name: "success fetch",
			storage: &FileStorage{
				baseStorage: BaseStorage{
					urls: map[string]models.URL{
						shortURL: url,
					},
				},
			},
			fetchedURLs: []models.URL{url},
		},
		{
			name: "urls not found",
			storage: &FileStorage{
				baseStorage: *NewBaseStorage(),
			},
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

func TestFileGetURL(t *testing.T) {
	ctx := context.Background()
	shortURL := "short_url"
	url := models.URL{
		ShortURL:    shortURL,
		OriginalURL: "some_url",
	}

	tests := []struct {
		name    string
		storage *FileStorage
		wantErr bool
	}{
		{
			name: "success get",
			storage: &FileStorage{
				baseStorage: BaseStorage{
					urls: map[string]models.URL{
						shortURL: url,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "short url not found",
			storage: &FileStorage{
				baseStorage: *NewBaseStorage(),
			},
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

func TestFileDeleteShortURLs(t *testing.T) {
	ctx := context.Background()
	shortURL := "short_url"
	urls := []string{shortURL}

	tests := []struct {
		name    string
		storage *FileStorage
		wantErr bool
		errText string
	}{
		{
			name: "success delete",
			storage: &FileStorage{
				baseStorage: BaseStorage{
					urls: map[string]models.URL{
						shortURL: {
							ShortURL:    shortURL,
							OriginalURL: "some_url",
						},
					},
				},
				fileStoragePath: "/tmp/short-url-db.json",
				logger:          zap.NewNop(),
			},
			wantErr: false,
		},
		{
			name: "failed open file",
			storage: &FileStorage{
				baseStorage:     *NewBaseStorage(),
				fileStoragePath: "",
				logger:          zap.NewNop(),
			},
			wantErr: true,
			errText: "failed to open file storage",
		},
		{
			name: "failed delete url",
			storage: &FileStorage{
				baseStorage:     *NewBaseStorage(),
				fileStoragePath: "/tmp/short-url-db.json",
				logger:          zap.NewNop(),
			},
			wantErr: true,
			errText: "failed to delete urls",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.storage.DeleteShortURLs(ctx, urls)

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, test.errText)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestFileDropDeletedURLs(t *testing.T) {
	storage := FileStorage{}
	ctx := context.Background()
	err := storage.DropDeletedURLs(ctx)

	require.NoError(t, err)
}

func TestFileFetchStats(t *testing.T) {
	ctx := context.Background()

	type want struct {
		urls  int
		users int
	}
	tests := []struct {
		name    string
		storage *FileStorage
		want    want
	}{
		{
			name: "success fetch",
			storage: &FileStorage{
				baseStorage: BaseStorage{
					urls: map[string]models.URL{
						"short_url": {
							ShortURL:    "short_url",
							OriginalURL: "some_url",
							UserID:      "1",
						},
						"short_url_2": {
							ShortURL:    "short_url_2",
							OriginalURL: "some_url",
							UserID:      "1",
						},
					},
				},
			},
			want: want{
				urls:  2,
				users: 1,
			},
		},
		{
			name: "short url not found",
			storage: &FileStorage{
				baseStorage: *NewBaseStorage(),
			},
			want: want{
				urls:  0,
				users: 0,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			urls, users, err := test.storage.FetchStats(ctx)

			require.NoError(t, err)
			assert.Equal(t, test.want.urls, urls)
			assert.Equal(t, test.want.users, users)
		})
	}
}

func TestFilePing(t *testing.T) {
	storage := FileStorage{}
	ctx := context.Background()
	err := storage.Ping(ctx)

	require.NoError(t, err)
}

func TestFileClose(t *testing.T) {
	storage := FileStorage{}
	err := storage.Close()

	require.NoError(t, err)
}

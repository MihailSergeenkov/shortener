package services

import (
	"context"
	"errors"
	"testing"

	"github.com/MihailSergeenkov/shortener/internal/app/common"
	"github.com/MihailSergeenkov/shortener/internal/app/data/mock"
	"github.com/MihailSergeenkov/shortener/internal/app/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAddShortURL_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ctx := context.Background()
	store := mock.NewMockStorager(mockCtrl)
	originalURL := "some_url"

	store.EXPECT().StoreShortURL(ctx, gomock.Any(), originalURL).Times(1).Return(nil)

	t.Run("add short URL success", func(t *testing.T) {
		_, err := AddShortURL(ctx, store, originalURL)
		assert.NoError(t, err)
	})
}

func TestAddShortURL_Failed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ctx := context.Background()
	store := mock.NewMockStorager(mockCtrl)
	originalURL := "some_url"
	errSome := errors.New("some error")

	store.EXPECT().StoreShortURL(ctx, gomock.Any(), originalURL).Times(1).Return(errSome)

	t.Run("add short URL failed", func(t *testing.T) {
		_, err := AddShortURL(ctx, store, originalURL)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to store short URL", "some error")
	})
}

func BenchmarkAddShortURL(b *testing.B) {
	mockCtrl := gomock.NewController(b)
	defer mockCtrl.Finish()

	ctx := context.Background()
	store := mock.NewMockStorager(mockCtrl)
	originalURL := "some_url"

	store.EXPECT().StoreShortURL(ctx, gomock.Any(), originalURL).AnyTimes().Return(nil)

	b.ResetTimer()

	b.Run("AddShortURL", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			AddShortURL(ctx, store, originalURL) //nolint:errcheck,gosec // Тест бенчмарка, ошибки проверяются в основном тесте
		}
	})
}

func TestAddBatchShortURL_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	currentUserID := "some_id"
	ctx := context.WithValue(context.Background(), common.KeyUserID, currentUserID)
	store := mock.NewMockStorager(mockCtrl)
	batch := models.BatchRequest{
		models.BatchDataRequest{
			CorrelationID: "some_id",
			OriginalURL:   "some_url",
		},
	}

	store.EXPECT().StoreShortURLs(ctx, gomock.Any()).Times(1).Return(nil)

	t.Run("add batch short URL success", func(t *testing.T) {
		_, err := AddBatchShortURL(ctx, store, batch)
		assert.NoError(t, err)
	})
}

func TestAddBatchShortURL_Failed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	currentUserID := "some_id"
	ctx := context.WithValue(context.Background(), common.KeyUserID, currentUserID)
	store := mock.NewMockStorager(mockCtrl)
	batch := models.BatchRequest{
		models.BatchDataRequest{
			CorrelationID: "some_id",
			OriginalURL:   "some_url",
		},
	}
	errSome := errors.New("some error")

	store.EXPECT().StoreShortURLs(ctx, gomock.Any()).Times(1).Return(errSome)

	t.Run("add batch short URL failed", func(t *testing.T) {
		_, err := AddBatchShortURL(ctx, store, batch)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to store short URLs", "some error")
	})
}

func BenchmarkAddBatchShortURL(b *testing.B) {
	mockCtrl := gomock.NewController(b)
	defer mockCtrl.Finish()

	currentUserID := "some_id"
	ctx := context.WithValue(context.Background(), common.KeyUserID, currentUserID)
	store := mock.NewMockStorager(mockCtrl)
	batch := models.BatchRequest{
		models.BatchDataRequest{
			CorrelationID: "some_id",
			OriginalURL:   "some_url",
		},
	}

	store.EXPECT().StoreShortURLs(ctx, gomock.Any()).AnyTimes().Return(nil)

	b.ResetTimer()

	b.Run("AddBatchShortURL", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			AddBatchShortURL(ctx, store, batch) //nolint:errcheck,gosec // Тест бенчмарка, ошибки проверяются в основном тесте
		}
	})
}

func TestFetchUserURLs_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ctx := context.Background()
	store := mock.NewMockStorager(mockCtrl)
	urls := []models.URL{
		{
			ShortURL:    "some_url",
			OriginalURL: "some_url",
			UserID:      "some_id",
			ID:          1,
			DeletedFlag: false,
		},
	}

	store.EXPECT().FetchUserURLs(ctx).Times(1).Return(urls, nil)

	t.Run("fetch user URLs success", func(t *testing.T) {
		_, err := FetchUserURLs(ctx, store)
		assert.NoError(t, err)
	})
}

func TestFetchUserURLs_Failed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ctx := context.Background()
	store := mock.NewMockStorager(mockCtrl)
	errSome := errors.New("some error")

	store.EXPECT().FetchUserURLs(ctx).Times(1).Return([]models.URL{}, errSome)

	t.Run("fetch user URLs failed", func(t *testing.T) {
		_, err := FetchUserURLs(ctx, store)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to fetch URLs", "some error")
	})
}

func BenchmarkFetchUserURLs(b *testing.B) {
	mockCtrl := gomock.NewController(b)
	defer mockCtrl.Finish()

	ctx := context.Background()
	store := mock.NewMockStorager(mockCtrl)
	urls := []models.URL{
		{
			ShortURL:    "some_url",
			OriginalURL: "some_url",
			UserID:      "some_id",
			ID:          1,
			DeletedFlag: false,
		},
	}

	store.EXPECT().FetchUserURLs(ctx).AnyTimes().Return(urls, nil)

	b.ResetTimer()

	b.Run("FetchUserURLs", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			FetchUserURLs(ctx, store) //nolint:errcheck,gosec // Тест бенчмарка, ошибки проверяются в основном тесте
		}
	})
}

func TestDeleteUserURLs_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	currentUserID := "some_id"
	ctx := context.WithValue(context.Background(), common.KeyUserID, currentUserID)
	logger := zap.NewNop()
	store := mock.NewMockStorager(mockCtrl)

	shortURL := "some_url"

	type request struct {
		urls []string
		mURL models.URL
	}
	type want struct {
		urls []string
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "delete user URLs success",
			request: request{
				urls: []string{shortURL},
				mURL: models.URL{
					ShortURL:    shortURL,
					OriginalURL: "some_url",
					UserID:      currentUserID,
					ID:          1,
					DeletedFlag: false,
				},
			},
			want: want{
				urls: []string{shortURL},
			},
		},
		{
			name: "no urls for delete",
			request: request{
				urls: []string{shortURL},
				mURL: models.URL{
					ShortURL:    shortURL,
					OriginalURL: "some_url",
					UserID:      "some_other_id",
					ID:          1,
					DeletedFlag: false,
				},
			},
			want: want{
				urls: []string{},
			},
		},
	}

	for _, test := range tests {
		store.EXPECT().GetURL(ctx, shortURL).Times(1).Return(test.request.mURL, nil)
		store.EXPECT().DeleteShortURLs(ctx, test.want.urls).Times(1).Return(nil)

		t.Run(test.name, func(t *testing.T) {
			err := DeleteUserURLs(ctx, logger, store, test.request.urls)
			assert.NoError(t, err)
		})
	}
}

func TestDeleteUserURLs_Failed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	currentUserID := "some_id"
	ctx := context.WithValue(context.Background(), common.KeyUserID, currentUserID)
	logger := zap.NewNop()
	store := mock.NewMockStorager(mockCtrl)
	shortURL := "some_url"
	urls := []string{shortURL}
	mURL := models.URL{
		ShortURL:    shortURL,
		OriginalURL: "some_url",
		UserID:      currentUserID,
		ID:          1,
		DeletedFlag: false,
	}
	errSome := errors.New("some error")

	store.EXPECT().GetURL(ctx, shortURL).Times(1).Return(mURL, nil)
	store.EXPECT().DeleteShortURLs(ctx, urls).Times(1).Return(errSome)

	t.Run("delete user URLs failed", func(t *testing.T) {
		err := DeleteUserURLs(ctx, logger, store, urls)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to delete URLs", "some error")
	})
}

func BenchmarkDeleteUserURLs(b *testing.B) {
	mockCtrl := gomock.NewController(b)
	defer mockCtrl.Finish()

	currentUserID := "some_id"
	ctx := context.WithValue(context.Background(), common.KeyUserID, currentUserID)
	logger := zap.NewNop()
	store := mock.NewMockStorager(mockCtrl)
	shortURL := "some_url"
	urls := []string{shortURL}
	mURL := models.URL{
		ShortURL:    shortURL,
		OriginalURL: "some_url",
		UserID:      currentUserID,
		ID:          1,
		DeletedFlag: false,
	}

	store.EXPECT().GetURL(ctx, shortURL).AnyTimes().Return(mURL, nil)
	store.EXPECT().DeleteShortURLs(ctx, urls).AnyTimes().Return(nil)

	b.ResetTimer()

	b.Run("DeleteUserURLs", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			DeleteUserURLs(ctx, logger, store, urls) //nolint:errcheck,gosec // Тест бенчмарка
		}
	})
}

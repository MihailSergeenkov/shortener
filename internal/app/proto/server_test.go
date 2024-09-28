package proto

import (
	context "context"
	"errors"
	"testing"

	"github.com/MihailSergeenkov/shortener/internal/app/common"
	"github.com/MihailSergeenkov/shortener/internal/app/data"
	"github.com/MihailSergeenkov/shortener/internal/app/data/mock"
	"github.com/MihailSergeenkov/shortener/internal/app/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestAuthInterceptor(t *testing.T) {
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return req, nil
	}
	tests := []struct {
		name     string
		method   string
		metadata map[string]string
		wantErr  bool
	}{
		{
			name:     "without auth",
			method:   "/shortener.Shortener/GetURL",
			metadata: map[string]string{},
			wantErr:  false,
		},
		{
			name:     "with auth",
			method:   "/shortener.Shortener/AddShortURL",
			metadata: map[string]string{"user_id": "12345"},
			wantErr:  false,
		},
		{
			name:     "without user id",
			method:   "/shortener.Shortener/AddShortURL",
			metadata: map[string]string{},
			wantErr:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			md := metadata.New(test.metadata)
			ctx := metadata.NewIncomingContext(context.Background(), md)

			info := &grpc.UnaryServerInfo{
				FullMethod: test.method,
			}
			_, err := authInterceptor(ctx, "test", info, handler)

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, "missing user id")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewGRPCServer(t *testing.T) {
	t.Run("init gRPC server", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		logger := zap.NewNop()
		storage := mock.NewMockStorager(mockCtrl)

		s := NewGRPCServer(logger, storage)
		assert.IsType(t, (*grpc.Server)(nil), s)
	})
}

func TestAddShortURL(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zap.NewNop()
	storage := mock.NewMockStorager(mockCtrl)

	server := ProtoServer{
		logger:  logger,
		storage: storage,
	}
	currentUserID := "some_id"
	ctx := context.WithValue(context.Background(), common.KeyUserID, currentUserID)
	shortURL := "short_url"
	originalURL := "some_url"

	tests := []struct {
		name    string
		wantErr bool
		err     error
	}{
		{
			name:    "success add",
			wantErr: false,
			err:     nil,
		},
		{
			name:    "when original url already exist",
			wantErr: false,
			err:     &data.OriginalURLAlreadyExistError{ShortURL: shortURL},
		},
		{
			name:    "failed add",
			wantErr: true,
			err:     errors.New("some error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.EXPECT().StoreShortURL(ctx, gomock.Any(), originalURL).Times(1).Return(test.err)

			resp, err := server.AddShortURL(ctx, &AddShortURLRequest{OriginalUrl: originalURL})

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, "failed to add URL to storage")
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, resp.GetShortUrl())
			}
		})
	}
}

func TestAddShortURLs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zap.NewNop()
	storage := mock.NewMockStorager(mockCtrl)

	server := ProtoServer{
		logger:  logger,
		storage: storage,
	}
	currentUserID := "some_id"
	ctx := context.WithValue(context.Background(), common.KeyUserID, currentUserID)
	correlationID := "correlation_id"
	originalURL := "some_url"

	tests := []struct {
		name    string
		wantErr bool
		err     error
	}{
		{
			name:    "success add",
			wantErr: false,
			err:     nil,
		},
		{
			name:    "failed add",
			wantErr: true,
			err:     errors.New("some error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.EXPECT().StoreShortURLs(ctx, gomock.Any()).Times(1).Return(test.err)

			resp, err := server.AddShortURLs(ctx, &AddShortURLsRequest{
				Urls: []*BatchRequest{
					{
						CorrelationId: correlationID,
						OriginalUrl:   originalURL,
					},
				},
			})

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, "failed to add URLs to storage")
			} else {
				require.NoError(t, err)
				assert.Equal(t, correlationID, resp.GetUrls()[0].GetCorrelationId())
				assert.NotEmpty(t, resp.GetUrls()[0].GetShortUrl())
			}
		})
	}
}

func TestGetURL(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zap.NewNop()
	storage := mock.NewMockStorager(mockCtrl)

	server := ProtoServer{
		logger:  logger,
		storage: storage,
	}
	ctx := context.Background()
	shortURL := "some_url"
	originalURL := "some_url"

	tests := []struct {
		name    string
		url     models.URL
		err     error
		wantErr bool
		errText string
	}{
		{
			name: "success get",
			url: models.URL{
				ShortURL:    shortURL,
				OriginalURL: originalURL,
				DeletedFlag: false,
			},
			err:     nil,
			wantErr: false,
			errText: "",
		},
		{
			name:    "failed get",
			url:     models.URL{},
			err:     errors.New("some error"),
			wantErr: true,
			errText: "ailed to fetch URL from storage",
		},
		{
			name:    "url not found",
			url:     models.URL{},
			err:     data.ErrURLNotFound,
			wantErr: true,
			errText: "URL not found",
		},
		{
			name: "when url deleted",
			url: models.URL{
				ShortURL:    shortURL,
				OriginalURL: originalURL,
				DeletedFlag: true,
			},
			err:     nil,
			wantErr: true,
			errText: "URL deleted",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.EXPECT().GetURL(ctx, shortURL).Times(1).Return(test.url, test.err)

			resp, err := server.GetURL(ctx, &GetURLRequest{ShortUrl: shortURL})

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, test.errText)
			} else {
				require.NoError(t, err)
				assert.Equal(t, originalURL, resp.GetOriginalUrl())
			}
		})
	}
}

func TestFetchUserURLs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zap.NewNop()
	storage := mock.NewMockStorager(mockCtrl)

	server := ProtoServer{
		logger:  logger,
		storage: storage,
	}
	ctx := context.Background()
	shortURL := "some_url"
	originalURL := "some_url"

	tests := []struct {
		name    string
		urls    []models.URL
		err     error
		wantErr bool
		errText string
	}{
		{
			name: "success fetch",
			urls: []models.URL{
				{
					ShortURL:    shortURL,
					OriginalURL: originalURL,
				},
			},
			err:     nil,
			wantErr: false,
			errText: "",
		},
		{
			name:    "failed fetch",
			urls:    []models.URL{},
			err:     errors.New("some error"),
			wantErr: true,
			errText: "failed to fetch URLs from storage",
		},
		{
			name:    "urls not found",
			urls:    []models.URL{},
			err:     nil,
			wantErr: true,
			errText: "URLs not found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.EXPECT().FetchUserURLs(ctx).Times(1).Return(test.urls, test.err)

			resp, err := server.FetchUserURLs(ctx, nil)

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, test.errText)
			} else {
				require.NoError(t, err)
				assert.Equal(t, shortURL, resp.GetUrls()[0].GetShortUrl())
				assert.Equal(t, originalURL, resp.GetUrls()[0].GetOriginalUrl())
			}
		})
	}
}

func TestDeleteUserURLs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zap.NewNop()
	storage := mock.NewMockStorager(mockCtrl)

	server := ProtoServer{
		logger:  logger,
		storage: storage,
	}

	currentUserID := "some_id"
	ctx := context.WithValue(context.Background(), common.KeyUserID, currentUserID)
	shortURL := "some_url"
	url := models.URL{
		ShortURL:    shortURL,
		OriginalURL: "some_url",
		UserID:      currentUserID,
		DeletedFlag: false,
	}
	urls := []string{shortURL}

	tests := []struct {
		name    string
		wantErr bool
		err     error
	}{
		{
			name:    "success delete",
			wantErr: false,
			err:     nil,
		},
		{
			name:    "failed delete",
			wantErr: true,
			err:     errors.New("some error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.EXPECT().GetURL(ctx, shortURL).Times(1).Return(url, nil)
			storage.EXPECT().DeleteShortURLs(ctx, urls).Times(1).Return(test.err)

			resp, err := server.DeleteUserURLs(ctx, &DeleteUserURLsRequest{Urls: urls})

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, "failed to delete URLs from storage")
			} else {
				require.NoError(t, err)
				assert.Equal(t, "Accepted", resp.GetText())
			}
		})
	}
}

func TestFetchStats(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zap.NewNop()
	storage := mock.NewMockStorager(mockCtrl)

	server := ProtoServer{
		logger:  logger,
		storage: storage,
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		wantErr bool
		err     error
		urls    int
		users   int
	}{
		{
			name:    "success fetch",
			wantErr: false,
			err:     nil,
			urls:    2,
			users:   1,
		},
		{
			name:    "failed fetch",
			wantErr: true,
			err:     errors.New("some error"),
			urls:    0,
			users:   0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.EXPECT().FetchStats(ctx).Times(1).Return(test.urls, test.users, test.err)

			resp, err := server.FetchStats(ctx, nil)

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, "failed to fetch stats from storage")
			} else {
				require.NoError(t, err)
				assert.Equal(t, int32(test.urls), resp.GetUrls())
				assert.Equal(t, int32(test.users), resp.GetUsers())
			}
		})
	}
}

func TestPing(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zap.NewNop()
	storage := mock.NewMockStorager(mockCtrl)

	server := ProtoServer{
		logger:  logger,
		storage: storage,
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
			storage.EXPECT().Ping(ctx).Times(1).Return(test.err)

			resp, err := server.Ping(ctx, nil)

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, "failed to connect to DB")
			} else {
				require.NoError(t, err)
				assert.Equal(t, "pong", resp.GetText())
			}
		})
	}
}

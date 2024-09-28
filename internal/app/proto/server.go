package proto

import (
	context "context"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/MihailSergeenkov/shortener/internal/app/common"
	"github.com/MihailSergeenkov/shortener/internal/app/config"
	"github.com/MihailSergeenkov/shortener/internal/app/data"
	"github.com/MihailSergeenkov/shortener/internal/app/models"
	"github.com/MihailSergeenkov/shortener/internal/app/services"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	status "google.golang.org/grpc/status"
)

// ProtoServer поддерживает все необходимые методы сервера.
type ProtoServer struct {
	UnimplementedShortenerServer

	logger  *zap.Logger
	storage data.Storager
}

//nolint:all // Функция взята из локументации к библиотеке
func loggerInterceptor(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields)/2)

		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]

			switch v := value.(type) {
			case string:
				f = append(f, zap.String(key.(string), v))
			case int:
				f = append(f, zap.Int(key.(string), v))
			case bool:
				f = append(f, zap.Bool(key.(string), v))
			default:
				f = append(f, zap.Any(key.(string), v))
			}
		}

		logger := l.WithOptions(zap.AddCallerSkip(1)).With(f...)

		switch lvl {
		case logging.LevelDebug:
			logger.Debug(msg)
		case logging.LevelInfo:
			logger.Info(msg)
		case logging.LevelWarn:
			logger.Warn(msg)
		case logging.LevelError:
			logger.Error(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

func authInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	methods := map[string]bool{
		"AddShortURL":    true,
		"AddShortURLs":   true,
		"FetchUserURLs":  true,
		"DeleteUserURLs": true,
	}

	method := strings.TrimPrefix(info.FullMethod, "/shortener.Shortener/")
	if _, ok := methods[method]; !ok {
		return handler(ctx, req)
	}

	var userID string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("user_id")
		if len(values) > 0 {
			userID = values[0]
		}
	}

	if len(userID) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing user id") //nolint:wrapcheck // FalsePositive
	}

	newContext := context.WithValue(ctx, common.KeyUserID, userID)

	return handler(newContext, req)
}

// NewGRPCServer функция инициализации gRPC сервера.
func NewGRPCServer(logger *zap.Logger, storage data.Storager) *grpc.Server {
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(loggerInterceptor(logger)),
			authInterceptor,
		),
	)
	RegisterShortenerServer(s, &ProtoServer{
		logger:  logger,
		storage: storage,
	})
	reflection.Register(s)

	return s
}

// AddShortURL реализует интерфейс сохранения короткой ссылки.
func (s *ProtoServer) AddShortURL(ctx context.Context, in *AddShortURLRequest) (*AddShortURLResponse, error) {
	var response AddShortURLResponse

	baseURL := config.Params.BaseURL
	shortURL, err := services.AddShortURL(ctx, s.storage, in.GetOriginalUrl())
	if err != nil {
		var origErr *data.OriginalURLAlreadyExistError
		if errors.As(err, &origErr) {
			newPath := path.Join(baseURL.Path, origErr.ShortURL)
			baseURL.Path = newPath
			response.ShortUrl = baseURL.String()

			return &response, nil
		}
		s.logger.Error("failed to add URL to storage", zap.Error(err))
		return nil, status.Error(codes.Aborted, "failed to add URL to storage") //nolint:wrapcheck // FalsePositive
	}

	newPath := path.Join(baseURL.Path, shortURL)
	baseURL.Path = newPath
	response.ShortUrl = baseURL.String()

	return &response, nil
}

// AddShortURLs реализует интерфейс сохранения нескольких коротких ссылок.
func (s *ProtoServer) AddShortURLs(ctx context.Context, in *AddShortURLsRequest) (*AddShortURLsResponse, error) {
	var req models.BatchRequest
	var response AddShortURLsResponse

	respURLs := make([]*BatchResponse, 0, len(in.GetUrls()))

	for _, u := range in.GetUrls() {
		r := models.BatchDataRequest{
			CorrelationID: u.GetCorrelationId(),
			OriginalURL:   u.GetOriginalUrl(),
		}
		req = append(req, r)
	}

	resp, err := services.AddBatchShortURL(ctx, s.storage, req)
	if err != nil {
		s.logger.Error("failed to add URLs to storage", zap.Error(err))
		return nil, status.Error(codes.Aborted, "failed to add URLs to storage") //nolint:wrapcheck // FalsePositive
	}

	for _, r := range resp {
		u := BatchResponse{
			CorrelationId: r.CorrelationID,
			ShortUrl:      r.ShortURL,
		}
		respURLs = append(respURLs, &u)
	}

	response.Urls = respURLs

	return &response, nil
}

// GetURL реализует интерфейс получения оригинальной ссылки по короткой.
func (s *ProtoServer) GetURL(ctx context.Context, in *GetURLRequest) (*GetURLResponse, error) {
	u, err := s.storage.GetURL(ctx, in.GetShortUrl())
	if err != nil {
		if errors.Is(err, data.ErrURLNotFound) {
			return nil, status.Error(codes.NotFound, "URL not found") //nolint:wrapcheck // FalsePositive
		}

		s.logger.Error("failed to fetch URL from storage", zap.Error(err))
		return nil, status.Error(codes.Aborted, "failed to fetch URL from storage") //nolint:wrapcheck // FalsePositive
	}

	if u.DeletedFlag {
		return nil, status.Error(codes.Unavailable, "URL deleted") //nolint:wrapcheck // FalsePositive
	}

	var response GetURLResponse
	response.OriginalUrl = u.OriginalURL

	return &response, nil
}

// FetchUserURLs реализует интерфейс получения всех сохраненных ссылок пользователя.
func (s *ProtoServer) FetchUserURLs(ctx context.Context, _ *FetchUserURLsRequest) (*FetchUserURLsResponse, error) {
	resp, err := services.FetchUserURLs(ctx, s.storage)
	if err != nil {
		s.logger.Error("failed to fetch URLs from storage", zap.Error(err))
		return nil, status.Error(codes.Aborted, "failed to fetch URLs from storage") //nolint:wrapcheck // FalsePositive
	}

	if len(resp) == 0 {
		return nil, status.Error(codes.NotFound, "URLs not found") //nolint:wrapcheck // FalsePositive
	}

	var response FetchUserURLsResponse
	respURLs := make([]*URL, 0, len(resp))

	for _, r := range resp {
		u := URL{
			OriginalUrl: r.OriginalURL,
			ShortUrl:    r.ShortURL,
		}
		respURLs = append(respURLs, &u)
	}

	response.Urls = respURLs

	return &response, nil
}

// DeleteUserURLs реализует интерфейс мягкого удаления ссылок.
func (s *ProtoServer) DeleteUserURLs(ctx context.Context, in *DeleteUserURLsRequest) (*DeleteUserURLsResponse, error) {
	err := services.DeleteUserURLs(ctx, s.logger, s.storage, in.GetUrls())
	if err != nil {
		s.logger.Error("failed to delete URLs from storage", zap.Error(err))
		return nil, status.Error(codes.Aborted, "failed to delete URLs from storage") //nolint:wrapcheck // FalsePositive
	}

	var response DeleteUserURLsResponse
	response.Text = "Accepted"

	return &response, nil
}

// FetchStats реализует интерфейс получения статистических данных.
func (s *ProtoServer) FetchStats(ctx context.Context, _ *FetchStatsRequest) (*FetchStatsResponse, error) {
	resp, err := services.FetchStats(ctx, s.storage)
	if err != nil {
		s.logger.Error("failed to fetch stats from storage", zap.Error(err))
		return nil, status.Error(codes.Aborted, "failed to fetch stats from storage") //nolint:wrapcheck // FalsePositive
	}

	var response FetchStatsResponse
	response.Urls = int32(resp.URLs)
	response.Users = int32(resp.Users)

	return &response, nil
}

// Ping реализует интерфейс проверки работоспособности БД.
func (s *ProtoServer) Ping(ctx context.Context, _ *PingRequest) (*PingResponse, error) {
	err := s.storage.Ping(ctx)
	if err != nil {
		s.logger.Error("failed to connect to DB", zap.Error(err))
		return nil, status.Error(codes.Aborted, "failed to connect to DB") //nolint:wrapcheck // FalsePositive
	}

	var response PingResponse
	response.Text = "pong"

	return &response, nil
}

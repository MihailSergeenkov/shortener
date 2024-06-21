package data

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"

	"github.com/MihailSergeenkov/shortener/internal/app/models"
	"go.uber.org/zap"
)

const (
	filePerm       fs.FileMode = 0o600
	openFileErrStr             = "failed to open file storage: %w"
)

type FileStorage struct {
	logger          *zap.Logger
	baseStorage     BaseStorage
	fileStoragePath string
}

func NewFileStorage(logger *zap.Logger, fsp string) (*FileStorage, error) {
	storage := FileStorage{
		baseStorage:     *NewBaseStorage(),
		fileStoragePath: fsp,
		logger:          logger,
	}

	file, err := os.OpenFile(fsp, os.O_RDONLY|os.O_CREATE, filePerm)

	if err != nil {
		return &FileStorage{}, fmt.Errorf(openFileErrStr, err)
	}
	defer closeFile(&storage, file)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		data := scanner.Bytes()

		url := models.URL{}
		err := json.Unmarshal(data, &url)
		if err != nil {
			return &FileStorage{}, fmt.Errorf("failed to parse file storage: %w", err)
		}

		storage.baseStorage.urls[url.ShortURL] = url
	}

	return &storage, nil
}

func (s *FileStorage) StoreShortURL(ctx context.Context, shortURL string, originalURL string) error {
	file, err := os.OpenFile(s.fileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, filePerm)
	if err != nil {
		return fmt.Errorf(openFileErrStr, err)
	}

	defer closeFile(s, file)

	baseStoreErr := s.baseStorage.StoreShortURL(ctx, shortURL, originalURL)
	if baseStoreErr != nil {
		return fmt.Errorf("failed to add url: %w", baseStoreErr)
	}

	encoder := json.NewEncoder(file)
	url := s.baseStorage.urls[shortURL]
	encoderErr := encoder.Encode(&url)

	if encoderErr != nil {
		return fmt.Errorf("failed to dump URL: %w", encoderErr)
	}

	return nil
}

func (s *FileStorage) StoreShortURLs(ctx context.Context, urls []models.URL) error {
	file, err := os.OpenFile(s.fileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, filePerm)
	if err != nil {
		return fmt.Errorf(openFileErrStr, err)
	}

	defer closeFile(s, file)

	baseStoreErr := s.baseStorage.StoreShortURLs(ctx, urls)
	if baseStoreErr != nil {
		return fmt.Errorf("failed to add urls: %w", baseStoreErr)
	}

	encoder := json.NewEncoder(file)

	for _, v := range urls {
		url := s.baseStorage.urls[v.ShortURL]
		encoderErr := encoder.Encode(&url)

		if encoderErr != nil {
			return fmt.Errorf("failed to dump URL: %w", encoderErr)
		}
	}

	return nil
}

func (s *FileStorage) FetchUserURLs(ctx context.Context) ([]models.URL, error) {
	return s.baseStorage.FetchUserURLs(ctx)
}

func (s *FileStorage) GetURL(ctx context.Context, shortURL string) (models.URL, error) {
	return s.baseStorage.GetURL(ctx, shortURL)
}

func (s *FileStorage) DeleteShortURLs(ctx context.Context, urls []string) error {
	file, err := os.OpenFile(s.fileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, filePerm)
	if err != nil {
		return fmt.Errorf(openFileErrStr, err)
	}

	defer closeFile(s, file)

	baseStoreErr := s.baseStorage.DeleteShortURLs(ctx, urls)
	if baseStoreErr != nil {
		return fmt.Errorf("failed to add urls: %w", baseStoreErr)
	}

	encoder := json.NewEncoder(file)

	for _, url := range urls {
		u := s.baseStorage.urls[url]
		encoderErr := encoder.Encode(&u)

		if encoderErr != nil {
			return fmt.Errorf("failed to dump URL: %w", encoderErr)
		}
	}

	return nil
}

func (s *FileStorage) Ping(_ context.Context) error {
	return nil
}

func (s *FileStorage) Close() error {
	return nil
}

func closeFile(s *FileStorage, file *os.File) {
	err := file.Close()

	if err != nil {
		s.logger.Error("failed to close file storage", zap.Error(err))
	}
}

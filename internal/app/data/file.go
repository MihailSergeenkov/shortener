package data

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"

	"github.com/MihailSergeenkov/shortener/internal/app/models"
	"go.uber.org/zap"
)

const filePerm fs.FileMode = 0o600

type FileStorage struct {
	baseStorage     BaseStorage
	fileStoragePath string
	logger          *zap.Logger
}

func NewFileStorage(logger *zap.Logger, fsp string) (*FileStorage, error) {
	storage := FileStorage{
		baseStorage:     *NewBaseStorage(),
		fileStoragePath: fsp,
		logger:          logger,
	}

	file, err := os.OpenFile(fsp, os.O_RDONLY|os.O_CREATE, filePerm)

	if err != nil {
		return &FileStorage{}, fmt.Errorf("failed to open file storage: %w", err)
	}

	defer func() {
		err := file.Close()

		if err != nil {
			storage.logger.Error("failed to close file storage", zap.Error(err))
		}
	}()

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

func (s *FileStorage) AddURL(originalURL string) (models.URL, error) {
	file, err := os.OpenFile(s.fileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, filePerm)
	if err != nil {
		return models.URL{}, fmt.Errorf("failed to open file storage: %w", err)
	}

	defer func() {
		err := file.Close()

		if err != nil {
			s.logger.Error("failed to close file storage", zap.Error(err))
		}
	}()

	encoder := json.NewEncoder(file)
	url, err := s.baseStorage.AddURL(originalURL)
	if err != nil {
		return models.URL{}, fmt.Errorf("failed to add url: %w", err)
	}

	encoderErr := encoder.Encode(&url)

	if encoderErr != nil {
		return models.URL{}, fmt.Errorf("failed to dump URL: %w", encoderErr)
	}

	return url, nil
}

func (s *FileStorage) FetchURL(shortURL string) (models.URL, error) {
	url, err := s.baseStorage.FetchURL(shortURL)

	if err != nil {
		return models.URL{}, fmt.Errorf("%w for short URL %s", ErrURLNotFound, shortURL)
	}

	return url, nil
}

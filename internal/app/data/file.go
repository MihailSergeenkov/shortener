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
	logger          *zap.Logger
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

func (s *FileStorage) StoreShortURL(shortURL string, originalURL string) error {
	file, err := os.OpenFile(s.fileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, filePerm)
	if err != nil {
		return fmt.Errorf("failed to open file storage: %w", err)
	}

	defer func() {
		err := file.Close()

		if err != nil {
			s.logger.Error("failed to close file storage", zap.Error(err))
		}
	}()

	baseStoreErr := s.baseStorage.StoreShortURL(shortURL, originalURL)
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

func (s *FileStorage) GetOriginalURL(shortURL string) (string, error) {
	return s.baseStorage.GetOriginalURL(shortURL)
}

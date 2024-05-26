package data

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

const (
	initSize int = 100
	keyBytes int = 8
	maxRetry int = 5
)

var (
	ErrURLNotFound = errors.New("url not found")
	ErrMaxRetry    = errors.New("generation attempts exceeded")
)

type URL struct {
	ID          uint   `json:"id"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Storage struct {
	FileStoragePath string
	URLs            map[string]URL
	dumpURLs        bool
}

func NewStorage(fsp string) (Storage, error) {
	storage := Storage{
		FileStoragePath: fsp,
		URLs:            make(map[string]URL, initSize),
		dumpURLs:        fsp != "",
	}

	if !storage.dumpURLs {
		return storage, nil
	}

	file, err := os.OpenFile(fsp, os.O_RDONLY|os.O_CREATE, 0666)

	if err != nil {
		return Storage{}, fmt.Errorf("failed to open file storage: %w", err)
	}

	defer func(f *os.File) {
		err := f.Close()

		if err != nil {
			log.Printf("failed to close file storage: %v", err)
		}
	}(file)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		data := scanner.Bytes()

		url := URL{}
		err := json.Unmarshal(data, &url)
		if err != nil {
			return Storage{}, fmt.Errorf("failed to parse file storage: %w", err)
		}

		storage.URLs[url.ShortURL] = url
	}

	return storage, nil
}

func (s *Storage) AddURL(originalURL string) (URL, error) {
	var encoder *json.Encoder

	if s.dumpURLs {
		file, err := os.OpenFile(s.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return URL{}, err
		}

		defer func(f *os.File) {
			err := f.Close()

			if err != nil {
				log.Printf("failed to close file storage: %v", err)
			}
		}(file)

		encoder = json.NewEncoder(file)
	}

	for range maxRetry {
		short_url, err := generateShortURL()
		if err != nil {
			return URL{}, err
		}

		if _, ok := s.URLs[short_url]; ok {
			continue
		}

		url := URL{
			ID:          uint(len(s.URLs) + 1),
			ShortURL:    short_url,
			OriginalURL: originalURL,
		}

		s.URLs[short_url] = url

		if s.dumpURLs {
			encoderErr := encoder.Encode(&url)

			if encoderErr != nil {
				return URL{}, encoderErr
			}
		}

		return url, nil
	}

	return URL{}, fmt.Errorf("%w for original URL %s", ErrMaxRetry, originalURL)
}

func (s *Storage) FetchURL(shortURL string) (URL, error) {
	u, ok := s.URLs[shortURL]

	if !ok {
		return URL{}, fmt.Errorf("%w for short URL %s", ErrURLNotFound, shortURL)
	}

	return u, nil
}

func generateShortURL() (string, error) {
	bytes := make([]byte, keyBytes)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate short URL error: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}

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

// type Urls map[string]string
type Url struct {
	ID          uint   `json:"id"`
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}

type Storage struct {
	FileStoragePath string
	Urls            map[string]Url
	dumpUrls        bool
}

func NewStorage(fsp string) (Storage, error) {
	storage := Storage{
		FileStoragePath: fsp,
		Urls:            make(map[string]Url, initSize),
		dumpUrls:        fsp != "",
	}

	if !storage.dumpUrls {
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

		url := Url{}
		err := json.Unmarshal(data, &url)
		if err != nil {
			return Storage{}, fmt.Errorf("failed to parse file storage: %w", err)
		}

		storage.Urls[url.ShortUrl] = url
	}

	return storage, nil
}

func (s *Storage) AddURL(original_url string) (Url, error) {
	var encoder *json.Encoder

	if s.dumpUrls {
		file, err := os.OpenFile(s.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return Url{}, err
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
		short_url, err := generateShortUrl()
		if err != nil {
			return Url{}, err
		}

		if _, ok := s.Urls[short_url]; ok {
			continue
		}

		url := Url{
			ID:          uint(len(s.Urls) + 1),
			ShortUrl:    short_url,
			OriginalUrl: original_url,
		}

		s.Urls[short_url] = url

		if s.dumpUrls {
			encoderErr := encoder.Encode(&url)

			if encoderErr != nil {
				return Url{}, encoderErr
			}
		}

		return url, nil
	}

	return Url{}, fmt.Errorf("%w for original url %s", ErrMaxRetry, original_url)
}

func (s *Storage) FetchURL(short_url string) (Url, error) {
	u, ok := s.Urls[short_url]

	if !ok {
		return Url{}, fmt.Errorf("%w for short url %s", ErrURLNotFound, short_url)
	}

	return u, nil
}

func generateShortUrl() (string, error) {
	bytes := make([]byte, keyBytes)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate short url error: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}

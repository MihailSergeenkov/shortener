package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/MihailSergeenkov/shortener/internal/app/data"
	"github.com/MihailSergeenkov/shortener/internal/app/models"
)

const (
	keyBytes int = 8
	maxRetry int = 5
)

var ErrMaxRetry = errors.New("generation attempts exceeded")

func AddShortURL(ctx context.Context, s data.Storager, originalURL string) (string, error) {
	for range maxRetry {
		shortURL, err := generateShortURL()
		if err != nil {
			return "", fmt.Errorf("failed to generate short URL: %w", err)
		}

		storeErr := s.StoreShortURL(ctx, shortURL, originalURL)

		if storeErr != nil {
			continue
		}

		return shortURL, nil
	}

	return "", fmt.Errorf("%w for original URL %s", ErrMaxRetry, originalURL)
}

func AddBatchShortURL(ctx context.Context, s data.Storager, req models.BatchRequest) (models.BatchResponse, error) {
	arrURLs := []models.URL{}
	resp := models.BatchResponse{}

	for _, reqData := range req {
		shortURL, err := generateShortURL()
		if err != nil {
			return models.BatchResponse{}, fmt.Errorf("failed to generate short URL: %w", err)
		}

		url := models.URL{
			ShortURL:    shortURL,
			OriginalURL: reqData.OriginalURL,
		}

		respData := models.BatchDataResponse{
			ShortURL:      shortURL,
			CorrelationID: reqData.CorrelationID,
		}

		arrURLs = append(arrURLs, url)
		resp = append(resp, respData)
	}

	if storeErr := s.StoreShortURLs(ctx, arrURLs); storeErr != nil {
		return models.BatchResponse{}, fmt.Errorf("failed to store short URLs: %w", storeErr)
	}

	return resp, nil
}

func generateShortURL() (string, error) {
	bytes := make([]byte, keyBytes)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate short URL error: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}

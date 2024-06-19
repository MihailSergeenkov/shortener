package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"

	"github.com/MihailSergeenkov/shortener/internal/app/config"
	"github.com/MihailSergeenkov/shortener/internal/app/data"
	"github.com/MihailSergeenkov/shortener/internal/app/models"
)

const keyBytes int = 8

func AddShortURL(ctx context.Context, s data.Storager, originalURL string) (string, error) {
	shortURL, err := generateShortURL()
	if err != nil {
		return "", fmt.Errorf("failed to generate short URL: %w", err)
	}

	storeErr := s.StoreShortURL(ctx, shortURL, originalURL)

	if storeErr != nil {
		return "", fmt.Errorf("failed to store short URL: %w", storeErr)
	}

	return shortURL, nil
}

func AddBatchShortURL(ctx context.Context, s data.Storager, req models.BatchRequest) (models.BatchResponse, error) {
	arrURLs := []models.URL{}
	resp := models.BatchResponse{}

	for _, reqData := range req {
		shortURL, err := generateShortURL()
		if err != nil {
			return models.BatchResponse{}, fmt.Errorf("failed to generate short URL: %w", err)
		}

		u := models.URL{
			ShortURL:    shortURL,
			OriginalURL: reqData.OriginalURL,
		}

		result, err := url.JoinPath(config.Params.BaseURL, shortURL)

		if err != nil {
			return models.BatchResponse{}, fmt.Errorf("failed to construct URL: %w", err)
		}

		respData := models.BatchDataResponse{
			ShortURL:      result,
			CorrelationID: reqData.CorrelationID,
		}

		arrURLs = append(arrURLs, u)
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

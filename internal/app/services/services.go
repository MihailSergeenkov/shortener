package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"

	"github.com/MihailSergeenkov/shortener/internal/app/config"
	"github.com/MihailSergeenkov/shortener/internal/app/constants"
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

	userID, ok := ctx.Value(constants.KeyUserID).(string)

	if !ok {
		return models.BatchResponse{}, fmt.Errorf("failed to fetch user id from context")
	}

	for _, reqData := range req {
		shortURL, err := generateShortURL()
		if err != nil {
			return models.BatchResponse{}, fmt.Errorf("failed to generate short URL: %w", err)
		}

		u := models.URL{
			ShortURL:    shortURL,
			OriginalURL: reqData.OriginalURL,
			UserID:      userID,
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

func FetchUserURLs(ctx context.Context, s data.Storager) (models.UserURLsResponse, error) {
	resp := models.UserURLsResponse{}

	urls, err := s.FetchUserURLs(ctx)

	if err != nil {
		return models.UserURLsResponse{}, fmt.Errorf("failed to fetch URLs: %w", err)
	}

	for _, u := range urls {
		result, err := url.JoinPath(config.Params.BaseURL, u.ShortURL)

		if err != nil {
			return models.UserURLsResponse{}, fmt.Errorf("failed to construct URL: %w", err)
		}

		respData := models.UserURLsDataResponse{
			ShortURL:    result,
			OriginalURL: u.OriginalURL,
		}

		resp = append(resp, respData)
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

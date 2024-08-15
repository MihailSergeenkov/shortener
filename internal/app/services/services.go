package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"path"

	"go.uber.org/zap"

	"github.com/MihailSergeenkov/shortener/internal/app/common"
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

	userID, ok := ctx.Value(common.KeyUserID).(string)
	if !ok {
		return models.BatchResponse{}, common.ErrFetchUserIDFromContext
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

		baseURL := config.Params.BaseURL
		newPath := path.Join(baseURL.Path, shortURL)
		baseURL.Path = newPath

		respData := models.BatchDataResponse{
			ShortURL:      baseURL.String(),
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
		baseURL := config.Params.BaseURL
		newPath := path.Join(baseURL.Path, u.ShortURL)
		baseURL.Path = newPath

		respData := models.UserURLsDataResponse{
			ShortURL:    baseURL.String(),
			OriginalURL: u.OriginalURL,
		}

		resp = append(resp, respData)
	}

	return resp, nil
}

func DeleteUserURLs(ctx context.Context, l *zap.Logger, s data.Storager, shortURLs []string) error {
	urls := make([]string, 0)

	inputCh := generator(ctx, shortURLs)
	checkResultCh := checkCh(ctx, l, s, inputCh)

	for url := range checkResultCh {
		urls = append(urls, url)
	}

	err := s.DeleteShortURLs(ctx, urls)
	if err != nil {
		return fmt.Errorf("failed to delete URLs: %w", err)
	}

	return nil
}

func generateShortURL() (string, error) {
	bytes := make([]byte, keyBytes)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate short URL error: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}

func generator(ctx context.Context, shortURLs []string) chan string {
	inputCh := make(chan string)

	go func() {
		defer close(inputCh)

		for _, shortURL := range shortURLs {
			select {
			case <-ctx.Done():
				return
			case inputCh <- shortURL:
			}
		}
	}()

	return inputCh
}

func checkCh(ctx context.Context, l *zap.Logger, s data.Storager, inputCh chan string) chan string {
	checkRes := make(chan string)

	go func() {
		defer close(checkRes)

		for url := range inputCh {
			u, err := checkURL(ctx, s, url)
			if err != nil {
				if errors.Is(err, common.ErrPermDenied) {
					continue
				}

				l.Error("failed to check URL", zap.Error(err))
				continue
			}

			select {
			case <-ctx.Done():
				return
			case checkRes <- u:
			}
		}
	}()
	return checkRes
}

func checkURL(ctx context.Context, s data.Storager, shortURL string) (string, error) {
	userID, ok := ctx.Value(common.KeyUserID).(string)
	if !ok {
		return "", common.ErrFetchUserIDFromContext
	}

	u, err := s.GetURL(ctx, shortURL)
	if err != nil {
		return "", fmt.Errorf("failed to get URL: %w", err)
	}

	if u.UserID != userID {
		return "", common.ErrPermDenied
	}

	return shortURL, nil
}

package storage

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
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

type Urls map[string]string

func Init() Urls {
	return make(Urls, initSize)
}

func (urls Urls) AddURL(u string) (string, error) {
	for i := 0; i < maxRetry; i++ {
		id, err := randomHex()
		if err != nil {
			return "", err
		}
		if _, ok := urls[id]; ok {
			continue
		}

		urls[id] = u
		return id, nil
	}

	return "", fmt.Errorf("%w for url %s", ErrMaxRetry, u)
}

func (urls Urls) FetchURL(id string) (string, error) {
	u, ok := urls[id]

	if !ok {
		return "", fmt.Errorf("%w for id %s", ErrURLNotFound, id)
	}

	return u, nil
}

func randomHex() (string, error) {
	bytes := make([]byte, keyBytes)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate key error: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}

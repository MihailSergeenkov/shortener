package storage

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
)

type Urls map[string]string

var ErrURLNotFound = errors.New("url not found")
var initSize = 100
var keyBytes = 8

func Init() Urls {
	return make(Urls, initSize)
}

func (urls Urls) AddURL(u string) (string, error) {
	for {
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

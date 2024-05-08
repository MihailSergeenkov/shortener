package storage

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
)

type Urls map[string]string

var ErrURLNotFound = errors.New("url not found")

func Init() Urls {
	return make(Urls, 100)
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
	bytes := make([]byte, 8)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

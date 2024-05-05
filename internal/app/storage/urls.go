package storage

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
)

var urls = make(map[string]string, 100)
var ErrURLNotFound = errors.New("url not found")

func AddURL(u string) (string, error) {
	h, err := randomHex()

	if err != nil {
		return "", err
	}

	urls[h] = u

	return h, nil
}

func FetchURL(h string) (string, error) {
	u, ok := urls[h]

	if !ok {
		return "", fmt.Errorf("%w for hash %s", ErrURLNotFound, h)
	}

	return u, nil
}

func randomHex() (string, error) {
	bytes := make([]byte, 4)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

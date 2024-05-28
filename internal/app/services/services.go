package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

const keyBytes int = 8

func GenerateShortURL() (string, error) {
	bytes := make([]byte, keyBytes)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate short URL error: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}

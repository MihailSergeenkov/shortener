package routes

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/MihailSergeenkov/shortener/internal/app/config"
	"github.com/MihailSergeenkov/shortener/internal/app/constants"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

const keyBytes int = 8

func setAuthMiddleware(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authCookie, cookieErr := r.Cookie("AUTH_TOKEN")

			if cookieErr != nil {
				if errors.Is(cookieErr, http.ErrNoCookie) {
					var err error
					authCookie, err = setAuthCookie(w)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						l.Error("failed to build auth token", zap.Error(err))
						return
					}
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					l.Error("failed to fetch cookie", zap.Error(cookieErr))
					return
				}
			}

			userID := getUserID(authCookie.Value)

			if userID == "" {
				authCookie, err := setAuthCookie(w)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					l.Error("failed to build auth token", zap.Error(err))
					return
				}
				userID = getUserID(authCookie.Value)
			}

			newContext := context.WithValue(r.Context(), constants.KeyUserID, userID)
			newRequest := r.WithContext(newContext)
			next.ServeHTTP(w, newRequest)
		})
	}
}

func checkAuthMiddleware(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authCookie, cookieErr := r.Cookie("AUTH_TOKEN")

			if cookieErr != nil {
				w.WriteHeader(http.StatusUnauthorized)
				l.Error("failed to fetch auth token", zap.Error(cookieErr))
				return
			}

			userID := getUserID(authCookie.Value)

			if userID == "" {
				w.WriteHeader(http.StatusUnauthorized)
				l.Error("failed to parse auth token")
				return
			}

			newContext := context.WithValue(r.Context(), constants.KeyUserID, userID)
			newRequest := r.WithContext(newContext)
			next.ServeHTTP(w, newRequest)
		})
	}
}

func setAuthCookie(w http.ResponseWriter) (*http.Cookie, error) {
	authToken, err := buildJWTString()
	if err != nil {
		return nil, fmt.Errorf("failed to build auth token: %w", err)
	}

	cookie := &http.Cookie{
		Name:     "AUTH_TOKEN",
		Value:    authToken,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	return cookie, nil
}

func buildJWTString() (string, error) {
	userID, err := generateUserID()
	if err != nil {
		return "", fmt.Errorf("failed to generate user id: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(config.Params.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func generateUserID() (string, error) {
	bytes := make([]byte, keyBytes)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate user id error: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}

func getUserID(tokenString string) string {
	claims := &Claims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Params.SecretKey), nil
	})

	if err != nil {
		return ""
	}

	return claims.UserID
}

package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSetAuthMiddleware(t *testing.T) {
	logger := zap.NewNop()
	someHandler := func(w http.ResponseWriter, r *http.Request) {}

	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	w := httptest.NewRecorder()
	m := setAuthMiddleware(logger)(http.HandlerFunc(someHandler))
	m.ServeHTTP(w, request)

	res := w.Result()
	defer closeBody(t, res)

	assert.Equal(t, 200, res.StatusCode)
}

func TestCheckAuthMiddleware_OK(t *testing.T) {
	logger := zap.NewNop()
	authToken, err := buildJWTString()
	require.NoError(t, err)

	someHandler := func(w http.ResponseWriter, r *http.Request) {}

	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	request.AddCookie(&http.Cookie{
		Name:     "AUTH_TOKEN",
		Value:    authToken,
		HttpOnly: true,
	})
	w := httptest.NewRecorder()
	m := checkAuthMiddleware(logger)(http.HandlerFunc(someHandler))
	m.ServeHTTP(w, request)

	res := w.Result()
	defer closeBody(t, res)

	assert.Equal(t, 200, res.StatusCode)
}

func TestCheckAuthMiddleware_Failed(t *testing.T) {
	logger := zap.NewNop()
	someHandler := func(w http.ResponseWriter, r *http.Request) {}

	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	w := httptest.NewRecorder()
	m := checkAuthMiddleware(logger)(http.HandlerFunc(someHandler))
	m.ServeHTTP(w, request)

	res := w.Result()
	defer closeBody(t, res)

	assert.Equal(t, 401, res.StatusCode)
}

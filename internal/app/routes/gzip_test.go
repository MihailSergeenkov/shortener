package routes

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGzipMiddleware(t *testing.T) {
	logger := zap.NewNop()
	someHandler := func(w http.ResponseWriter, r *http.Request) {}

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write([]byte("some data"))
	require.NoError(t, err)
	err = zw.Close()
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodGet, "/", &buf)
	request.Header.Add("Accept-Encoding", "gzip")
	request.Header.Add("Content-Encoding", "gzip")
	request.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	m := gzipMiddleware(logger)(http.HandlerFunc(someHandler))
	m.ServeHTTP(w, request)

	res := w.Result()
	defer closeBody(t, res)

	assert.Equal(t, 200, res.StatusCode)
}

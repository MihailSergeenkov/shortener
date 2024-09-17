package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWithRequestLogging(t *testing.T) {
	logger := zap.NewNop()
	someHandler := func(w http.ResponseWriter, r *http.Request) {}

	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	w := httptest.NewRecorder()
	m := withRequestLogging(logger)(http.HandlerFunc(someHandler))
	m.ServeHTTP(w, request)

	res := w.Result()
	defer closeBody(t, res)

	assert.Equal(t, 200, res.StatusCode)
}

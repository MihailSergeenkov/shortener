package test_helpers

import (
	"net/http"
	"testing"
)

func CloseBody(t *testing.T, r *http.Response) {
	t.Helper()
	err := r.Body.Close()

	if err != nil {
		t.Log(err)
	}
}

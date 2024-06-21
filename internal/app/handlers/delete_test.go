package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAPIDeleteUserURLsHandler(t *testing.T) {
	logger := zap.NewNop()
	storage := MockStorage{}

	type request struct {
		body string
	}

	type want struct {
		code int
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "success add url",
			request: request{
				body: `["6qxTVvsy", "RTfd56hn", "Jlfd67ds"]`,
			},
			want: want{
				code: http.StatusAccepted,
			},
		},
		{
			name: "bad request",
			request: request{
				body: `sdfsdfsdfsdf`,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(test.request.body))
			w := httptest.NewRecorder()
			APIDeleteUserURLsHandler(logger, &storage)(w, request)

			res := w.Result()
			defer func() {
				err := res.Body.Close()

				if err != nil {
					t.Log(err)
				}
			}()

			assert.Equal(t, test.want.code, res.StatusCode)

			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
		})
	}
}

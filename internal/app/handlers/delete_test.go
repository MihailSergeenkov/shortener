package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MihailSergeenkov/shortener/internal/app/data/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAPIDeleteUserURLsHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zap.NewNop()
	storage := mock.NewMockStorager(mockCtrl)

	body := `["6qxTVvsy"]`

	t.Run("success delete url", func(t *testing.T) {
		storage.EXPECT().DeleteShortURLs(gomock.Any(), []string{}).Times(1).Return(nil)

		request := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(body))
		w := httptest.NewRecorder()
		APIDeleteUserURLsHandler(logger, storage)(w, request)

		res := w.Result()
		defer closeBody(t, res)

		assert.Equal(t, http.StatusAccepted, res.StatusCode)
		_, err := io.ReadAll(res.Body)
		require.NoError(t, err)
	})
}

func TestAPIDeleteUserURLsHandler_Failed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zap.NewNop()
	storage := mock.NewMockStorager(mockCtrl)

	body := `["6qxTVvsy"]`

	type request struct {
		body string
	}
	type want struct {
		err  error
		code int
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "failed delete url",
			request: request{
				body: body,
			},
			want: want{
				err:  errors.New("some error"),
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "bad request",
			request: request{
				body: `sdfsdfsdfsdf`,
			},
			want: want{
				err:  nil,
				code: http.StatusBadRequest,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.want.err != nil {
				storage.EXPECT().DeleteShortURLs(gomock.Any(), []string{}).Times(1).Return(test.want.err)
			}

			request := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(test.request.body))
			w := httptest.NewRecorder()
			APIDeleteUserURLsHandler(logger, storage)(w, request)

			res := w.Result()
			defer closeBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
		})
	}
}

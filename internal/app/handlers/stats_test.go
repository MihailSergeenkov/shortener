package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MihailSergeenkov/shortener/internal/app/data/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAPIFetchStatsHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zap.NewNop()
	storage := mock.NewMockStorager(mockCtrl)

	type storeResp struct {
		urls  int
		users int
		err   error
	}
	type want struct {
		code int
	}
	tests := []struct {
		name      string
		storeResp storeResp
		want      want
	}{
		{
			name: "success fetch stats",
			storeResp: storeResp{
				urls:  1,
				users: 1,
				err:   nil,
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "when fetch failed",
			storeResp: storeResp{
				urls:  0,
				users: 0,
				err:   errors.New("some error"),
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.EXPECT().FetchStats(gomock.Any()).Times(1).
				Return(test.storeResp.urls, test.storeResp.users, test.storeResp.err)

			request := httptest.NewRequest(http.MethodGet, "/api/internal/stats", http.NoBody)
			w := httptest.NewRecorder()
			APIFetchStatsHandler(logger, storage)(w, request)

			res := w.Result()
			defer closeBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)

			resBody, err := io.ReadAll(res.Body)

			if test.want.code == http.StatusOK {
				require.NoError(t, err)
				assert.NotEmpty(t, resBody)
			}
		})
	}
}

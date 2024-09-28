package routes

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCheckSubnetMiddleware(t *testing.T) {
	logger := zap.NewNop()

	_, trustedSubnet, err := net.ParseCIDR("192.168.1.0/24")
	require.NoError(t, err)

	type want struct {
		code int
	}
	tests := []struct {
		name          string
		trustedSubnet *net.IPNet
		clientIP      string
		want          want
	}{
		{
			name:          "success check",
			trustedSubnet: trustedSubnet,
			clientIP:      "192.168.1.33",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:          "without trusted subnet",
			trustedSubnet: nil,
			clientIP:      "192.168.1.33",
			want: want{
				code: http.StatusForbidden,
			},
		},
		{
			name:          "failed parse client IP",
			trustedSubnet: trustedSubnet,
			clientIP:      "some string",
			want: want{
				code: http.StatusInternalServerError,
			},
		},
		{
			name:          "client IP not from trusted subnet",
			trustedSubnet: trustedSubnet,
			clientIP:      "192.168.100.33",
			want: want{
				code: http.StatusForbidden,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			someHandler := func(w http.ResponseWriter, r *http.Request) {}

			request := httptest.NewRequest(http.MethodGet, "/api/internal/stats", http.NoBody)
			request.Header.Add("X-Real-IP", test.clientIP)

			w := httptest.NewRecorder()

			m := checkSubnetMiddleware(logger, test.trustedSubnet)(http.HandlerFunc(someHandler))
			m.ServeHTTP(w, request)

			res := w.Result()
			defer closeBody(t, res)

			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}

package main

import (
	"errors"
	"net/http"
	"testing"

	webmock "github.com/MihailSergeenkov/shortener/cmd/shortener/mock"

	"github.com/MihailSergeenkov/shortener/internal/app/data/mock"
	"github.com/MihailSergeenkov/shortener/internal/app/routes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestConfigureServer(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zap.NewNop()
	storage := mock.NewMockStorager(mockCtrl)
	router := routes.NewRouter(logger, storage)
	runAddr := "localhost:8080"

	tests := []struct {
		name        string
		enableHTTPS bool
	}{
		{
			name:        "HTTP",
			enableHTTPS: false,
		},
		{
			name:        "HTTPS",
			enableHTTPS: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := configureServer(router, test.enableHTTPS, runAddr)

			assert.IsType(t, (*http.Server)(nil), server)
			assert.Equal(t, runAddr, server.Addr)
			assert.Equal(t, router, server.Handler)

			if test.enableHTTPS {
				assert.NotEmpty(t, server.TLSConfig)
			} else {
				assert.Empty(t, server.TLSConfig)
			}
		})
	}
}

func TestRunServer_OK(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	server := webmock.NewMockWebServer(mockCtrl)

	tests := []struct {
		name        string
		enableHTTPS bool
	}{
		{
			name:        "run HTTP server",
			enableHTTPS: false,
		},
		{
			name:        "run HTTPS server",
			enableHTTPS: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.enableHTTPS {
				server.EXPECT().ListenAndServeTLS("", "").Times(1).Return(nil)
			} else {
				server.EXPECT().ListenAndServe().Times(1).Return(nil)
			}

			err := runServer(server, test.enableHTTPS)
			require.NoError(t, err)
		})
	}
}

func TestRunServer_Failed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	server := webmock.NewMockWebServer(mockCtrl)
	errSome := errors.New("some error")
	server.EXPECT().ListenAndServe().Times(1).Return(errSome)

	t.Run("failed run server", func(t *testing.T) {
		err := runServer(server, false)
		require.Error(t, err)
		require.ErrorContains(t, err, "listen and server has failed")
	})
}

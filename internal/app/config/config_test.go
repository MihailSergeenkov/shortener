package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetConfigData(t *testing.T) {
	tests := []struct {
		setEnv      func()
		name        string
		wantErr     bool
		presentData bool
	}{
		{
			name: "config present",
			setEnv: func() {
				require.NoError(t, os.Setenv("CONFIG", "testdata/settings.json"))
			},
			wantErr:     false,
			presentData: true,
		},
		{
			name: "config absent",
			setEnv: func() {
				require.NoError(t, os.Setenv("CONFIG", ""))
			},
			wantErr:     false,
			presentData: false,
		},
		{
			name: "config failed",
			setEnv: func() {
				require.NoError(t, os.Setenv("CONFIG", "set.json"))
			},
			wantErr:     true,
			presentData: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setEnv()

			_, presentData, err := getConfigData()

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, "failed to read config file")
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.presentData, presentData)
			}
		})
	}
}

func TestParseConfigData(t *testing.T) {
	tests := []struct {
		name       string
		wantErr    bool
		configFile string
	}{
		{
			name:       "success parsed config",
			wantErr:    false,
			configFile: "testdata/settings.json",
		},
		{
			name:       "bad config",
			wantErr:    true,
			configFile: "testdata/bad_settings.json",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := os.ReadFile(test.configFile)
			require.NoError(t, err)

			err = parseConfigData(data)

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, "failed to unmarshal json")
			} else {
				require.NoError(t, err)
				assert.Equal(t, "http://localhost:8081", os.Getenv("BASE_URL"))
				assert.Equal(t, "localhost:8081", os.Getenv("SERVER_ADDRESS"))
				assert.Equal(t, "/tmp/url-db.json", os.Getenv("FILE_STORAGE_PATH"))
				assert.Equal(t, "", os.Getenv("DATABASE_DSN"))
				assert.Equal(t, "12345", os.Getenv("SECRET_KEY"))
				assert.Equal(t, "1m", os.Getenv("DROP_URLS_PERIOD"))
				assert.Equal(t, "ERROR", os.Getenv("LOG_LEVEL"))
				assert.Equal(t, "false", os.Getenv("ENABLE_HTTPS"))
			}
		})
	}
}

func TestParseEnv(t *testing.T) {
	tests := []struct {
		setEnv  func()
		name    string
		wantErr bool
	}{
		{
			name: "valid config",
			setEnv: func() {
				require.NoError(t, os.Setenv("SERVER_ADDRESS", "localhost:8081"))
				require.NoError(t, os.Setenv("BASE_URL", "http://localhost:8081"))
			},
			wantErr: false,
		},
		{
			name: "invalid config",
			setEnv: func() {
				require.NoError(t, os.Setenv("SERVER_ADDRESS", "localhost:8081"))
				require.NoError(t, os.Setenv("LOG_LEVEL", "some string"))
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setEnv()

			err := Params.parseEnv()

			if test.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, "env error")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

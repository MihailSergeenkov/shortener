package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseEnv(t *testing.T) {
	tests := []struct {
		name    string
		setEnv  func()
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

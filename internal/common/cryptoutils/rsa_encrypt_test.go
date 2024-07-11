package cryptoutils

import (
	"testing"

	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
)

func TestEncryptRSA(t *testing.T) {
	type Config struct {
		Crypto config.Crypto `mapstructure:"crypto"`
	}

	configFile, err := config.FilePath()
	if err != nil {
		t.Errorf("failed to get config file path: %v", err)
		return
	}

	cfg := Config{}
	err = config.Load(configFile, &cfg)
	if err != nil {
		t.Errorf("failed to load config: %v", err)
		return
	}

	type args struct {
		data []byte
		key  []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "generate encrypted key",
			args: args{
				data: []byte("07e76313-d119-4ba5-9a3e-d90f71c4c001"),
				key:  cfg.Crypto.PublicKeyBytes(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncryptRSA(tt.args.data, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptRSA() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("EncryptRSA() got = %v", got)
		})
	}
}

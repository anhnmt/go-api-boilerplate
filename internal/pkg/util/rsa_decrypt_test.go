package util

import (
	"reflect"
	"testing"

	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
)

func TestDecryptRSA(t *testing.T) {
	type Config struct {
		Crypto config.Crypto `mapstructure:"crypto"`
	}

	cfg := Config{}
	err := config.Load(&cfg)
	if err != nil {
		t.Errorf("failed to load config: %v", err)
		return
	}

	type args struct {
		data       string
		privateKey []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "decrypt rsa success",
			args: args{
				data:       "McjPiQnUd1F2GUsGGAToMuH+gKC6TuttFpBC1EL+Smk4I5syICbryUPTNu0N0Q16ICrAIZWzDLxP+xyaR5pjLIFTqlHV79sVmEGpaoD+syAHMUTw4LnDccqrDXCnEDu1fkUQlzIDEasiP2nYiaqY0cKDlCIMPH1pwX2+Mb5Cl5s=",
				privateKey: cfg.Crypto.PrivateKeyBytes(),
			},
			wantErr: false,
			want:    []byte("07e76313-d119-4ba5-9a3e-d90f71c4c001"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecryptRSA(tt.args.data, tt.args.privateKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptRSA() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecryptRSA() got = %v, want %v", got, tt.want)
			}
		})
	}
}

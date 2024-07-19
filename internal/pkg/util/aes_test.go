package util

import (
	"reflect"
	"testing"

	"github.com/bytedance/sonic"

	"github.com/anhnmt/go-api-boilerplate/gen/pb"
)

func TestDecryptAES(t *testing.T) {
	type args struct {
		data string
		key  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "decrypt AES success",
			args: args{
				data: "QbzS4pS4RJZC8QYgyIIPqDovMEmrJm8OAhjBiESiSyU=",
				key:  "07e76313-d119-4ba5-9a3e-d90f71c4c001",
			},
			wantErr: false,
			want:    `{"data":"ahihi"}`,
		},
		{
			name: "decrypt AES success 2",
			args: args{
				data: "U/p3pHW9DPzidwnSMcYXMQ==",
				key:  "07e76313-d119-4ba5-9a3e-d90f71c4c001",
			},
			wantErr: false,
			want:    `{"data":"alo"}`,
		},
		{
			name: "error unexpected EOF",
			args: args{
				data: "SLZtrykSmAbgzol7aJj9+7vQoHcI/sSFxeF9VnG6h0ebPsf0qPPVNe2zttb+Iasp",
				key:  "07e76313-d119-4ba5-9a3e-d90f71c4c001",
			},
			wantErr: false,
			want:    `{"code":"unknown","message":"unexpected EOF"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecryptAES(tt.args.data, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptAES() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotStr := string(got)
			if !reflect.DeepEqual(gotStr, tt.want) {
				t.Errorf("DecryptAES() got = %v, want %v", gotStr, tt.want)
			}
		})
	}
}

func TestEncryptAES(t *testing.T) {
	marshal, err := sonic.Marshal(&pb.EncryptRequest{
		Data: "ahihi",
	})
	if err != nil {
		t.Errorf("failed to marshal: %v", err)
		return
	}

	type args struct {
		plaintext []byte
		key       string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "encrypt AES success",
			args: args{
				plaintext: marshal,
				key:       "07e76313-d119-4ba5-9a3e-d90f71c4c001",
			},
			wantErr: false,
			want:    "QbzS4pS4RJZC8QYgyIIPqDovMEmrJm8OAhjBiESiSyU=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncryptAES(tt.args.plaintext, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptAES() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EncryptAES() got = %v, want %v", got, tt.want)
			}
		})
	}
}

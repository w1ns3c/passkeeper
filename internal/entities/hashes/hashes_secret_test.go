package hashes

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"testing"

	"passkeeper/internal/entities/config"
)

func TestEncryptDecryptSecret(t *testing.T) {
	type args struct {
		secret string
		key    string
	}

	secret1, err := GenerateSecret(config.UserSecretLen)
	if err != nil {
		fmt.Println(err)
		return
	}
	secret2, err := GenerateSecret(config.UserSecretLen)
	if err != nil {
		fmt.Println(err)
		return
	}
	secret3, err := GenerateSecret(config.UserSecretLen)
	if err != nil {
		fmt.Println(err)
		return
	}

	key := sha512.Sum512([]byte("12312"))
	key1 := hex.EncodeToString(key[:])

	key = sha512.Sum512([]byte("supersecret"))
	key2 := hex.EncodeToString(key[:])

	key = sha512.Sum512([]byte("incredible P@ssw0rd "))
	key3 := hex.EncodeToString(key[:])

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Simple1",
			args: args{
				secret: secret1,
				key:    key1,
			},
			wantErr: false,
		},
		{
			name: "Simple2",
			args: args{
				secret: secret2,
				key:    key2,
			},
			wantErr: false,
		},
		{
			name: "Simple3",
			args: args{
				secret: secret3,
				key:    key3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCryptSecret, err := EncryptSecret(tt.args.secret, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotSecret, err := DecryptSecret(gotCryptSecret, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSecret != tt.args.secret {
				t.Errorf("Encrypt/DecryptSecret() not work input secret = %s, %d,  got secret %s", tt.args.secret, len(tt.args.secret), gotSecret)
			}
		})
	}
}

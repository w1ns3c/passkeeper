package hashes

import (
	"testing"
	"time"

	"passkeeper/internal/entities/config"
)

func TestCheckToken(t *testing.T) {
	type args struct {
		token  string
		secret string
	}

	var (
		id1        = "user1"
		secret1, _ = GenerateSecret(config.UserSecretLen)
		token1, _  = GenerateToken(id1, secret1, time.Second*200)
	)

	tests := []struct {
		name       string
		args       args
		wantUserID string
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "Test 1",
			args: args{
				token:  token1,
				secret: secret1,
			},
			wantUserID: id1,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := CheckToken(tt.args.token, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("CheckToken() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestExtractUserID(t *testing.T) {
	var (
		id1        = "user1"
		secret1, _ = GenerateSecret(config.UserSecretLen)
		token1, _  = GenerateToken(id1, secret1, time.Second*200)
	)

	tests := []struct {
		name       string
		token      string
		wantUserID string
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name:       "Test 1",
			token:      token1,
			wantUserID: id1,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ExtractUserID(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ExtractUserID() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

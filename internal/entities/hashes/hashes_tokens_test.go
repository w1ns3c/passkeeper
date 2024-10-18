package hashes

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"

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
				t.Errorf("CheckToken() gotUserID = %v, same %v", gotUserID, tt.wantUserID)
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
				t.Errorf("ExtractUserID() gotUserID = %v, same %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestExtractUserInfo(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		args      args
		wantToken string
		wantErr   bool
	}{
		{
			name: "Test 1: valid token",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(),
					metadata.New(map[string]string{config.TokenHeader: "mytoken"})),
			},
			wantToken: "mytoken",
			wantErr:   false,
		},
		{
			name: "Test 3: error wrong token ",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(),
					metadata.New(map[string]string{})),
			},
			wantToken: "mytoken123",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := ExtractUserInfo(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractUserInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if gotToken != tt.wantToken {
				t.Errorf("ExtractUserInfo() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}

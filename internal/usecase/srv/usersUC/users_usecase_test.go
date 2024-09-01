package usersUC

import (
	"passkeeper/internal/entities/hashes"
	"testing"
)

func TestGenerateCryptoHash(t *testing.T) {
	type args struct {
		password string
		salt     string
	}
	tests := []struct {
		name     string
		args     args
		wantHash string
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name: "simple test with empty salt",
			args: args{
				password: "123123",
				salt:     "",
			},
		},
		{
			name: "simple test",
			args: args{
				password: "123123",
				salt:     "supersalt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHash, err := hashes.GenerateCryptoHash(tt.args.password, tt.args.salt)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateCryptoHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Errorf("GenerateCryptoHash() gotHash = %v, want %v", gotHash, tt.wantHash)

			//if gotHash != tt.wantHash {
			//	t.Errorf("GenerateCryptoHash() gotHash = %v, want %v", gotHash, tt.wantHash)
			//}
		})
	}
}

func TestComparePassAndCryptoHash(t *testing.T) {
	type args struct {
		password string
		hash     string
		salt     string
	}
	tests := []struct {
		name  string
		args  args
		equal bool
	}{
		// TODO: Add test cases.
		{
			name: "simple test with empty salt",
			args: args{
				password: "123123",
				salt:     "",
				hash:     "$2a$10$HJj3MwoxaFLAighnqNWLU.9JHbYZTMrjexCw4/zYKANuilggdm9HC",
			},
			equal: true,
		},
		{
			name: "simple test",
			args: args{
				password: "123123",
				salt:     "supersalt",
				hash:     "$2a$10$17Flv.5YK4c1DDHPOLrrYuL7RSRLzOnC6WOqW1wMgedGDcEo24eZ2",
			},
			equal: true,
		},
		{
			name: "wrong password",
			args: args{
				password: "1231231",
				salt:     "supersalt",
				hash:     "$2a$10$17Flv.5YK4c1DDHPOLrrYuL7RSRLzOnC6WOqW1wMgedGDcEo24eZ2",
			},
			equal: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hashes.ComparePassAndCryptoHash(tt.args.password, tt.args.hash, tt.args.salt); got != tt.equal {
				t.Errorf("ComparePassAndCryptoHash() = %v, want %v", got, tt.equal)
			}
		})
	}
}

package hashes

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"testing"
)

func TestComparePassAndCryptoHash(t *testing.T) {
	type args struct {
		password string
		hash     string
		salt     string
	}
	tests := []struct {
		name string
		args args
		same bool
	}{
		{
			name: "TestComparePassAndCryptoHash",
			args: args{
				password: "password",
				hash:     "password",
				salt:     "salt",
			},
			same: false,
		},
		{
			name: "TestComparePassAndCryptoHash",
			args: args{
				password: "fa585d89c851dd338a70dcf535aa2a92fee7836dd6aff1226583e88e0996293f16bc009c652826e0fc5c706695a03cddce372f139eff4d13959da6f1f5d3eabe",
				hash:     "$2a$10$Bltccu96qhJX/FnkixAZ3eeK6Zw/JUxUHz5WXZefG68GgY733hoIS",
				salt:     "86961cd1a7152bc9d747be71dee7946e",
			},
			same: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComparePassAndCryptoHash(tt.args.password, tt.args.hash, tt.args.salt); got != tt.same {
				t.Errorf("ComparePassAndCryptoHash() = %v, same %v", got, tt.same)
			}
		})
	}
}

func TestGenerateCryptoHash(t *testing.T) {
	type args struct {
		password string
		salt     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestGenerateCryptoHash",
			args: args{
				password: "password",
				salt:     "salt",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHash, err := GenerateCryptoHash(tt.args.password, tt.args.salt)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateCryptoHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !ComparePassAndCryptoHash(tt.args.password, gotHash, tt.args.salt) {
				t.Errorf("GenerateCryptoHash() gotHash = %v, not for pass %v", gotHash, tt.args.password)
			}
		})
	}
}

func TestGenerateHash(t *testing.T) {
	type args struct {
		password string
		salt     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestGenerateHash",
			args: args{
				password: "password",
				salt:     "salt",
			},
			want: "f9203a341650f7e335b16b3118c8cecd8a3e399bd9659ae2b075ab7e22e95204",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateHash(tt.args.password, tt.args.salt); got != tt.want {
				t.Errorf("GenerateHash() = %v, same %v", got, tt.want)
			}
		})
	}
}

func TestGenerateUserID(t *testing.T) {
	type args struct {
		secret string
		salt   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestGenerateUserID",
			args: args{
				secret: "secret",
				salt:   "salt",
			},
			want: "e929573f254e2cf8624772896cbd4cd2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateUserID(tt.args.secret, tt.args.salt); got != tt.want {
				t.Errorf("GenerateUserID() = %v, same %v", got, tt.want)
			}
		})
	}
}

func Test_genID(t *testing.T) {
	type args struct {
		secret string
		salt   string
		h      hash.Hash
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test_genID sha512",
			args: args{
				secret: "secret",
				salt:   "salt",
				h:      sha512.New(),
			},
			want: "a4abd33735934d93792a9e5a9dba088c5a15176943c885da23d455fbf30c89e1a4633baac7a23c52af0dffba6fde8e84f4a332a5fbf2df24304528520bd5833c",
		},
		{
			name: "Test_genID sha256",
			args: args{
				secret: "secret",
				salt:   "salt",
				h:      sha256.New(),
			},
			want: "a1fe12ea2744e1ec94576be88e9e10513a4cfb11f08dac007c07282cd8913f27",
		},
		{
			name: "Test_genID md5",
			args: args{
				secret: "secret",
				salt:   "salt",
				h:      md5.New(),
			},
			want: "e929573f254e2cf8624772896cbd4cd2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := genID(tt.args.secret, tt.args.salt, tt.args.h)
			if got != tt.want {
				t.Errorf("genID() = %v, same %v", got, tt.want)
			}
		})
	}
}

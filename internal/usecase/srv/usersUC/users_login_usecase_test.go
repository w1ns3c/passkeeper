package usersUC

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/w1ns3c/go-examples/crypto"

	"passkeeper/internal/entities/config"
	"passkeeper/internal/entities/hashes"
	"passkeeper/internal/entities/myerrors"
	"passkeeper/internal/entities/structs"
	"passkeeper/internal/storage"
	mocks "passkeeper/mock"
)

func TestUserUsecase_LoginUser(t *testing.T) {
	type fields struct {
		ctx             context.Context
		storage         storage.UserStorage
		tokenLifeTime   time.Duration
		userSecretLen   int
		userPassSaltLen int
		log             *zerolog.Logger
	}
	type args struct {
		login    string
		password string
	}

	var (
		salt1, _  = crypto.GenRandStr(config.UserPassSaltLen)
		login1    = "login1"
		password1 = "password"
		password2 = "password22"
		hash, _   = hashes.GenerateCryptoHash(password1, salt1)
		//hash2, _        = hashes.GenerateCryptoHash(password2, salt1)
		userID1         = hashes.GenerateUserID(password1, salt1)
		secret1, _      = hashes.GenerateSecret(config.UserSecretLen)
		secureSecret, _ = hashes.EncryptSecret(secret1, password1)
		user1           = &structs.User{
			ID:     userID1,
			Login:  login1,
			Secret: secureSecret,
			Salt:   salt1,
			Hash:   hash,
		}
		user3 = &structs.User{
			ID:     userID1,
			Login:  login1,
			Secret: secureSecret,
			Salt:   salt1,
			Hash:   hash,
		}
		login2 = "login2"
	)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ctx := context.Background()

	mockStorage := mocks.NewMockUserStorage(mockCtrl)
	mockStorage.EXPECT().GetUserByLogin(ctx, login1).Return(user1, nil).Times(1)
	mockStorage.EXPECT().GetUserByLogin(ctx, login2).Return(nil, myerrors.ErrUserNotFound).Times(1)
	mockStorage.EXPECT().GetUserByLogin(ctx, login1).Return(user3, nil).Times(1)

	tests := []struct {
		name       string
		fields     fields
		args       args
		wantUserID string
		wantSecret string
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "Test 1: valid user",
			fields: fields{
				ctx:             context.Background(),
				storage:         mockStorage,
				tokenLifeTime:   time.Second * 100,
				userSecretLen:   config.UserSecretLen,
				userPassSaltLen: config.UserPassSaltLen,
				log:             new(zerolog.Logger),
			},
			args: args{
				login:    login1,
				password: password1,
			},
			wantUserID: user1.ID,
			wantSecret: user1.Secret,
			wantErr:    false,
		},
		{
			name: "Test 2: user not exist",
			fields: fields{
				ctx:             context.Background(),
				storage:         mockStorage,
				tokenLifeTime:   time.Second * 100,
				userSecretLen:   config.UserSecretLen,
				userPassSaltLen: config.UserPassSaltLen,
				log:             new(zerolog.Logger),
			},
			args: args{
				login:    login2,
				password: password1,
			},
			wantUserID: "",
			wantSecret: "",
			wantErr:    true,
		},
		{
			name: "Test 3: user with invalid password",
			fields: fields{
				ctx:             context.Background(),
				storage:         mockStorage,
				tokenLifeTime:   time.Second * 100,
				userSecretLen:   config.UserSecretLen,
				userPassSaltLen: config.UserPassSaltLen,
				log:             new(zerolog.Logger),
			},
			args: args{
				login:    login1,
				password: password2,
			},
			wantUserID: "",
			wantSecret: "",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserUsecase{
				ctx:             tt.fields.ctx,
				storage:         tt.fields.storage,
				tokenLifeTime:   tt.fields.tokenLifeTime,
				userSecretLen:   tt.fields.userSecretLen,
				userPassSaltLen: tt.fields.userPassSaltLen,
				log:             tt.fields.log,
			}
			gotToken, gotSecret, err := u.LoginUser(context.Background(), tt.args.login, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			userID, err := hashes.ExtractUserID(gotToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginUser() gotUser from token = %v, want %v", userID, tt.wantUserID)
			}
			require.Equal(t, tt.wantUserID, userID)

			if gotSecret != tt.wantSecret {
				t.Errorf("LoginUser() gotSecret = %v, want %v", gotSecret, tt.wantSecret)
			}
		})
	}
}

func TestUserUsecase_GetTokenSalt(t *testing.T) {
	type fields struct {
		ctx             context.Context
		storage         storage.UserStorage
		tokenLifeTime   time.Duration
		userSecretLen   int
		userPassSaltLen int
		log             *zerolog.Logger
	}

	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorage := mocks.NewMockUserStorage(mockCtrl)

	var (
		salt1, _   = crypto.GenRandStr(config.UserPassSaltLen)
		login1     = "login1"
		password   = "password"
		hash, _    = hashes.GenerateCryptoHash(password, salt1)
		userID1    = hashes.GenerateUserID(password, salt1)
		secret1, _ = hashes.GenerateSecret(config.UserSecretLen)
		user1      = &structs.User{
			ID:     userID1,
			Login:  login1,
			Secret: secret1,
			Salt:   salt1,
			Hash:   hash,
		}
		userID2 = hashes.GenerateUserID(password, salt1+"something else")
	)

	mockStorage.EXPECT().GetUserByID(ctx, userID1).Return(user1, nil)
	mockStorage.EXPECT().GetUserByID(ctx, userID2).Return(nil, myerrors.ErrUserNotFound)

	tests := []struct {
		name   string
		fields fields
		userID string
		want   string
	}{
		{
			name: "Test 1: user exist",
			fields: fields{
				ctx:             context.Background(),
				storage:         mockStorage,
				tokenLifeTime:   time.Second * 100,
				userSecretLen:   config.UserSecretLen,
				userPassSaltLen: config.UserPassSaltLen,
				log:             new(zerolog.Logger),
			},
			userID: userID1,
			want:   salt1,
		},

		{
			name: "Test 2: user not exist",
			fields: fields{
				ctx:             context.Background(),
				storage:         mockStorage,
				tokenLifeTime:   time.Second * 100,
				userSecretLen:   config.UserSecretLen,
				userPassSaltLen: config.UserPassSaltLen,
				log:             new(zerolog.Logger),
			},
			userID: userID2,
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserUsecase{
				ctx:             tt.fields.ctx,
				storage:         tt.fields.storage,
				tokenLifeTime:   tt.fields.tokenLifeTime,
				userSecretLen:   tt.fields.userSecretLen,
				userPassSaltLen: tt.fields.userPassSaltLen,
				log:             tt.fields.log,
			}
			if got := u.GetTokenSalt(ctx, tt.userID); got != tt.want {
				t.Errorf("GetTokenSalt() = %v, want %v", got, tt.want)
			}
		})
	}
}

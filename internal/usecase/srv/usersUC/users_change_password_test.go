package usersUC

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/w1ns3c/go-examples/crypto"

	"passkeeper/internal/entities/config"
	"passkeeper/internal/entities/hashes"
	"passkeeper/internal/entities/myerrors"
	"passkeeper/internal/entities/structs"
	"passkeeper/internal/storage"
	mocks "passkeeper/mock"
)

func TestUserUsecase_ChangePassword(t *testing.T) {
	type fields struct {
		ctx             context.Context
		storage         storage.UserStorage
		tokenLifeTime   time.Duration
		userSecretLen   int
		userPassSaltLen int
		log             *zerolog.Logger
	}
	type args struct {
		ctx       context.Context
		userID    string
		oldPass   string
		newPass   string
		reNewPass string
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockUserStorage(ctrl)

	var (
		ctx = context.Background()
		u   = &UserUsecase{
			ctx:             ctx,
			tokenLifeTime:   time.Hour,
			userSecretLen:   config.UserSecretLen,
			userPassSaltLen: config.UserPassSaltLen,
			log:             &zerolog.Logger{},
			storage:         mockStorage,
		}

		old_password = "old_password"

		salt1, _ = crypto.GenRandStr(config.UserPassSaltLen)
		salt2, _ = crypto.GenRandStr(config.UserPassSaltLen)
		salt3, _ = crypto.GenRandStr(config.UserPassSaltLen)

		login1            = "login1"
		valid_password    = "valid_password"
		valid_re_password = "valid_password"
		hash1, _          = hashes.GenerateCryptoHash(old_password, salt1)

		userID1    = hashes.GenerateUserID(old_password, salt1)
		secret1, _ = hashes.GenerateSecret(config.UserSecretLen)

		user1 = &structs.User{
			ID:     userID1,
			Login:  login1,
			Secret: secret1,
			Salt:   salt1,
			Hash:   hash1,
		}

		login2              = "login2"
		invalid_password    = "invalid_password"
		invalid_re_password = "invalid_password"
		userID2             = hashes.GenerateUserID(old_password, salt2)

		login3               = "login3"
		not_same_password    = "pass"
		not_same_re_password = "not_same_password"
		userID3              = hashes.GenerateUserID(old_password, salt3)

		hash2, _   = hashes.GenerateCryptoHash(invalid_password, salt2)
		hash3, _   = hashes.GenerateCryptoHash(old_password, salt3)
		secret2, _ = hashes.GenerateSecret(config.UserSecretLen)
		secret3, _ = hashes.GenerateSecret(config.UserSecretLen)

		user2 = &structs.User{
			ID:     userID2,
			Login:  login2,
			Secret: secret2,
			Salt:   salt2,
			Hash:   hash2,
		}

		user3 = &structs.User{
			ID:     userID3,
			Login:  login3,
			Secret: secret3,
			Salt:   salt3,
			Hash:   hash3,
		}
		userID4 = "not_exist"
	)

	mockStorage.EXPECT().GetUserByID(ctx, userID1).Return(user1, nil)
	mockStorage.EXPECT().SaveUser(ctx, gomock.Any()).Return(nil)

	mockStorage.EXPECT().GetUserByID(ctx, userID2).Return(user2, nil)

	mockStorage.EXPECT().GetUserByID(ctx, userID3).Return(user3, nil)
	mockStorage.EXPECT().GetUserByID(ctx, userID4).Return(nil, myerrors.ErrUserNotExist)

	tests := []struct {
		name    string
		usecase *UserUsecase
		args    args
		wantErr bool
	}{
		{
			name:    "Test 1: valid password change",
			usecase: u,
			args: args{
				ctx:       ctx,
				userID:    userID1,
				oldPass:   old_password,
				newPass:   valid_password,
				reNewPass: valid_re_password,
			},
			wantErr: false,
		},
		{
			name:    "Test 2: invalid password",
			usecase: u,
			args: args{
				ctx:       ctx,
				userID:    userID2,
				oldPass:   old_password,
				newPass:   invalid_password,
				reNewPass: invalid_re_password,
			},
			wantErr: true,
		},
		{
			name:    "Test 3:  password and rePassword are not the same change",
			usecase: u,
			args: args{
				ctx:       ctx,
				userID:    userID3,
				oldPass:   old_password,
				newPass:   not_same_password,
				reNewPass: not_same_re_password,
			},
			wantErr: true,
		},
		{
			name:    "Test 4:  userNotExist",
			usecase: u,
			args: args{
				ctx:       ctx,
				userID:    userID4,
				oldPass:   old_password,
				newPass:   valid_password,
				reNewPass: valid_re_password,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := u.ChangePassword(tt.args.ctx, tt.args.userID, tt.args.oldPass, tt.args.newPass, tt.args.reNewPass); (err != nil) != tt.wantErr {
				t.Errorf("ChangePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

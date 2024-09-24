package usersUC

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"

	"passkeeper/internal/entities/structs"
	"passkeeper/internal/storage"
	"passkeeper/internal/storage/memstorage"
)

func TestUserUsecase_RegisterUser(t *testing.T) {
	type fields struct {
		ctx             context.Context
		storage         storage.UserStorage
		tokenLifeTime   time.Duration
		userSecretLen   int
		userPassSaltLen int
		log             *zerolog.Logger
	}
	type args struct {
		ctx      context.Context
		login    string
		password string
		rePass   string
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		ctx       = context.Background()
		login1    = "login1"
		password1 = "password"

		login2 = "login2"
	)

	storage := memstorage.NewMemStorage(ctx, memstorage.WithUsers(
		map[string]*structs.User{
			login2: nil,
		}))

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
			name: "Test 1: register valid user",
			fields: fields{
				ctx:             ctx,
				storage:         storage,
				tokenLifeTime:   110,
				userSecretLen:   222,
				userPassSaltLen: 333,
				log:             &zerolog.Logger{},
			},
			args: args{
				ctx:      ctx,
				login:    login1,
				password: password1,
				rePass:   password1,
			},
			wantErr: false,
		},
		{
			name: "Test 2: register existed user",
			fields: fields{
				ctx:             ctx,
				storage:         storage,
				tokenLifeTime:   110,
				userSecretLen:   222,
				userPassSaltLen: 333,
				log:             &zerolog.Logger{},
			},
			args: args{
				ctx:      ctx,
				login:    login2,
				password: "test",
				rePass:   "test",
			},
			wantErr: true,
		},
		{
			name: "Test 3: password and rePass are not match",
			fields: fields{
				ctx:             ctx,
				storage:         storage,
				tokenLifeTime:   110,
				userSecretLen:   222,
				userPassSaltLen: 333,
				log:             &zerolog.Logger{},
			},
			args: args{
				ctx:      ctx,
				login:    login2,
				password: "test",
				rePass:   "test1",
			},
			wantErr: true,
		},
		{
			name: "Test 4: password is empty",
			fields: fields{
				ctx:             ctx,
				storage:         storage,
				tokenLifeTime:   110,
				userSecretLen:   222,
				userPassSaltLen: 333,
				log:             &zerolog.Logger{},
			},
			args: args{
				ctx:      ctx,
				login:    login2,
				password: "",
				rePass:   "",
			},
			wantErr: true,
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

			_, _, err := u.RegisterUser(tt.args.ctx, tt.args.login, tt.args.password, tt.args.rePass)
			if err != nil {
				if (err != nil) != tt.wantErr {
					t.Errorf("LoginUser() error = %v, wantErr %v", err, tt.wantErr)
					return
				} else {
					return
				}
			}

		})
	}
}

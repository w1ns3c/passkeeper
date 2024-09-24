package usersUC

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"passkeeper/internal/storage"
	mocks "passkeeper/mock"
)

func TestNewUserUsecase(t *testing.T) {
	type args struct {
		ctx       context.Context
		log       *zerolog.Logger
		saltLen   int
		secretLen int
		tokenLife time.Duration
		storage   storage.Storage
	}

	mockCtr := gomock.NewController(t)
	defer mockCtr.Finish()

	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "Test: success creation UserUsecase",
			args: args{
				ctx:       context.Background(),
				log:       &zerolog.Logger{},
				saltLen:   20,
				secretLen: 32,
				tokenLife: time.Hour,
				storage:   mocks.NewMockStorage(mockCtr),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUC := NewUserUsecase(tt.args.ctx)
			gotUC.SetSaltLen(tt.args.saltLen).
				SetSecretLen(tt.args.secretLen).
				SetContext(tt.args.ctx).
				SetLog(tt.args.log).
				SetTokenLifeTime(tt.args.tokenLife).
				SetStorage(tt.args.storage)

			require.Equal(t, tt.args.saltLen, gotUC.userPassSaltLen)
			require.Equal(t, tt.args.secretLen, gotUC.userSecretLen)
			require.Equal(t, tt.args.ctx, gotUC.ctx)
			require.Equal(t, tt.args.log, gotUC.log)
			require.Equal(t, tt.args.saltLen, gotUC.userPassSaltLen)
			require.Equal(t, tt.args.storage, gotUC.storage)
		})
	}
}

func TestUserUsecase_SetContext(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
	}{
		{
			name: "Test set salt len",
			ctx:  context.Background(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserUsecase{}
			gotUC := u.SetContext(tt.ctx)
			require.Equal(t, tt.ctx, gotUC.ctx)

		})
	}
}

func TestUserUsecase_SetLog(t *testing.T) {

	tests := []struct {
		name string
		log  *zerolog.Logger
	}{
		{
			name: "Test set logger",
			log:  &zerolog.Logger{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserUsecase{}
			gotUC := u.SetLog(tt.log)
			require.Equal(t, tt.log, gotUC.log)

		})
	}
}

func TestUserUsecase_SetSaltLen(t *testing.T) {
	tests := []struct {
		name    string
		saltlen int
	}{
		{
			name:    "Test set salt len",
			saltlen: 120,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserUsecase{}
			gotUC := u.SetSaltLen(tt.saltlen)
			require.Equal(t, tt.saltlen, gotUC.userPassSaltLen)

		})
	}
}

func TestUserUsecase_SetSecretLen(t *testing.T) {
	tests := []struct {
		name      string
		secretlen int
	}{
		{
			name:      "Test set secretlen",
			secretlen: 120,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserUsecase{}
			gotUC := u.SetSecretLen(tt.secretlen)
			require.Equal(t, tt.secretlen, gotUC.userSecretLen)

		})
	}
}

func TestUserUsecase_SetStorage(t *testing.T) {

	mockCtr := gomock.NewController(t)
	defer mockCtr.Finish()

	tests := []struct {
		name    string
		storage storage.Storage
	}{
		{
			name:    "Test set storage",
			storage: mocks.NewMockStorage(mockCtr),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserUsecase{}
			gotUC := u.SetStorage(tt.storage)
			require.Equal(t, tt.storage, gotUC.storage)

		})
	}
}

func TestUserUsecase_SetTokenLifeTime(t *testing.T) {
	tests := []struct {
		name          string
		tokenLifeTime time.Duration
	}{
		// TODO: Add test cases.
		{
			name:          "Test set token life time",
			tokenLifeTime: time.Second * 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserUsecase{}
			gotUC := u.SetTokenLifeTime(tt.tokenLifeTime)
			require.Equal(t, tt.tokenLifeTime, gotUC.tokenLifeTime)

		})
	}
}

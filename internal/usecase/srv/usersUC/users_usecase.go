package usersUC

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	"passkeeper/internal/entities/config"
	"passkeeper/internal/storage"
	"passkeeper/internal/storage/memstorage"
)

// UserUsecaseInf describe main user functional on server side
//
//go:generate mockgen -source users_usecase.go -destination ../../../../mocks/usecase/users_usecase/users_usecase.go -package mocks
type UserUsecaseInf interface {
	RegisterUser(ctx context.Context, login string, password string, rePass string) (token, secret string, err error)

	ChangePassword(ctx context.Context, userID, oldPass, newPass, reNewPass string) (err error)
	GetTokenSalt(ctx context.Context, userID string) string

	LoginUser(ctx context.Context, login string, password string) (token, secret string, err error)
}

// UserUsecase implement UserUsecaseInf
type UserUsecase struct {
	ctx             context.Context
	storage         storage.UserStorage
	tokenLifeTime   time.Duration
	userSecretLen   int
	userPassSaltLen int
	log             *zerolog.Logger
}

// NewUserUsecase is constructo for UserUsecase
func NewUserUsecase(ctx context.Context) *UserUsecase {
	return &UserUsecase{
		storage:         memstorage.NewMemStorage(ctx),
		tokenLifeTime:   config.TokenLifeTime,
		userPassSaltLen: config.UserPassSaltLen,
		userSecretLen:   config.UserSecretLen,
	}
}

// SetStorage set storage to UserUsecase
func (u *UserUsecase) SetStorage(storage storage.UserStorage) *UserUsecase {
	u.storage = storage

	return u
}

// SetTokenLifeTime set secretlen to UserUsecase
func (u *UserUsecase) SetTokenLifeTime(tokenLifeTime time.Duration) *UserUsecase {
	u.tokenLifeTime = tokenLifeTime

	return u
}

// SetLog add logger to UserUsecase
func (u *UserUsecase) SetLog(logger *zerolog.Logger) *UserUsecase {
	u.log = logger

	return u
}

// SetContext add context to UserUsecase
func (u *UserUsecase) SetContext(ctx context.Context) *UserUsecase {
	u.ctx = ctx

	return u
}

// SetSecretLen set len of user's secret string
func (u *UserUsecase) SetSecretLen(length int) *UserUsecase {
	u.userSecretLen = length

	return u
}

// SetSaltLen set len of user's salt string
func (u *UserUsecase) SetSaltLen(length int) *UserUsecase {
	u.userPassSaltLen = length

	return u
}

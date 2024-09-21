package usersUC

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"

	"passkeeper/internal/config"
	"passkeeper/internal/storage"
	"passkeeper/internal/storage/memstorage"
)

var (
	ErrGetUser = fmt.Errorf("can't get user by ID")
	ErrGenHash = fmt.Errorf("can't generate hash of password")

	ErrWrongOldPassword = fmt.Errorf("old password is wrong")
	ErrRepassNotSame    = fmt.Errorf("new pass and repeat not the same")

	ErrWrongPassword = fmt.Errorf("wrong password for username")
	ErrWrongUsername = fmt.Errorf("username not found")
	ErrWrongAuth     = fmt.Errorf("wrong user/password")

	ErrUserSecret = fmt.Errorf("can't generate user secret hash")

	//ErrPassNotTheSame = fmt.Errorf("passwords not match")
	ErrPassIsEmpty   = fmt.Errorf("password is empty")
	ErrRePassIsEmpty = fmt.Errorf("password repeat is empty")
)

type UserUsecaseInf interface {
	RegisterUser(ctx context.Context, login string, password string, rePass string) (token, secret string, err error)

	ChangePassword(ctx context.Context, userID, oldPass, newPass, reNewPass string) (err error)
	GetTokenSalt(ctx context.Context, userID string) string

	LoginUser(ctx context.Context, login string, password string) (token, secret string, err error)
}

type UserUsecase struct {
	ctx             context.Context
	storage         storage.UserStorage
	tokenLifeTime   time.Duration
	userSecretLen   int
	userPassSaltLen int
	log             *zerolog.Logger
}

func (u *UserUsecase) GetTokenSalt(ctx context.Context, userID string) string {
	user, err := u.storage.GetUserByID(ctx, userID)
	if err != nil {
		u.log.Error().
			Err(err).Msg("can't get user's salt from storage")

		return ""
	}

	return user.Salt
}

func NewUserUsecase() *UserUsecase {
	return &UserUsecase{
		storage:         memstorage.NewMemStorage(),
		tokenLifeTime:   config.TokenLifeTime,
		userPassSaltLen: config.UserPassSaltLen,
		userSecretLen:   config.UserSecretLen,
	}
}

func (u *UserUsecase) SetStorage(storage storage.UserStorage) *UserUsecase {
	u.storage = storage

	return u
}

func (u *UserUsecase) SetTokenLifeTime(tokenLifeTime time.Duration) *UserUsecase {
	u.tokenLifeTime = tokenLifeTime

	return u
}

func (u *UserUsecase) SetLog(logger *zerolog.Logger) *UserUsecase {
	u.log = logger

	return u
}

func (u *UserUsecase) SetContext(ctx context.Context) *UserUsecase {
	u.ctx = ctx

	return u
}

func (u *UserUsecase) SetSecretLen(length int) *UserUsecase {
	u.userSecretLen = length

	return u
}

func (u *UserUsecase) SetSaltLen(length int) *UserUsecase {
	u.userPassSaltLen = length

	return u
}

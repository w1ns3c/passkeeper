package usersUC

import (
	"context"
	"fmt"
	"time"

	"passkeeper/internal/storage"

	"github.com/rs/zerolog"
)

var (
	ErrGetUser       = fmt.Errorf("can't get user by ID")
	ErrGenHash       = fmt.Errorf("can't generate hash of password")
	ErrWrongPassword = fmt.Errorf("old password is wrong")
	ErrRepassNotSame = fmt.Errorf("new pass and repeat not the same")

	ErrWrongAuth    = fmt.Errorf("wrong user/password")
	ErrInvalidToken = fmt.Errorf("token sign is not valid")

	ErrUserSecret = fmt.Errorf("can't generate user secret hash")
)

type UserUsecaseInf interface {
	RegisterUser(ctx context.Context, login string, password string, rePass string) (token, secret string, err error)

	ChangePassword(ctx context.Context, userID, oldPass, newPass, reNewPass string) (err error)
	GetTokenSalt() string

	LoginUser(ctx context.Context, login string, password string) (token, secret string, err error)
}

type UserUsecase struct {
	storage       storage.UserStorage
	salt          string
	tokenLifeTime time.Duration
	userSecretLen int
	log           *zerolog.Logger
}

func NewUserUsecase(storage storage.UserStorage, salt string, tokenLifeTime time.Duration, userSecretLen int, log *zerolog.Logger) *UserUsecase {
	return &UserUsecase{storage: storage, salt: salt, tokenLifeTime: tokenLifeTime, userSecretLen: userSecretLen, log: log}
}

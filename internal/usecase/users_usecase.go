package usecase

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/w1nsec/go-examples/crypto"
	"golang.org/x/crypto/bcrypt"

	"github.com/w1nsec/passkeeper/internal/entities"
	"github.com/w1nsec/passkeeper/internal/storage"
)

type UserUsecaseInf interface {
	LoginUser(ctx context.Context, login string, password string) (token string, err error)
	RegisterUser(ctx context.Context, login string, password string, rePass string) (token string, err error)
}

type UserUsecase struct {
	storage       storage.UserStorage
	salt          string
	tokenLifeTime time.Duration
	userSecretLen int
}

func (u *UserUsecase) LoginUser(ctx context.Context, login string, password string) (token string, err error) {
	hash, err := GenerateHash(password, u.salt)
	if err != nil {
		return "", fmt.Errorf("can't generate hash of password: %v", err)
	}

	user, err := u.storage.LoginUser(ctx, login, hash)
	if err != nil {
		return "", fmt.Errorf("wrong login/password: %v", err)
	}

	return GenerateToken(user.ID, u.salt, u.tokenLifeTime)
}

func (u *UserUsecase) RegisterUser(ctx context.Context, login string, password string, rePass string) (token string, err error) {
	if password != rePass {
		return "", entities.ErrPassNotTheSame
	}

	// checking login free
	exist, err := u.storage.CheckUserExist(ctx, login)
	if !errors.Is(err, entities.ErrUserNotFound) || exist {
		return "", fmt.Errorf("user is already exist:%v", err)
	}

	id := GenerateID(login, u.salt)

	hash, err := GenerateHash(password, u.salt)
	if err != nil {
		return "", fmt.Errorf("can't generate hash of password: %v", err)
	}

	secret, err := GenerateSecret(u.userSecretLen)
	if err != nil {
		return "", fmt.Errorf("can't generate secret for user: %v", err)
	}

	user := &entities.User{
		ID:    id,
		Login: login,
		//Credential: password,
		Hash:   hash,
		Secret: secret,
	}

	err = u.storage.SaveUser(ctx, user)
	if err != nil {
		return "", err
	}

	token, err = GenerateToken(user.ID, u.salt, u.tokenLifeTime)
	if err != nil {
		return "", fmt.Errorf("can't generate user token: %v", err)
	}

	return token, nil
}

func GenerateHash(password, salt string) (hash string, err error) {
	password = fmt.Sprintf("%s-%s.%s.%s", string(salt), string(password), string(password), string(salt))
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(h), nil
}

func GenerateID(login, salt string) string {

	hash := md5.Sum([]byte(fmt.Sprintf("%s.%s.%s", salt, login, salt)))
	return hex.EncodeToString(hash[:])
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func GenerateToken(userid string, secret string, lifetime time.Duration) (token string, err error) {
	tokenLife := time.Now().Add(lifetime)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(tokenLife),
		},
		UserID: userid,
	})
	token, err = jwtToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func GenerateSecret(secretLen int) (secret string, err error) {
	sl, err := crypto.GenRandSlice(secretLen)
	if err != nil {
		return "", nil
	}

	return hex.EncodeToString(sl), nil
}

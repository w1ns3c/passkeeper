package srv

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/w1ns3c/go-examples/crypto"
	"github.com/w1ns3c/passkeeper/internal/entities"
	"github.com/w1ns3c/passkeeper/internal/storage"
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

func (u *UserUsecase) ChangePassword(ctx context.Context, userID, oldPass, newPass, reNewPass string) (err error) {
	user, err := u.storage.GetUserByID(ctx, userID)
	if err != nil {
		u.log.Error().Err(err).
			Msg(ErrGetUser.Error())

		return ErrGetUser
	}

	equal := ComparePassAndCryptoHash(oldPass, user.Hash, u.salt)
	if !equal {
		u.log.Error().
			Err(ErrWrongPassword).Send()

		return ErrWrongPassword
	}

	if newPass != reNewPass {
		u.log.Error().
			Err(ErrRepassNotSame).Send()

		return ErrRepassNotSame
	}

	hNew1, err := GenerateCryptoHash(newPass, u.salt)
	if err != nil {
		u.log.Error().Err(err).
			Msg(ErrGenHash.Error())

		return ErrGenHash
	}

	user.Hash = hNew1
	return u.storage.SaveUser(ctx, user)

}

func (u *UserUsecase) RegisterUser(ctx context.Context, login string,
	password string, rePass string) (token, secretForCreds string, err error) {

	if password != rePass {
		return "", "", entities.ErrPassNotTheSame
	}

	// checking login free
	exist, err := u.storage.CheckUserExist(ctx, login)
	if !errors.Is(err, entities.ErrUserNotFound) || exist {
		return "", "", fmt.Errorf("user is already exist:%v", err)
	}

	hash, err := GenerateCryptoHash(password, u.salt)
	if err != nil {
		return "", "", fmt.Errorf("can't generate hash of password: %v", err)
	}

	m := md5.Sum([]byte(hash))
	id := GenerateID(hex.EncodeToString(m[:]), u.salt)

	secret, err := GenerateSecret(u.userSecretLen)
	if err != nil {
		return "", "", fmt.Errorf("can't generate secret for user: %v", err)
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
		return "", "", err
	}

	hashedSecret, err := HashSecret(user.Secret)
	if err != nil {
		u.log.Error().Err(err).
			Msg(ErrUserSecret.Error())

		return "", "", ErrWrongPassword
	}

	token, err = GenerateToken(user.ID, user.Secret, u.tokenLifeTime)
	if err != nil {
		return "", "", fmt.Errorf("can't generate user token: %v", err)
	}

	return token, hashedSecret, nil
}

// HashSecret save secret before sent to client
// User secret
// Send secret: 		md5(aes256(user.secret, key:user.secret))
// Secret for token: 	user.secret
func HashSecret(secret string) (hash string, err error) {
	key := sha256.Sum256([]byte(secret))
	secretAES, err := crypto.EncryptAES([]byte(secret), key[:])
	if err != nil {
		return "", err
	}

	hashedSecret := fmt.Sprintf("%x", md5.Sum(secretAES))

	return hashedSecret, nil
}

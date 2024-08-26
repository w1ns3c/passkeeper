package usersUC

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"passkeeper/internal/config"
	"passkeeper/internal/utils/hashes"

	"passkeeper/internal/entities"

	"github.com/w1ns3c/go-examples/crypto"
)

// RegisterUser function for register user in app
func (u *UserUsecase) RegisterUser(ctx context.Context, login string,
	password string, rePass string) (token, secretForCreds string, err error) {

	if password != rePass {
		return "", "", entities.ErrPassNotTheSame
	}

	// checking login free
	exist, err := u.storage.CheckUserExist(ctx, login)
	if (!errors.Is(err, entities.ErrUserNotFound) && err != nil) || exist {
		return "", "", fmt.Errorf("user is already exist:%v", err)
	}

	userSalt, err := crypto.GenRandStr(config.UserPassSaltLen)
	if err != nil {
		return "", "", fmt.Errorf("user is already exist:%v", err)
	}
	hash, err := hashes.GenerateCryptoHash(password, userSalt)
	if err != nil {
		return "", "", fmt.Errorf("can't generate hash of password: %v", err)
	}

	m := md5.Sum([]byte(hash))
	id := hashes.GenerateUserID(hex.EncodeToString(m[:]), userSalt)

	secret, err := hashes.GenerateSecret(u.userSecretLen)
	if err != nil {
		return "", "", fmt.Errorf("can't generate secret for user: %v", err)
	}

	user := &entities.User{
		ID:     id,
		Login:  login,
		Hash:   hash,
		Salt:   userSalt,
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

	token, err = hashes.GenerateToken(user.ID, user.Secret, u.tokenLifeTime)
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

package usersUC

import (
	"context"
	"errors"
	"fmt"

	"github.com/w1ns3c/go-examples/crypto"

	"passkeeper/internal/entities/config"
	"passkeeper/internal/entities/hashes"
	"passkeeper/internal/entities/myerrors"
	"passkeeper/internal/entities/structs"
)

// RegisterUser function for register user in app on server side
func (u *UserUsecase) RegisterUser(ctx context.Context, login string,
	password string, rePass string) (token, secretForCreds string, err error) {

	// here password should be hashes.Hash(password)
	if password != rePass {
		u.log.Error().Err(myerrors.ErrRepassNotSame).
			Str("login", login).
			Msg("can't register user: ")

		return "", "", myerrors.ErrRepassNotSame
	}

	if password == "" || rePass == "" {
		u.log.Error().Err(myerrors.ErrPassIsEmpty).
			Str("login", login).
			Msg("can't register user:")

		return "", "", myerrors.ErrPassIsEmpty
	}

	// checking login free
	exist, err := u.storage.CheckUserExist(ctx, login)
	if (!errors.Is(err, myerrors.ErrUserNotExist) && err != nil) || exist {
		u.log.Error().Err(err).
			Str("login", login).
			Msg("can't register user, possible it exist")

		return "", "", fmt.Errorf("user is already exist:%v", myerrors.ErrAlreadyExist)
	}

	userSalt, err := crypto.GenRandStr(config.UserPassSaltLen)
	if err != nil {
		return "", "", fmt.Errorf("user is already exist:%v", err)
	}

	hash, err := hashes.GenerateCryptoHash(password, userSalt)
	if err != nil {
		return "", "", fmt.Errorf("can't generate hash of password: %v", err)
	}

	id := hashes.GenerateUserID(password, userSalt)

	secret, err := hashes.GenerateSecret(u.userSecretLen)
	if err != nil {
		return "", "", fmt.Errorf("can't generate secret for user: %v", err)
	}

	secureSecret, err := hashes.EncryptSecret(secret, password)

	user := &structs.User{
		ID:     id,
		Login:  login,
		Hash:   hash,
		Salt:   userSalt,
		Secret: secureSecret,
	}

	err = u.storage.SaveUser(ctx, user)
	if err != nil {
		return "", "", err
	}

	token, err = hashes.GenerateToken(user.ID, userSalt, u.tokenLifeTime)
	if err != nil {
		return "", "", fmt.Errorf("can't generate user token: %v", err)
	}

	// user.CredsSecret is hashes.EncryptSecret(secret, hashes.Hash(clearPassword))
	return token, user.Secret, nil
}

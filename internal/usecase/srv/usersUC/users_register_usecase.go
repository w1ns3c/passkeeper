package usersUC

import (
	"context"
	"errors"
	"fmt"

	"github.com/w1ns3c/go-examples/crypto"

	"passkeeper/internal/entities"
	"passkeeper/internal/entities/config"
	"passkeeper/internal/entities/hashes"
	"passkeeper/internal/entities/myerrors"
)

// RegisterUser function for register user in app on server side
func (u *UserUsecase) RegisterUser(ctx context.Context, login string,
	password string, rePass string) (token, secretForCreds string, err error) {

	// here password should be hashes.Hash(password)
	if password != rePass {
		return "", "", myerrors.ErrRepassNotSame
	}

	if password == "" {
		return "", "", myerrors.ErrPassIsEmpty
	}

	if rePass == "" {
		return "", "", myerrors.ErrRePassIsEmpty
	}

	// checking login free
	exist, err := u.storage.CheckUserExist(ctx, login)
	if (!errors.Is(err, myerrors.ErrUserNotExist) && err != nil) || exist {
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

	id := hashes.GenerateUserID(password, userSalt)

	secret, err := hashes.GenerateSecret(u.userSecretLen)
	if err != nil {
		return "", "", fmt.Errorf("can't generate secret for user: %v", err)
	}

	secureSecret, err := hashes.EncryptSecret(secret, password)

	user := &entities.User{
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

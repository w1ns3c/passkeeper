package usersUC

import (
	"context"
	"errors"
	"fmt"

	"github.com/w1ns3c/go-examples/crypto"

	"passkeeper/internal/entities"
	"passkeeper/internal/entities/config"
	"passkeeper/internal/entities/hashes"
)

// RegisterUser function for register user in app
func (u *UserUsecase) RegisterUser(ctx context.Context, login string,
	password string, rePass string) (token, secretForCreds string, err error) {

	// here password should be hashes.Hash(password)
	if password != rePass {
		return "", "", ErrRepassNotSame
	}

	if password == "" {
		return "", "", ErrPassIsEmpty
	}

	if rePass == "" {
		return "", "", ErrRePassIsEmpty
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

	//hashedSecret, err := HashSecret(user.FullSecret)
	//if err != nil {
	//	u.log.Error().Err(err).
	//		Msg(ErrUserSecret.Error())
	//
	//	return "", "", ErrWrongOldPassword
	//}

	token, err = hashes.GenerateToken(user.ID, userSalt, u.tokenLifeTime)
	if err != nil {
		return "", "", fmt.Errorf("can't generate user token: %v", err)
	}

	// user.CredsSecret is hashes.EncryptSecret(secret, hashes.Hash(clearPassword))
	return token, user.Secret, nil
}

//
//// HashSecret save secret before sent to client
//// User secret
//// Send secret: 		md5(aes256(user.secret, key:user.secret))
//// FullSecret for token: 	user.secret
//func HashSecret(secret string) (hash string, err error) {
//	key := sha256.Sum256([]byte(secret))
//	secretAES, err := crypto.EncryptAES([]byte(secret), key[:])
//	if err != nil {
//		return "", err
//	}
//
//	hashedSecret := fmt.Sprintf("%x", md5.Sum(secretAES))
//
//	return hashedSecret, nil
//}

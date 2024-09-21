package usersUC

import (
	"context"

	"passkeeper/internal/entities/hashes"
)

func (u *UserUsecase) LoginUser(ctx context.Context, login string, password string) (token string, secret string, err error) {

	// here password should be hashes.Hash(password)

	user, err := u.storage.GetUserByLogin(ctx, login)
	if err != nil {
		return "", "", ErrWrongUsername
	}

	same := hashes.ComparePassAndCryptoHash(password, user.Hash, user.Salt)
	if !same {
		//u.log.Error().Err(err).
		//	Msg(ErrWrongPassword.Error())

		return "", "", ErrWrongPassword
	}

	//hashedSecret, err := HashSecret(user.FullSecret)
	//if err != nil {
	//	u.log.Error().Err(err).
	//		Msg(ErrUserSecret.Error())
	//
	//	return "", "", ErrWrongOldPassword
	//}

	//clearSecret, err := hashes.DecryptSecret(user.CredsSecret, password)
	//if err != nil {
	//	return "", "", err
	//}

	//token, err = hashes.GenerateToken(user.ID, clearSecret, u.tokenLifeTime)

	token, err = hashes.GenerateToken(user.ID, user.Salt, u.tokenLifeTime)

	// user.CredsSecret is hashes.EncryptSecret(secret, hashes.Hash(clearPassword))
	return token, user.Secret, err
}

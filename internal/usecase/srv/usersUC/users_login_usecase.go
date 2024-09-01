package usersUC

import (
	"context"
	"passkeeper/internal/entities/hashes"
)

func (u *UserUsecase) LoginUser(ctx context.Context, login string, password string) (token string, secret string, err error) {

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

	clearSecret, err := hashes.DecryptSecret(user.Secret, password)
	if err != nil {
		return "", "", err
	}

	token, err = hashes.GenerateToken(user.ID, clearSecret, u.tokenLifeTime)

	return token, user.Secret, err
}

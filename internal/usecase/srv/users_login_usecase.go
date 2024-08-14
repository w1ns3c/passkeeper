package srv

import "context"

func (u *UserUsecase) LoginUser(ctx context.Context, login string, password string) (token string, secret string, err error) {

	user, err := u.storage.GetUserByLogin(ctx, login)
	if err != nil {
		return "", "", ErrWrongAuth
	}

	same := ComparePassAndCryptoHash(password, user.Hash, u.salt)
	if !same {
		u.log.Error().Err(err).
			Msg(ErrWrongPassword.Error())

		return "", "", ErrWrongPassword
	}

	hashedSecret, err := HashSecret(user.Secret)
	if err != nil {
		u.log.Error().Err(err).
			Msg(ErrUserSecret.Error())

		return "", "", ErrWrongPassword
	}

	token, err = GenerateToken(user.ID, user.Secret, u.tokenLifeTime)
	return token, hashedSecret, err
}

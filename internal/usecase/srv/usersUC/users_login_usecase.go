package usersUC

import (
	"context"

	"passkeeper/internal/entities/hashes"
	"passkeeper/internal/entities/myerrors"
)

// LoginUser function for register user in app on server side
func (u *UserUsecase) LoginUser(ctx context.Context, login string, password string) (token string, secret string, err error) {

	// here password should be hashes.Hash(password)

	user, err := u.storage.GetUserByLogin(ctx, login)
	if err != nil {
		u.log.Error().Err(err).
			Msg(myerrors.ErrUserNotExist.Error())

		return "", "", myerrors.ErrWrongUsername
	}

	same := hashes.ComparePassAndCryptoHash(password, user.Hash, user.Salt)
	if !same {
		u.log.Error().Err(err).
			Msg(myerrors.ErrWrongPassword.Error())

		return "", "", myerrors.ErrWrongPassword
	}

	token, err = hashes.GenerateToken(user.ID, user.Salt, u.tokenLifeTime)

	// user.CredsSecret is hashes.EncryptSecret(secret, hashes.Hash(clearPassword))
	return token, user.Secret, err
}

// GetTokenSalt return user's salt from storage
// important, because salt need to encrypt/decrypt user token
func (u *UserUsecase) GetTokenSalt(ctx context.Context, userID string) string {
	user, err := u.storage.GetUserByID(ctx, userID)
	if err != nil {
		u.log.Error().
			Err(err).Msg("can't get user's salt from storage")

		return ""
	}

	return user.Salt
}

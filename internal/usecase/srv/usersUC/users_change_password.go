package usersUC

import "context"

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

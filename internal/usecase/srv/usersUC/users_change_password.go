package usersUC

import (
	"context"

	"passkeeper/internal/entities/hashes"
)

func (u *UserUsecase) ChangePassword(ctx context.Context, userID, oldPass, newPass, reNewPass string) (err error) {
	user, err := u.storage.GetUserByID(ctx, userID)
	if err != nil {
		u.log.Error().Err(err).
			Msg(ErrGetUser.Error())

		return ErrGetUser
	}

	equal := hashes.ComparePassAndCryptoHash(oldPass, user.Hash, user.Salt)
	if !equal {
		u.log.Error().
			Err(ErrWrongOldPassword).Send()

		return ErrWrongOldPassword
	}

	if newPass != reNewPass {
		u.log.Error().
			Err(ErrRepassNotSame).Send()

		return ErrRepassNotSame
	}

	hNew1, err := hashes.GenerateCryptoHash(newPass, user.Salt)
	if err != nil {
		u.log.Error().Err(err).
			Msg(ErrGenHash.Error())

		return ErrGenHash
	}

	user.Hash = hNew1
	return u.storage.SaveUser(ctx, user)

}

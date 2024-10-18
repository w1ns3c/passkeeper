package usersUC

import (
	"context"

	"passkeeper/internal/entities/hashes"
	"passkeeper/internal/entities/myerrors"
)

// ChangePassword func to change user password on server side
func (u *UserUsecase) ChangePassword(ctx context.Context, userID, oldPass, newPass, reNewPass string) (err error) {
	user, err := u.storage.GetUserByID(ctx, userID)
	if err != nil {
		u.log.Error().Err(err).
			Msg(myerrors.ErrGetUser.Error())

		return myerrors.ErrGetUser
	}

	if newPass != reNewPass {
		u.log.Error().
			Err(myerrors.ErrRepassNotSame).Send()

		return myerrors.ErrRepassNotSame
	}

	equal := hashes.ComparePassAndCryptoHash(oldPass, user.Hash, user.Salt)
	if !equal {
		u.log.Error().
			Err(myerrors.ErrWrongOldPassword).Send()

		return myerrors.ErrWrongOldPassword
	}

	hNew1, err := hashes.GenerateCryptoHash(newPass, user.Salt)
	if err != nil {
		u.log.Error().Err(err).
			Msg(myerrors.ErrGenHash.Error())

		return myerrors.ErrGenHash
	}

	user.Hash = hNew1
	return u.storage.SaveUser(ctx, user)

}

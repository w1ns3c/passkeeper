package credentialsUC

import (
	"context"
	"passkeeper/internal/utils/hashes"
	"time"

	"passkeeper/internal/entities"
)

func (u *CredUsecase) GetCredential(ctx context.Context, userToken, credID string) (cred *entities.Credential, err error) {
	cred, err = u.storage.GetCredential(ctx, userToken, credID)
	if err != nil {
		return nil, err
	}

	cred.Password, err = DecryptPass(cred.Password)
	if err != nil {
		return nil, err
	}

	return cred, nil
}

func (u *CredUsecase) AddCredential(ctx context.Context,
	userID string, cred *entities.Credential) error {

	//cred.Password, err = EncryptPass(cred.Password)
	//cred.ID = hashes.GeneratePassID(sec, salt)

	cred.ID = hashes.GeneratePassID2()
	VerifyCredDate(cred)

	return u.storage.AddCredential(ctx, userID, cred)
}

func (u *CredUsecase) UpdateCredential(ctx context.Context,
	userID string, cred *entities.Credential) error {

	//sec, err := usersUC.GenerateSecret(config.UserSecretLen)
	//if err != nil {
	//	return err
	//}

	//cred.ID = usersUC.GenerateID(sec, u.salt)
	//cred.Password, err = EncryptPass(cred.Password)
	VerifyCredDate(cred)

	return u.storage.UpdateCredential(ctx, userID, cred)
}

func (u *CredUsecase) DeleteCredential(ctx context.Context,
	userToken, credID string) error {

	return u.storage.DeleteCredential(ctx, userToken, credID)
}

func (u *CredUsecase) ListCredentials(ctx context.Context,
	userToken string) (creds []*entities.Credential, err error) {

	creds, err = u.storage.GetAllCredentials(ctx, userToken)
	if err != nil {
		return nil, err
	}

	for ind := 0; ind < len(creds); ind++ {
		creds[ind].Password, err = DecryptPass(creds[ind].Password)
		if err != nil {
			u.log.Error().Err(err).
				Msgf("%s with ID: %s (user: %s)", ErrNoDecrypt, creds[ind].ID, userToken)
			creds[ind].Password = ErrNoDecrypt
		}
	}

	return creds, nil
}

// VerifyCredDate verify date/time in received credential
// if date too old, update it
func VerifyCredDate(cred *entities.Credential) {
	now := time.Now()
	if cred.Date.Sub(now) > time.Hour*24 {
		cred.Date = now
	}
}

func (u *CredUsecase) GetUserSecret(ctx context.Context, userToken string) string {
	//TODO implement me
	panic("implement me")
}

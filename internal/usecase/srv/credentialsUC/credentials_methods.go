package credentialsUC

import (
	"context"
	"passkeeper/internal/entities/hashes"
	"time"

	"passkeeper/internal/entities"
)

func (u *CredUsecase) GetCredential(ctx context.Context, userToken, credID string) (cred *entities.CredBlob, err error) {
	return u.storage.GetCredential(ctx, userToken, credID)
}

func (u *CredUsecase) AddCredential(ctx context.Context,
	userID string, cred *entities.CredBlob) error {

	//cred.Password, err = EncryptPass(cred.Password)
	//cred.ID = hashes.GeneratePassID(sec, salt)

	cred.ID = hashes.GeneratePassID2()
	//VerifyCredDate(cred)

	return u.storage.AddCredential(ctx, userID, cred)
}

func (u *CredUsecase) UpdateCredential(ctx context.Context,
	userID string, cred *entities.CredBlob) error {

	//sec, err := usersUC.GenerateSecret(config.UserSecretLen)
	//if err != nil {
	//	return err
	//}

	//cred.ID = usersUC.GenerateID(sec, u.salt)
	//cred.Password, err = EncryptPass(cred.Password)
	//VerifyCredDate(cred)

	return u.storage.UpdateCredential(ctx, userID, cred)
}

func (u *CredUsecase) DeleteCredential(ctx context.Context,
	userToken, credID string) error {

	return u.storage.DeleteCredential(ctx, userToken, credID)
}

func (u *CredUsecase) ListCredentials(ctx context.Context,
	userID string) (creds []*entities.CredBlob, err error) {

	return u.storage.GetAllCredentials(ctx, userID)
}

// VerifyCredDate verify date/time in received credential
// if date too old, update it
func VerifyCredDate(cred *entities.Credential) {
	now := time.Now()
	if cred.Date.Sub(now) > time.Hour*24 {
		cred.Date = now
	}
}

func (u *CredUsecase) GetUserSalt(ctx context.Context, userID string) string {
	//TODO implement me
	panic("implement me")
}

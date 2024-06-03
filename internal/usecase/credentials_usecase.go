package usecase

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/w1ns3c/passkeeper/internal/config"
	"github.com/w1ns3c/passkeeper/internal/entities"
	"github.com/w1ns3c/passkeeper/internal/storage"
)

var (
	ErrNoDecrypt = "can't decrypt password"
)

type CredUsecaseInf interface {
	GetCredential(ctx context.Context, userToken, credID string) (cred *entities.Credential, err error)
	AddCredential(ctx context.Context, userToken string, cred *entities.Credential) error
	UpdateCredential(ctx context.Context, userToken string, cred *entities.Credential) error
	DeleteCredential(ctx context.Context, userToken, credID string) error
	ListCredentials(ctx context.Context, userToken string) (creds []*entities.Credential, err error)

	//VerifyCredDate(cred *entities.Credential) // check date/time in received credential

	//EncryptPwd(ctx context.Context, password string) (encPwd string, err error)
	//DecryptPass(ctx context.Context, encPwd string) (password string, err error)
}

type CredUsecase struct {
	storage storage.CredentialStorage
	salt    string
	log     *zerolog.Logger
}

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
	userToken string, cred *entities.Credential) error {

	sec, err := GenerateSecret(config.UserSecretLen)
	if err != nil {
		return err
	}

	cred.ID = GenerateID(sec, u.salt)
	cred.Password, err = EncryptPass(cred.Password)
	VerifyCredDate(cred)

	return u.storage.AddCredential(ctx, userToken, cred)
}

func (u *CredUsecase) UpdateCredential(ctx context.Context,
	userToken string, cred *entities.Credential) error {

	sec, err := GenerateSecret(config.UserSecretLen)
	if err != nil {
		return err
	}

	cred.ID = GenerateID(sec, u.salt)
	cred.Password, err = EncryptPass(cred.Password)
	VerifyCredDate(cred)

	return u.storage.UpdateCredential(ctx, userToken, cred)
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

func EncryptPass(password string) (encPass string, err error) {
	//TODO implement me
	return password, nil
}

func DecryptPass(encPass string) (password string, err error) {
	//TODO implement me
	return encPass, nil
}

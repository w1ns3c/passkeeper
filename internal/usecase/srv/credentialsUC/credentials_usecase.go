package credentialsUC

import (
	"context"

	"passkeeper/internal/storage/memstorage"

	"github.com/rs/zerolog"

	"passkeeper/internal/entities"
	"passkeeper/internal/storage"
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

type CredOption func(usecase *CredUsecase)

func NewCredUsecase() *CredUsecase {
	return &CredUsecase{
		salt:    "rand",
		storage: memstorage.NewMemStorage(),
	}
}

func WithSalt(salt string) CredOption {
	return func(usecase *CredUsecase) {
		usecase.salt = salt
	}
}

func WithStorage(db storage.Storage) CredOption {
	return func(usecase *CredUsecase) {
		usecase.storage = db
	}
}

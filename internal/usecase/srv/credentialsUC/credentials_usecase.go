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
	GetUserSalt(ctx context.Context, userID string) string // return secret user string
	GetCredential(ctx context.Context, userID, credID string) (cred *entities.CredBlob, err error)
	AddCredential(ctx context.Context, userID string, cred *entities.CredBlob) error
	UpdateCredential(ctx context.Context, userID string, cred *entities.CredBlob) error
	DeleteCredential(ctx context.Context, userID, credID string) error
	ListCredentials(ctx context.Context, userID string) (creds []*entities.CredBlob, err error)

	//VerifyCredDate(cred *entities.Credential) // check date/time in received credential

	//EncryptPwd(ctx context.Context, password string) (encPwd string, err error)
	//DecryptPass(ctx context.Context, encPwd string) (password string, err error)
}

type CredUsecase struct {
	ctx     context.Context
	storage storage.CredentialStorage
	//salt    string
	log *zerolog.Logger
}

type CredOption func(usecase *CredUsecase)

func newCredUsecase() *CredUsecase {
	return &CredUsecase{
		storage: memstorage.NewMemStorage(),
	}
}

func NewCredUCWithOpts(opts ...CredOption) *CredUsecase {
	credsUC := newCredUsecase()
	for _, opt := range opts {
		opt(credsUC)
	}

	return credsUC
}

func WithStorage(db storage.Storage) CredOption {
	return func(usecase *CredUsecase) {
		usecase.storage = db
	}
}

func WithContext(ctx context.Context) CredOption {
	return func(usecase *CredUsecase) {
		usecase.ctx = ctx
	}
}

func WithLogger(logger *zerolog.Logger) CredOption {
	return func(usecase *CredUsecase) {
		usecase.log = logger
	}
}

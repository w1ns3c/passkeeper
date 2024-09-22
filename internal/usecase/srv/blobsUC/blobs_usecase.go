package blobsUC

import (
	"context"

	"github.com/rs/zerolog"

	"passkeeper/internal/entities"
	"passkeeper/internal/storage"
)

var (
	ErrNoDecrypt = "can't decrypt password"
)

type BlobUsecaseInf interface {
	GetUserSalt(ctx context.Context, blobID string) string // return secret user string
	GetBlob(ctx context.Context, userID, blobID string) (cred *entities.CryptoBlob, err error)
	AddBlob(ctx context.Context, userID string, blob *entities.CryptoBlob) error
	UpdBlob(ctx context.Context, userID string, blob *entities.CryptoBlob) error
	DelBlob(ctx context.Context, userID, blobID string) error
	ListBlobs(ctx context.Context, userID string) (blobs []*entities.CryptoBlob, err error)

	//VerifyCredDate(cred *entities.Credential) // check date/time in received credential

	//EncryptPwd(ctx context.Context, password string) (encPwd string, err error)
	//DecryptPass(ctx context.Context, encPwd string) (password string, err error)
}

type BlobUsecase struct {
	ctx     context.Context
	storage storage.CredentialStorage
	//salt    string
	log *zerolog.Logger
}

type CredOption func(usecase *BlobUsecase)

func newBlobUC() *BlobUsecase {
	return &BlobUsecase{}
}

func NewBlobUCWithOpts(opts ...CredOption) *BlobUsecase {
	credsUC := newBlobUC()
	for _, opt := range opts {
		opt(credsUC)
	}

	return credsUC
}

func WithStorage(db storage.Storage) CredOption {
	return func(usecase *BlobUsecase) {
		usecase.storage = db
	}
}

func WithContext(ctx context.Context) CredOption {
	return func(usecase *BlobUsecase) {
		usecase.ctx = ctx
	}
}

func WithLogger(logger *zerolog.Logger) CredOption {
	return func(usecase *BlobUsecase) {
		usecase.log = logger
	}
}

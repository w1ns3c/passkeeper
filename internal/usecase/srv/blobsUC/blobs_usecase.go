package blobsUC

import (
	"context"

	"github.com/rs/zerolog"

	"passkeeper/internal/entities"
	"passkeeper/internal/storage"
)

// BlobUsecaseInf describe main blob functional on server side
type BlobUsecaseInf interface {
	GetBlob(ctx context.Context, userID, blobID string) (cred *entities.CryptoBlob, err error)
	AddBlob(ctx context.Context, userID string, blob *entities.CryptoBlob) error
	UpdBlob(ctx context.Context, userID string, blob *entities.CryptoBlob) error
	DelBlob(ctx context.Context, userID, blobID string) error
	ListBlobs(ctx context.Context, userID string) (blobs []*entities.CryptoBlob, err error)
}

// BlobUsecase implement BlobUsecaseInf
type BlobUsecase struct {
	ctx     context.Context
	storage storage.BlobStorage
	//salt    string
	log *zerolog.Logger
}

// BlobsOption new type to use in constructor of BlobUsecase
type BlobsOption func(usecase *BlobUsecase)

// newBlobUC is an empty constructor for BlobUsecase
func newBlobUC() *BlobUsecase {
	return &BlobUsecase{}
}

// NewBlobUCWithOpts is a constructor that can add BlobsOption to BlobUsecase
func NewBlobUCWithOpts(opts ...BlobsOption) *BlobUsecase {
	credsUC := newBlobUC()
	for _, opt := range opts {
		opt(credsUC)
	}

	return credsUC
}

// WithStorage add storage to BlobUsecase
func WithStorage(db storage.Storage) BlobsOption {
	return func(usecase *BlobUsecase) {
		usecase.storage = db
	}
}

// WithContext add context to BlobUsecase
func WithContext(ctx context.Context) BlobsOption {
	return func(usecase *BlobUsecase) {
		usecase.ctx = ctx
	}
}

// WithLogger add logger to BlobUsecase
func WithLogger(logger *zerolog.Logger) BlobsOption {
	return func(usecase *BlobUsecase) {
		usecase.log = logger
	}
}

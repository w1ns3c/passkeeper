package storage

import (
	"context"

	"passkeeper/internal/entities"
)

// Storage describe main storage functionality
type Storage interface {
	Init(ctx context.Context) error
	Close() error
	CheckConnection() error

	UserStorage
	BlobStorage
}

// UserStorage describe store that keep entities.User
type UserStorage interface {
	GetUserByID(cxt context.Context, userID string) (user *entities.User, err error)
	GetUserByLogin(cxt context.Context, login string) (user *entities.User, err error)
	CheckUserExist(ctx context.Context, login string) (exist bool, err error)
	SaveUser(ctx context.Context, u *entities.User) error
}

// BlobStorage describe store that keep entities.CryptoBlob
type BlobStorage interface {
	AddBlob(ctx context.Context, blob *entities.CryptoBlob) error
	GetBlob(ctx context.Context, userID, blobID string) (blob *entities.CryptoBlob, err error)
	GetAllBlobs(ctx context.Context, userID string) (blobs []*entities.CryptoBlob, err error)
	DeleteBlob(ctx context.Context, userID, blobID string) error
	UpdateBlob(ctx context.Context, blob *entities.CryptoBlob) error
}

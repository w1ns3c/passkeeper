package storage

import (
	"context"

	"passkeeper/internal/entities/structs"
)

// Storage describe main storage functionality
//
//go:generate mockgen -source storage.go -destination mocks/storage_mock.go -package=mocks
type Storage interface {
	Init(ctx context.Context) error
	Close() error
	CheckConnection() error

	UserStorage
	BlobStorage
}

// UserStorage describe store that keep entities.User
//
//go:generate mockgen -source internal/storage/storage.go -destination mocks/storage_mock.go -package=mocks
type UserStorage interface {
	GetUserByID(cxt context.Context, userID string) (user *structs.User, err error)
	GetUserByLogin(cxt context.Context, login string) (user *structs.User, err error)
	CheckUserExist(ctx context.Context, login string) (exist bool, err error)
	SaveUser(ctx context.Context, u *structs.User) error
}

// BlobStorage describe store that keep entities.CryptoBlob
//
//go:generate mockgen -source internal/storage/storage.go -destination mocks/storage_mock.go -package=mocks
type BlobStorage interface {
	AddBlob(ctx context.Context, blob *structs.CryptoBlob) error
	GetBlob(ctx context.Context, userID, blobID string) (blob *structs.CryptoBlob, err error)
	GetAllBlobs(ctx context.Context, userID string) (blobs []*structs.CryptoBlob, err error)
	DeleteBlob(ctx context.Context, userID, blobID string) error
	UpdateBlob(ctx context.Context, blob *structs.CryptoBlob) error
}

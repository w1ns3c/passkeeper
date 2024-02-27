package storage

import (
	"context"

	"github.com/w1nsec/passkeeper/internal/entities"
)

type Storage interface {
	Init(ctx context.Context) error
	UserStorage
	CredentialStorage
}

type UserStorage interface {
	CheckUserExist(ctx context.Context, login string) (exist bool, err error)
	RefreshToken(ctx context.Context, login string) error
	LoginUser(ctx context.Context, login, hash string) (user *entities.User, err error)
	SaveUser(ctx context.Context, u *entities.User) error
}

type CredentialStorage interface {
	AddCredential(ctx context.Context, userID string, password *entities.Credential) error
	GetCredential(ctx context.Context, userID, passwordID string) (password *entities.Credential, err error)
	GetAllCredentials(ctx context.Context, userID string) (passwords []*entities.Credential, err error)
	DeleteCredential(ctx context.Context, userID, passwordID string) error
	UpdateCredential(ctx context.Context, userID string, password *entities.Credential) error
}

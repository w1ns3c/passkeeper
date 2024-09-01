package storage

import (
	"context"

	"passkeeper/internal/entities"
)

type Storage interface {
	Init(ctx context.Context) error
	UserStorage
	CredentialStorage
}

type UserStorage interface {
	GetUserByID(cxt context.Context, userID string) (user *entities.User, err error)
	GetUserByLogin(cxt context.Context, login string) (user *entities.User, err error)
	CheckUserExist(ctx context.Context, login string) (exist bool, err error)
	SaveUser(ctx context.Context, u *entities.User) error

	//LoginUser(ctx context.Context, login, hash string) (user *entities.User, err error)
	//RefreshToken(ctx context.Context, login string) error
	//SaveChallenge(ctx context.Context, challenge string) error
}

type CredentialStorage interface {
	AddCredential(ctx context.Context, userID string, password *entities.CredBlob) error
	GetCredential(ctx context.Context, userID, passwordID string) (password *entities.CredBlob, err error)
	GetAllCredentials(ctx context.Context, userID string) (passwords []*entities.CredBlob, err error)
	DeleteCredential(ctx context.Context, userID, passwordID string) error
	UpdateCredential(ctx context.Context, userID string, password *entities.CredBlob) error
}

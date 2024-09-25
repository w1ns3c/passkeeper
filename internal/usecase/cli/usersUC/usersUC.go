package usersUC

import "context"

type UsersUsecase interface {
	Login(ctx context.Context, login, password string) (err error)
	Register(ctx context.Context, email, login, password, repeat string) error
	Logout()
	IsAuthed() bool
	GetToken() string
}

package myerrors

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (

	// token
	ErrNoTokenMsg = "no token in context"
	ErrNoToken    = status.Error(codes.Unauthenticated, ErrNoTokenMsg)

	ErrEmptyTokenMsg = "token is empty"
	ErrEmptyToken    = status.Error(codes.Unauthenticated, ErrEmptyTokenMsg)

	// DB
	ErrUserNotExist      = fmt.Errorf("user with this login not found")
	ErrWrongResultValues = errors.New("wrong count of results")
	ErrUsersWrongResult  = fmt.Errorf("can't return user: %v", ErrWrongResultValues)
	ErrBlobWrongResult   = fmt.Errorf("can't return blob: %v", ErrWrongResultValues)
	ErrRepoNotInit       = errors.New("repo not initialize")
	ErrDBConnect         = errors.New("can't connect to datebase")

	// Server
	ErrBlobsUserIDdifferent = errors.New("blob.UserID and userID in JWT are not the same")

	ErrGetUser = fmt.Errorf("can't get user by ID")
	ErrGenHash = fmt.Errorf("can't generate hash of password")

	ErrWrongOldPassword = fmt.Errorf("old password is wrong")
	ErrRepassNotSame    = fmt.Errorf("new pass and repeat not the same")

	ErrWrongPassword = fmt.Errorf("wrong password for username")
	ErrWrongUsername = fmt.Errorf("username not found")

	ErrPassIsEmpty   = fmt.Errorf("password is empty")
	ErrRePassIsEmpty = fmt.Errorf("password repeat is empty")
)

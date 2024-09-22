package myerrors

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUserNotFound = fmt.Errorf("user with this login not found")

	// token
	ErrNoTokenMsg = "no token in context"
	ErrNoToken    = status.Error(codes.Unauthenticated, ErrNoTokenMsg)

	ErrEmptyTokenMsg = "token is empty"
	ErrEmptyToken    = status.Error(codes.Unauthenticated, ErrEmptyTokenMsg)

	// DB
	ErrWrongResultValues = errors.New("wrong count of results")
	ErrUsersWrongResult  = fmt.Errorf("can't return user: %v", ErrWrongResultValues)
	ErrUsersNotExist     = errors.New("user not exist")
	ErrRepoNotInit       = errors.New("repo not initialize")
	ErrDBConnect         = errors.New("can't connect to datebase")

	// Server
	ErrBlobsUserIDdifferent = errors.New("blob.UserID and userID in JWT are not the same")
)

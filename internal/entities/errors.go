package entities

import (
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
)

package myerrors

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// Hashes
	ErrInvalidToken = fmt.Errorf("token sign is not valid")

	// token
	ErrNoTokenMsg = "no token in context"
	ErrNoToken    = status.Error(codes.Unauthenticated, ErrNoTokenMsg)

	ErrEmptyTokenMsg = "token is empty"
	ErrEmptyToken    = status.Error(codes.Unauthenticated, ErrEmptyTokenMsg)

	// DB storage
	ErrUserNotExist      = fmt.Errorf("user with this login not found")
	ErrWrongResultValues = errors.New("wrong count of results")
	ErrUsersWrongResult  = fmt.Errorf("can't return user: %v", ErrWrongResultValues)
	ErrBlobWrongResult   = fmt.Errorf("can't return blob: %v", ErrWrongResultValues)
	ErrRepoNotInit       = errors.New("repo not initialize")
	ErrDBConnect         = errors.New("can't connect to datebase")

	// Memory storage
	ErrBlobNotFound = fmt.Errorf("blob not exist")
	ErrUserNotFound = fmt.Errorf("user not exist")

	// Server
	ErrNoUCusers = fmt.Errorf("no users usecase")
	ErrNoUCcreds = fmt.Errorf("no creds usecase")

	// Server Usecase
	ErrBlobsUserIDdifferent = errors.New("blob.UserID and userID in JWT are not the same")

	ErrGetUser = fmt.Errorf("can't get user by ID")
	ErrGenHash = fmt.Errorf("can't generate hash of password")

	ErrWrongOldPassword = fmt.Errorf("old password is wrong")
	ErrRepassNotSame    = fmt.Errorf("new pass and repeat not the same")

	ErrWrongPassword = fmt.Errorf("wrong password for username")
	ErrWrongUsername = fmt.Errorf("username not found")

	ErrPassIsEmpty = fmt.Errorf("password is empty")

	// Server Transport
	ErrHndNotRegistered = fmt.Errorf("can't register some handlers")
	ErrNotEnoughOptions = fmt.Errorf("not enough options for grpc constructor")
	ErrWrongAuth        = "not authorized"

	// Server Transport Handlers
	ErrBlobAddMsg = "cred not added"
	ErrBlobAdd    = status.Error(codes.Internal, ErrBlobAddMsg)

	ErrBlobGetMsg = "cred can't get"
	ErrBlobGet    = status.Error(codes.Internal, ErrBlobGetMsg)

	ErrBlobUpdMsg = "cred not updated"
	ErrBlobUpd    = status.Error(codes.Internal, ErrBlobUpdMsg)

	ErrBlobDelMsg = "cred not deleted"
	ErrBlobDel    = status.Error(codes.Internal, ErrBlobDelMsg)

	ErrBlobListMsg = "creds not listed"
	ErrBlobList    = status.Error(codes.Internal, ErrBlobListMsg)

	ErrAlreadyExistMsg = "user already exist"
	ErrAlreadyExist    = status.Error(codes.AlreadyExists, ErrAlreadyExistMsg)

	ErrRegisterMsg = "can't register user"
	ErrRegister    = status.Error(codes.Internal, ErrRegisterMsg)

	ErrWrongLoginMsg = "can't login user, wrong login/password"
	ErrWrongLogin    = status.Errorf(codes.PermissionDenied, ErrWrongLoginMsg)

	// Client usecase
	ErrPassNotSame         = fmt.Errorf("pass and pass repeat are not the same")
	ErrInvalidEmail        = fmt.Errorf("email is not valid by new regexp")
	ErrPassDiff            = fmt.Errorf("passwords are not the same")
	ErrEmptyUsername       = fmt.Errorf("username is empty")
	ErrEmptyEmail          = fmt.Errorf("email is empty")
	ErrEmptyPassword       = fmt.Errorf("pass is empty")
	ErrEmptyPasswordRepeat = fmt.Errorf("pass repeat is empty")
)

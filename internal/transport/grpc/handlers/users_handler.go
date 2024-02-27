package handlers

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"github.com/w1nsec/passkeeper/internal/entities"
	pb "github.com/w1nsec/passkeeper/internal/transport/grpc/protofiles"
	"github.com/w1nsec/passkeeper/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// response errors
var (
	ErrAlreadyExistMsg = "user already exist"
	ErrAlreadyExist    = status.Error(codes.AlreadyExists, ErrAlreadyExistMsg)

	ErrRegisterMsg = "can't register user"
	ErrRegister    = status.Error(codes.Internal, ErrRegisterMsg)

	ErrWrongLoginMsg = "can't login user, wrong login/password"
	ErrWrongLogin    = status.Errorf(codes.PermissionDenied, ErrWrongLoginMsg)
)

type UsersHandler struct {
	pb.UnimplementedUserSvcServer
	service usecase.UserUsecaseInf
	log     *zerolog.Logger
}

//rpc RegisterUser(UserRegisterRequest) returns (UserRegisterResponse);
//rpc LoginUser(UserLoginRequest) returns (UserLoginResponse);

func (h *UsersHandler) RegisterUser(ctx context.Context, request *pb.UserRegisterRequest) (resp *pb.UserRegisterResponse, err error) {
	token, err := h.service.RegisterUser(ctx, request.Login, request.Password, request.RePassword)
	if err != nil {
		if !errors.Is(err, entities.ErrUserNotFound) {
			h.log.Error().
				Err(err).Msg(ErrAlreadyExistMsg)
			return nil, ErrAlreadyExist
		}

		h.log.Error().
			Err(err).Msg(ErrRegisterMsg)
		return nil, ErrRegister
	}

	resp = &pb.UserRegisterResponse{
		Token: token,
	}

	return resp, nil
}

func (h *UsersHandler) LoginUser(ctx context.Context, request *pb.UserLoginRequest) (resp *pb.UserLoginResponse, err error) {
	token, err := h.service.LoginUser(ctx, request.Login, request.Password)
	if err != nil {
		h.log.Error().
			Err(err).Msg(ErrWrongLoginMsg)

		return nil, ErrWrongLogin
	}

	resp = &pb.UserLoginResponse{
		Token: token,
	}

	return resp, nil
}

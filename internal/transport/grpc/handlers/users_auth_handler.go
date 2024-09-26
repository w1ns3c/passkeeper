package handlers

import (
	"context"
	"errors"

	"github.com/rs/zerolog"

	"passkeeper/internal/entities/myerrors"
	pb "passkeeper/internal/transport/grpc/protofiles/proto"
	"passkeeper/internal/usecase/srv/usersUC"
)

// UsersHandler handle user auth requests
//
//go:generate mockgen -source=../protofiles/proto/users_service_grpc.pb.go -destination=../../../../mocks/gservice/user_auth.go -package=mocks
type UsersHandler struct {
	pb.UnimplementedUserSvcServer
	service usersUC.UserUsecaseInf
	log     *zerolog.Logger
}

// NewUsersHandler is a constructor for UsersHandler
func NewUsersHandler(logger *zerolog.Logger, service usersUC.UserUsecaseInf) *UsersHandler {
	return &UsersHandler{
		UnimplementedUserSvcServer: pb.UnimplementedUserSvcServer{},
		service:                    service,
		log:                        logger,
	}
}

// RegisterUser handle user register request
func (h *UsersHandler) RegisterUser(ctx context.Context, request *pb.UserRegisterRequest) (resp *pb.UserRegisterResponse, err error) {
	token, secret, err := h.service.RegisterUser(ctx, request.Login, request.Password, request.RePassword)
	if err != nil {
		if errors.Is(err, myerrors.ErrAlreadyExist) {
			h.log.Error().
				Err(err).Msg(myerrors.ErrAlreadyExistMsg)

			return nil, myerrors.ErrAlreadyExist
		}

		h.log.Error().
			Err(err).Msg(myerrors.ErrRegisterMsg)

		return nil, myerrors.ErrRegister
	}

	resp = &pb.UserRegisterResponse{
		Token:     token,
		SrvSecret: secret,
	}
	h.log.Info().
		Msgf("User \"%s\" successfully registered!", request.Login)

	return resp, nil
}

// LoginUser handle user login request
func (h *UsersHandler) LoginUser(ctx context.Context,
	req *pb.UserLoginRequest) (resp *pb.UserLoginResponse, err error) {

	token, secret, err := h.service.LoginUser(ctx, req.Login, req.Password)
	if err != nil {
		h.log.Error().
			Err(err).Msg(myerrors.ErrWrongLoginMsg)

		return nil, myerrors.ErrWrongLogin
	}

	resp = &pb.UserLoginResponse{
		Token:     token,
		SrvSecret: secret,
	}
	h.log.Info().
		Msgf("User \"%s\" successfully logged in!", req.Login)

	return resp, nil
}

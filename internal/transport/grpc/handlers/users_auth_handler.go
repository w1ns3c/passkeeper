package handlers

import (
	"context"
	"errors"

	"github.com/rs/zerolog"

	errors2 "passkeeper/internal/entities/myerrors"
	pb "passkeeper/internal/transport/grpc/protofiles/proto"
	"passkeeper/internal/usecase/srv/usersUC"
)

// UsersHandler handle user auth requests
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
		if !errors.Is(err, errors2.ErrUserNotExist) {
			h.log.Error().
				Err(err).Msg(errors2.ErrAlreadyExistMsg)
			return nil, errors2.ErrAlreadyExist
		}

		h.log.Error().
			Err(err).Msg(errors2.ErrRegisterMsg)
		return nil, errors2.ErrRegister
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
			Err(err).Msg(errors2.ErrWrongLoginMsg)

		return nil, errors2.ErrWrongLogin
	}

	resp = &pb.UserLoginResponse{
		Token:     token,
		SrvSecret: secret,
	}
	h.log.Info().
		Msgf("User \"%s\" successfully logged in!", req.Login)

	return resp, nil
}

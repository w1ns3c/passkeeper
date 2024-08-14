package handlers

import (
	"context"
	"errors"
	"github.com/w1ns3c/passkeeper/internal/config"
	"github.com/w1ns3c/passkeeper/internal/usecase/srv"
	"google.golang.org/grpc/metadata"

	"github.com/rs/zerolog"
	"github.com/w1ns3c/passkeeper/internal/entities"
	pb "github.com/w1ns3c/passkeeper/internal/transport/grpc/protofiles/proto"
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

	ErrGenChallengeMsg = "can't generate challenge"
	ErrGenChallenge    = status.Errorf(codes.Internal, ErrGenChallengeMsg)
)

type UsersHandler struct {
	pb.UnimplementedUserSvcServer
	service srv.UserUsecaseInf
	log     *zerolog.Logger
}

//rpc RegisterUser(UserRegisterRequest) returns (UserRegisterResponse);
//rpc LoginUser(UserLoginRequest) returns (UserLoginResponse);

func (h *UsersHandler) RegisterUser(ctx context.Context, request *pb.UserRegisterRequest) (resp *pb.UserRegisterResponse, err error) {
	token, secret, err := h.service.RegisterUser(ctx, request.Login, request.Password, request.RePassword)
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

	ctx = metadata.AppendToOutgoingContext(ctx, config.TokenHeader, token)

	resp = &pb.UserRegisterResponse{
		Secret: secret,
	}

	return resp, nil
}

func (h *UsersHandler) LoginUser(ctx context.Context,
	req *pb.UserLoginRequest) (resp *pb.UserLoginResponse, err error) {

	token, secret, err := h.service.LoginUser(ctx, req.Login, req.Password)
	if err != nil {
		h.log.Error().
			Err(err).Msg(ErrWrongLoginMsg)

		return nil, ErrWrongLogin
	}

	ctx = metadata.AppendToOutgoingContext(ctx, config.TokenHeader, token)

	resp = &pb.UserLoginResponse{
		//Token: token,
		Secret: secret,
	}

	return resp, nil

}

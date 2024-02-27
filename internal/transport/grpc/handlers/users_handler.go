package handlers

import (
	"context"
	"errors"
	"github.com/w1nsec/passkeeper/internal/entities"
	pb "github.com/w1nsec/passkeeper/internal/transport/grpc/protofiles"
	"github.com/w1nsec/passkeeper/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UsersHandler struct {
	pb.UnimplementedUserSvcServer
	service usecase.UserUsecase
}

//rpc RegisterUser(UserRegisterRequest) returns (UserRegisterResponse);
//rpc LoginUser(UserLoginRequest) returns (UserLoginResponse);

func (s *UsersHandler) RegisterUser(ctx context.Context, request *pb.UserRegisterRequest) (resp *pb.UserRegisterResponse, err error) {
	token, err := s.service.RegisterUser(ctx, request.Login, request.Password, request.RePassword)
	if err != nil {
		if !errors.Is(err, entities.ErrUserNotFound) {
			return nil, status.Errorf(codes.AlreadyExists, "user already exist: %v", err)
		}
		return nil, status.Errorf(codes.Canceled, "can't register user: %v", err)
	}

	resp = &pb.UserRegisterResponse{
		// TODO change Err to Token
		Err: token,
	}

	return resp, nil
}

func (s *UsersHandler) LoginUser(ctx context.Context, request *pb.UserLoginRequest) (resp *pb.UserLoginResponse, err error) {
	token, err := s.service.LoginUser(ctx, request.Login, request.Password)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "can't register user: %v", err)

	}

	resp = &pb.UserLoginResponse{
		Token: token,
	}

	return resp, nil
}

package handlers

import (
	"context"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/metadata"

	"passkeeper/internal/entities/myerrors"
	pb "passkeeper/internal/transport/grpc/protofiles/proto"
	"passkeeper/internal/usecase/srv/usersUC"
	mocks "passkeeper/mocks/gservice"
	mocksusecase "passkeeper/mocks/usecase/users_usecase"
)

func TestNewUsersHandler(t *testing.T) {
	type args struct {
		service usersUC.UserUsecaseInf
	}

	var (
		logger = zerolog.New(io.Discard)
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocksusecase.NewMockUserUsecaseInf(ctrl)

	tests := []struct {
		name string
		args args
		want *UsersHandler
	}{
		{
			name: "Test NewUsersHandler 1: Nil service",
			args: args{
				service: nil,
			},
			want: &UsersHandler{
				UnimplementedUserSvcServer: pb.UnimplementedUserSvcServer{},
				log:                        &logger,
				service:                    nil,
			},
		},

		{
			name: "Test NewUsersHandler 2: All valid",
			args: args{
				service: m,
			},
			want: &UsersHandler{
				UnimplementedUserSvcServer: pb.UnimplementedUserSvcServer{},
				log:                        &logger,
				service:                    m,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUsersHandler(&logger, tt.args.service); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUsersHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsersHandler_LoginUser(t *testing.T) {

	var (
		logger = zerolog.New(io.Discard)
	)

	type args struct {
		ctx context.Context
		req *pb.UserLoginRequest
	}
	tests := []struct {
		name     string
		args     args
		prepare  func(server *mocks.MockUserSvcServer, uc *mocksusecase.MockUserUsecaseInf)
		wantResp *pb.UserLoginResponse
		wantErr  bool
	}{
		{
			name: "LoginUser: Valid Login",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(),
					metadata.New(map[string]string{})),
				req: &pb.UserLoginRequest{
					Login:    "user1",
					Password: "password1",
				},
			},
			prepare: func(server *mocks.MockUserSvcServer, uc *mocksusecase.MockUserUsecaseInf) {
				uc.EXPECT().LoginUser(gomock.Any(), "user1", "password1").
					Return("token", "secret", nil)
			},
			wantResp: &pb.UserLoginResponse{
				Token:     "token",
				SrvSecret: "secret",
			},
			wantErr: false,
		},
		{
			name: "LoginUser: invalid user",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(),
					metadata.New(map[string]string{})),
				req: &pb.UserLoginRequest{
					Login:    "user12",
					Password: "password1",
				},
			},
			prepare: func(server *mocks.MockUserSvcServer, uc *mocksusecase.MockUserUsecaseInf) {
				uc.EXPECT().LoginUser(gomock.Any(), "user12", "password1").
					Return("", "", myerrors.ErrUserNotExist)
			},
			wantResp: nil,

			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srv := mocks.NewMockUserSvcServer(ctrl)
			usecase := mocksusecase.NewMockUserUsecaseInf(ctrl)

			if &tt.prepare != nil {
				tt.prepare(srv, usecase)
			}

			h := &UsersHandler{
				service: usecase,
				log:     &logger,
			}
			gotResp, err := h.LoginUser(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("LoginUser() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func TestUsersHandler_RegisterUser(t *testing.T) {
	var (
		logger = zerolog.New(io.Discard)
	)

	type args struct {
		ctx context.Context
		req *pb.UserRegisterRequest
	}

	tests := []struct {
		name     string
		args     args
		prepare  func(server *mocks.MockUserSvcServer, uc *mocksusecase.MockUserUsecaseInf)
		wantResp *pb.UserRegisterResponse
		wantErr  error
	}{
		{
			name: "Register User: valid register",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(),
					metadata.New(map[string]string{})),
				req: &pb.UserRegisterRequest{
					Login:      "user12",
					Password:   "password1",
					RePassword: "password1",
				},
			},
			prepare: func(server *mocks.MockUserSvcServer, uc *mocksusecase.MockUserUsecaseInf) {
				gomock.InOrder(
					uc.EXPECT().RegisterUser(gomock.Any(), "user12", "password1", "password1").
						Return("token", "cryptosecret", nil),
				)
			},
			wantResp: &pb.UserRegisterResponse{
				Token:     "token",
				SrvSecret: "cryptosecret",
			},
			wantErr: nil,
		},
		{
			name: "Register User: user exist",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(),
					metadata.New(map[string]string{})),
				req: &pb.UserRegisterRequest{
					Login:      "exist",
					Password:   "password1",
					RePassword: "password1",
				},
			},
			prepare: func(server *mocks.MockUserSvcServer, uc *mocksusecase.MockUserUsecaseInf) {
				gomock.InOrder(
					uc.EXPECT().RegisterUser(gomock.Any(), "exist", "password1", "password1").
						Return("", "", myerrors.ErrAlreadyExist),
				)
			},
			wantResp: nil,
			wantErr:  myerrors.ErrAlreadyExist,
		},
		{
			name: "Register User: user pass not the same",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(),
					metadata.New(map[string]string{})),
				req: &pb.UserRegisterRequest{
					Login:      "exist",
					Password:   "password1",
					RePassword: "123",
				},
			},
			prepare: func(server *mocks.MockUserSvcServer, uc *mocksusecase.MockUserUsecaseInf) {
				gomock.InOrder(
					uc.EXPECT().RegisterUser(gomock.Any(), "exist", "password1", "123").
						Return("", "", myerrors.ErrRepassNotSame),
				)
			},
			wantResp: nil,
			wantErr:  myerrors.ErrRegister,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocksusecase.NewMockUserUsecaseInf(ctrl)

			if &tt.prepare != nil {
				tt.prepare(nil, service)
			}

			h := &UsersHandler{
				service: service,
				log:     &logger,
			}

			gotResp, err := h.RegisterUser(tt.args.ctx, tt.args.req)
			if err != nil {
				if !errors.Is(tt.wantErr, err) {
					t.Errorf("RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("RegisterUser() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

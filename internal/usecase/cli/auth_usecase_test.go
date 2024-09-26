package cli

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"passkeeper/internal/entities/hashes"
	"passkeeper/internal/entities/myerrors"
	"passkeeper/internal/entities/structs"
	pb "passkeeper/internal/transport/grpc/protofiles/proto"
	mocks "passkeeper/mocks/gservice"
)

func TestFilterEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "Valid simple email",
			email:   "email@email.com",
			wantErr: false,
		},
		{
			name:    "Valid hard email",
			email:   "email44.11-email.me@email99.corp.com",
			wantErr: false,
		},
		{
			name:    "Invalid email (dot end username)",
			email:   "email.@email.com",
			wantErr: true,
		},
		{
			name:    "Invalid email (short email code)",
			email:   "email@email.c",
			wantErr: true,
		},
		{
			name:    "Invalid email (double dots username)",
			email:   "email..email@email.com",
			wantErr: true,
		},
		{
			name:    "Invalid email (double dots email)",
			email:   "email@email..c",
			wantErr: true,
		},
		{
			name:    "Invalid email (double @)",
			email:   "email@alex@email.com",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := FilterEmail(tt.email); (err != nil) != tt.wantErr {
				t.Errorf("FilterEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientUC_Login(t *testing.T) {
	type args struct {
		login    string
		password string
	}

	tests := []struct {
		name    string
		args    args
		prepare func(*mocks.MockUserSvcClient)
		wantErr bool
	}{
		{
			name: "Test Login 1: success",
			args: args{
				login:    "login",
				password: "password",
			},
			prepare: func(cli *mocks.MockUserSvcClient) {
				cli.EXPECT().LoginUser(gomock.Any(), &pb.UserLoginRequest{
					Login:    "login",
					Password: hashes.Hash("password"),
				}).Return(&pb.UserLoginResponse{
					Token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc0MjM0NjksIlVzZXJJRCI6ImMwMjIwNzUxNGIzYmJhZjQ3OGE4Yjg1MDIyMTAxNzk1In0.U_BjCncYFhovh_lB9JZmNs69luT7UiH71WAtCBCnHEs",
					SrvSecret: "12a8c6e3494a48ea60d8890b8714e8230854a94b98280db465248966d9aa2a90ee983d3027abbafbaa4b0ddfdba8ca473e1a99ce8ac198d07d4bbbf7",
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "Test Login 2: invalid pass",
			args: args{
				login:    "login",
				password: "invalid_pass",
			},
			prepare: func(cli *mocks.MockUserSvcClient) {
				cli.EXPECT().LoginUser(gomock.Any(), &pb.UserLoginRequest{
					Login:    "login",
					Password: hashes.Hash("invalid_pass"),
				}).Return(nil, fmt.Errorf("new error"))
			},
			wantErr: true,
		},
		{
			name: "Test Login 3: can't extract userID",
			args: args{
				login:    "login",
				password: "password",
			},
			prepare: func(cli *mocks.MockUserSvcClient) {
				cli.EXPECT().LoginUser(gomock.Any(), &pb.UserLoginRequest{
					Login:    "login",
					Password: hashes.Hash("password"),
				}).Return(&pb.UserLoginResponse{
					Token:     "wrong_token",
					SrvSecret: "",
				}, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var (
				mock    = mocks.NewMockUserSvcClient(ctrl)
				usecase = &ClientUC{}
			)

			if tt.prepare != nil {
				tt.prepare(mock)
			}

			usecase.userSvc = mock

			if err := usecase.Login(context.Background(), tt.args.login, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientUC_Register(t *testing.T) {
	type args struct {
		email    string
		login    string
		password string
		repeat   string
	}

	tests := []struct {
		name    string
		args    args
		prepare func(*mocks.MockUserSvcClient)
		wantErr error
	}{
		{
			name: "Test Register 1: success",
			args: args{
				email:    "email@email.com",
				login:    "login",
				password: "password",
				repeat:   "password",
			},
			prepare: func(cli *mocks.MockUserSvcClient) {
				cli.EXPECT().RegisterUser(gomock.Any(), &pb.UserRegisterRequest{
					Email:      "email@email.com",
					Login:      "login",
					Password:   hashes.Hash("password"),
					RePassword: hashes.Hash("password"),
				}).Return(&pb.UserRegisterResponse{}, nil)
			},
			wantErr: nil,
		},
		{
			name: "Test Login 2: exist user",
			args: args{
				email:    "email@email.com",
				login:    "exist_user",
				password: "password",
				repeat:   "password",
			},
			prepare: func(cli *mocks.MockUserSvcClient) {
				cli.EXPECT().RegisterUser(gomock.Any(), &pb.UserRegisterRequest{
					Email:      "email@email.com",
					Login:      "exist_user",
					Password:   hashes.Hash("password"),
					RePassword: hashes.Hash("password"),
				}).Return(&pb.UserRegisterResponse{}, myerrors.ErrAlreadyExist)
			},
			wantErr: myerrors.ErrAlreadyExist,
		},
		{
			name: "Test Login 3: passwords are not the same",
			args: args{
				email:    "email@email.com",
				login:    "exist_user",
				password: "password",
				repeat:   "other_password",
			},
			wantErr: myerrors.ErrPassNotSame,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var (
				mock    = mocks.NewMockUserSvcClient(ctrl)
				usecase = &ClientUC{}
			)

			if tt.prepare != nil {
				tt.prepare(mock)
			}

			usecase.userSvc = mock

			err := usecase.Register(context.Background(),
				tt.args.email, tt.args.login, tt.args.password, tt.args.repeat)

			if err != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientUC_FilterUserRegValues(t *testing.T) {
	type args struct {
		username   string
		password   string
		passRepeat string
		email      string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Test 1: empty email",
			args: args{
				username:   "username",
				password:   "password",
				passRepeat: "password",
				email:      "",
			},
			wantErr: myerrors.ErrEmptyEmail,
		},
		{
			name: "Test 2: empty password",
			args: args{
				username:   "username",
				password:   "",
				passRepeat: "password",
				email:      "username@email.com",
			},
			wantErr: myerrors.ErrEmptyPassword,
		},
		{
			name: "Test 3: empty password repeat",
			args: args{
				username:   "username",
				password:   "password",
				passRepeat: "",
				email:      "username@email.com",
			},
			wantErr: myerrors.ErrEmptyPasswordRepeat,
		},
		{
			name: "Test 4: empty username",
			args: args{
				username:   "",
				password:   "password",
				passRepeat: "password",
				email:      "username@email.com",
			},
			wantErr: myerrors.ErrEmptyUsername,
		},
		{
			name: "Test 5: wrong email",
			args: args{
				username:   "username",
				password:   "password",
				passRepeat: "password",
				email:      "usern ame@email.com",
			},
			wantErr: myerrors.ErrInvalidEmail,
		},
		{
			name: "Test 6: short password",
			args: args{
				username:   "username",
				password:   "pass",
				passRepeat: "pass",
				email:      "username@email.com",
			},
			wantErr: myerrors.ErrMinPasswordLen,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClientUC{}

			if err := c.FilterUserRegValues(tt.args.username,
				tt.args.password, tt.args.passRepeat, tt.args.email); err != tt.wantErr {
				t.Errorf("FilterUserRegValues() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientUC_Logout(t *testing.T) {

	c := &ClientUC{
		Authed: true,
		User:   &structs.User{},
		Token:  "some token",

		Creds: []*structs.Credential{},
		Cards: []*structs.Card{},
		Notes: []*structs.Note{},
		Files: []*structs.File{},

		viewPageFocus: false,

		m: &sync.RWMutex{},
	}

	c.Logout()

	require.Equal(t, c.Authed, false)
	require.Nil(t, c.User)

	require.Nil(t, c.Creds)
	require.Nil(t, c.Cards)
	require.Nil(t, c.Notes)
	require.Nil(t, c.Files)

	require.Equal(t, c.viewPageFocus, true)
}

func TestClientUC_IsAuthed(t *testing.T) {
	c := &ClientUC{
		User: &structs.User{},
		m:    &sync.RWMutex{},
	}

	require.Equal(t, c.IsAuthed(), true)

	c.Logout()
	require.Equal(t, c.IsAuthed(), false)
}

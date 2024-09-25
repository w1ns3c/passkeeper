package memstorage

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"passkeeper/internal/entities/structs"
)

func TestMemStorage_CheckUserExist(t *testing.T) {

	type args struct {
		users map[string]*structs.User
		login string
	}
	tests := []struct {
		name      string
		args      args
		wantExist bool
		wantErr   bool
	}{
		{
			name: "TestCheckUserExist 1: Exist",
			args: args{
				users: map[string]*structs.User{"test": {
					Login: "test",
				}},
				login: "test",
			},
			wantExist: true,
			wantErr:   false,
		},
		{
			name: "TestCheckUserExist 2: Not exist",
			args: args{
				users: map[string]*structs.User{"test": {
					Login: "test",
				}},
				login: "test1",
			},
			wantExist: false,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				users: tt.args.users,
			}
			m.Init(context.Background())

			gotExist, err := m.CheckUserExist(context.Background(), tt.args.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckUserExist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotExist != tt.wantExist {
				t.Errorf("CheckUserExist() gotExist = %v, want %v", gotExist, tt.wantExist)
			}
		})
	}
}

func TestMemStorage_GetUserByID(t *testing.T) {

	type args struct {
		users  map[string]*structs.User
		userID string
	}

	ctx := context.Background()

	tests := []struct {
		name     string
		args     args
		wantUser *structs.User
		wantErr  bool
	}{
		{
			name: "Test_GetUserByID 1: Exist",
			args: args{
				users: map[string]*structs.User{"1111": {
					ID:     "test",
					Login:  "1111",
					Hash:   "123123123",
					Phone:  "8888888888888",
					Email:  "test@test.com",
					Salt:   "salt",
					Secret: "secret",
				}},
				userID: "test",
			},
			wantUser: &structs.User{
				ID:     "test",
				Login:  "1111",
				Hash:   "123123123",
				Phone:  "8888888888888",
				Email:  "test@test.com",
				Salt:   "salt",
				Secret: "secret",
			},
			wantErr: false,
		},
		{
			name: "Test_GetUserByID 2: NOT Exist",
			args: args{
				users: map[string]*structs.User{"1111": {
					ID:     "test",
					Login:  "1111",
					Hash:   "123123123",
					Phone:  "8888888888888",
					Email:  "test@test.com",
					Salt:   "salt",
					Secret: "secret",
				}},
				userID: "test111",
			},
			wantUser: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				users: tt.args.users,
			}
			m.Init(ctx)

			gotUser, err := m.GetUserByID(ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("GetUserByID() gotUser = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}

func TestMemStorage_GetUserByLogin(t *testing.T) {
	type args struct {
		users map[string]*structs.User
		login string
	}

	ctx := context.Background()

	tests := []struct {
		name     string
		args     args
		wantUser *structs.User
		wantErr  bool
	}{
		{
			name: "Test_GetUserByID 1: Exist",
			args: args{
				users: map[string]*structs.User{"1111": {
					ID:     "test",
					Login:  "1111",
					Hash:   "123123123",
					Phone:  "8888888888888",
					Email:  "test@test.com",
					Salt:   "salt",
					Secret: "secret",
				}},
				login: "1111",
			},
			wantUser: &structs.User{
				ID:     "test",
				Login:  "1111",
				Hash:   "123123123",
				Phone:  "8888888888888",
				Email:  "test@test.com",
				Salt:   "salt",
				Secret: "secret",
			},
			wantErr: false,
		},
		{
			name: "Test_GetUserByID 2: NOT Exist",
			args: args{
				users: map[string]*structs.User{"1111": {
					ID:     "test",
					Login:  "1111",
					Hash:   "123123123",
					Phone:  "8888888888888",
					Email:  "test@test.com",
					Salt:   "salt",
					Secret: "secret",
				}},
				login: "test111",
			},
			wantUser: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				users: tt.args.users,
			}
			m.Init(ctx)

			gotUser, err := m.GetUserByLogin(ctx, tt.args.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("GetUserByLogin() gotUser = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}

func TestMemStorage_SaveUser(t *testing.T) {
	type args struct {
		users      map[string]*structs.User
		savingUser *structs.User
	}

	ctx := context.Background()

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test_SaveUser 1: valid",
			args: args{
				users: map[string]*structs.User{"1111": {
					ID:     "test",
					Login:  "1111",
					Hash:   "123123123",
					Phone:  "8888888888888",
					Email:  "test@test.com",
					Salt:   "salt",
					Secret: "secret",
				}},
				savingUser: &structs.User{
					ID:     "test123",
					Login:  "1111555",
					Hash:   "12312312311",
					Phone:  "8888888888888",
					Email:  "test@test.com",
					Salt:   "salt",
					Secret: "secret",
				},
			},

			wantErr: false,
		},
		{
			name: "Test_SaveUser 2: user exist",
			args: args{
				users: map[string]*structs.User{"1111": {
					ID:     "test",
					Login:  "1111",
					Hash:   "123123123",
					Phone:  "8888888888888",
					Email:  "test@test.com",
					Salt:   "salt",
					Secret: "secret",
				}},
				savingUser: &structs.User{
					ID:     "test",
					Login:  "1111",
					Hash:   "123123123",
					Phone:  "8888888888888",
					Email:  "test@test.com",
					Salt:   "salt",
					Secret: "secret",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				users: tt.args.users,
			}
			m.Init(ctx)

			err := m.SaveUser(ctx, tt.args.savingUser)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			savedUser, err := m.GetUserByLogin(ctx, tt.args.savingUser.Login)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveUser() error = %v, can't return saved user", err)
			}

			require.Equal(t, savedUser, tt.args.savingUser)

		})
	}
}

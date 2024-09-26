package handlers

import (
	"context"
	"io"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/metadata"

	"passkeeper/internal/entities/config"
	pb "passkeeper/internal/transport/grpc/protofiles/proto"
	"passkeeper/internal/usecase/srv/usersUC"
	mocks "passkeeper/mocks/gservice"
	mocksusecase "passkeeper/mocks/usecase/users_usecase"
)

func TestNewUserChangePassHandler(t *testing.T) {
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
		want *UserChangePassHandler
	}{
		{
			name: "Test NewUserChangePass 1: Nil service",
			args: args{
				service: nil,
			},
			want: &UserChangePassHandler{
				UnimplementedUserChangePassSvcServer: pb.UnimplementedUserChangePassSvcServer{},
				log:                                  &logger,
				service:                              nil,
			},
		},

		{
			name: "Test NewUserChangePass 2: All valid",
			args: args{
				service: m,
			},
			want: &UserChangePassHandler{
				UnimplementedUserChangePassSvcServer: pb.UnimplementedUserChangePassSvcServer{},
				log:                                  &logger,
				service:                              m,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserChangePassHandler(&logger, tt.args.service); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserChangePassHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserChangePassHandler_ChangePass(t *testing.T) {

	var (
		ctx1 = metadata.NewIncomingContext(context.Background(),
			metadata.New(map[string]string{config.TokenHeader: "userID1"}))
		emptyCtx = metadata.NewIncomingContext(context.Background(),
			metadata.New(map[string]string{}))

		logger = zerolog.New(io.Discard)
	)
	tests := []struct {
		name    string
		ctx     context.Context
		req     *pb.UserChangePassReq
		prepare func(*mocks.MockUserChangePassSvcServer, *mocksusecase.MockUserUsecaseInf)
		want    *empty.Empty
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "ChangePass all valid",
			ctx:  ctx1,
			req: &pb.UserChangePassReq{
				OldPass: "old",
				NewPass: "new",
				Repeat:  "new",
			},
			prepare: func(srv *mocks.MockUserChangePassSvcServer, uc *mocksusecase.MockUserUsecaseInf) {
				gomock.InOrder(
					uc.EXPECT().ChangePassword(ctx1, "userID1", "old", "new", "new").
						Return(nil),
				)
			},
			wantErr: false,
		},
		{
			name: "ChangePass without token",
			ctx:  emptyCtx,
			req: &pb.UserChangePassReq{
				OldPass: "old",
				NewPass: "new",
				Repeat:  "new",
			},
			prepare: func(srv *mocks.MockUserChangePassSvcServer, uc *mocksusecase.MockUserUsecaseInf) {
				gomock.InOrder()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srv := mocks.NewMockUserChangePassSvcServer(ctrl)
			usecase := mocksusecase.NewMockUserUsecaseInf(ctrl)

			if &tt.prepare != nil {
				tt.prepare(srv, usecase)
			}

			h := &UserChangePassHandler{
				//UnimplementedUserChangePassSvcServer: srv,
				service: usecase,
				log:     &logger,
			}

			_, err := h.ChangePass(tt.ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangePass() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				return
			}

		})
	}
}

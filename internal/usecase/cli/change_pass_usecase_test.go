package cli

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"passkeeper/internal/entities/hashes"
	pb "passkeeper/internal/transport/grpc/protofiles/proto"
	mocks "passkeeper/mocks/gservice"
)

func TestClientUC_ChangePass(t *testing.T) {
	type args struct {
		curPass string
		newPass string
		repeat  string
	}

	tests := []struct {
		name    string
		args    args
		prepare func(*mocks.MockUserChangePassSvcClient)
		wantErr bool
	}{
		{
			name: "Test ChangePass 1: success",
			args: args{
				curPass: "123123",
				newPass: "newpass",
				repeat:  "newpass",
			},
			prepare: func(cli *mocks.MockUserChangePassSvcClient) {
				cli.EXPECT().ChangePass(gomock.Any(), &pb.UserChangePassReq{
					OldPass: hashes.Hash("123123"),
					NewPass: hashes.Hash("newpass"),
					Repeat:  hashes.Hash("newpass"),
				}).Return(nil, nil)
			},
			wantErr: false,
		},
		{
			name: "Test ChangePass 2: pass not the same",
			args: args{
				curPass: "123123",
				newPass: "newpass",
				repeat:  "anotherpassword",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var (
				mock    = mocks.NewMockUserChangePassSvcClient(ctrl)
				usecase = &ClientUC{}
			)

			if tt.prepare != nil {
				tt.prepare(mock)
			}

			usecase.passSvc = mock

			if err := usecase.ChangePass(context.Background(),
				tt.args.curPass, tt.args.newPass, tt.args.repeat); (err != nil) != tt.wantErr {
				t.Errorf("ChangePass() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

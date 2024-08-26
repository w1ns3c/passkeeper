package cli

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "passkeeper/internal/transport/grpc/protofiles/proto"
)

var (
	ErrPassNotSame = fmt.Errorf("pass and pass repeat are not the same")
)

type ClientUsecase interface {
	Login(ctx context.Context, login, password string) error
	Register(ctx context.Context, email, login, password, repeat string) error
}

type ClientUC struct {
	ctx        context.Context
	Token      string
	SecretHash string

	userSvc pb.UserSvcClient
	passSvc pb.UserChangePassSvcClient
}

func NewClientUC(addr string) (cli *ClientUC, err error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &ClientUC{
		userSvc: pb.NewUserSvcClient(conn),
		passSvc: pb.NewUserChangePassSvcClient(conn),
	}, nil
}

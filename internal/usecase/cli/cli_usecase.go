package cli

import (
	"context"
	"fmt"
	pb "github.com/w1ns3c/passkeeper/internal/transport/grpc/protofiles/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ErrPassNotSame = fmt.Errorf("pass and pass repeat are not the same")
)

type ClientUsecase interface {
	Login(ctx context.Context, login, password string) error
	Register(ctx context.Context, login, password, repeat, email string) error
}

type ClientUC struct {
	ctx        context.Context
	Token      string
	SecretHash string

	userSvc pb.UserSvcClient
	passSvc pb.UserPassSvcClient
}

func NewClientUC(addr string) (cli *ClientUC, err error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &ClientUC{
		userSvc: pb.NewUserSvcClient(conn),
		passSvc: pb.NewUserPassSvcClient(conn),
	}, nil
}

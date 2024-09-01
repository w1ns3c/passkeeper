package cli

import (
	"context"
	"fmt"
	"passkeeper/internal/entities"

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
	GetCreds(ctx context.Context) (creds []*entities.Credential, err error)
}

type ClientUC struct {
	ctx    context.Context
	Token  string // JWT token
	UserID string // userID from JWT token
	Secret string // full secret for decrypt user's creds, contains sha1(clear_pass+secret_from_server)

	userSvc  pb.UserSvcClient
	passSvc  pb.UserChangePassSvcClient
	credsSvc pb.CredSvcClient
}

func NewClientUC(addr string) (cli *ClientUC, err error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &ClientUC{
		userSvc:  pb.NewUserSvcClient(conn),
		passSvc:  pb.NewUserChangePassSvcClient(conn),
		credsSvc: pb.NewCredSvcClient(conn),
	}, nil
}

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
	Login(ctx context.Context, login, password string) (secret, userID string, err error)
	Register(ctx context.Context, email, login, password, repeat string) error
	Logout()

	GetCreds(ctx context.Context) (creds []*entities.Credential, err error)
	EditCred(ctx context.Context, cred *entities.Credential) (err error)
	AddCred(ctx context.Context, cred *entities.Credential) (err error)
	DelCred(ctx context.Context, credID string) (err error)

	// moved from tuiApp
	GetToken() string
	GetHeader() string
}

type ClientUC struct {
	Token       string // JWT token
	UserID      string // userID from JWT token
	CredsSecret string // full secret for decrypt user's creds, contains sha1(clear_pass+secret_from_server)

	userSvc  pb.UserSvcClient
	passSvc  pb.UserChangePassSvcClient
	credsSvc pb.CredSvcClient

	TokenHeader string
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

func (c *ClientUC) GetToken() string {
	return c.Token
}

func (c *ClientUC) GetHeader() string {
	return c.TokenHeader
}

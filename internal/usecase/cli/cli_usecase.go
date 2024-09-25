package cli

import (
	"sync"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"passkeeper/internal/entities/config"
	"passkeeper/internal/entities/structs"
	pb "passkeeper/internal/transport/grpc/protofiles/proto"
	"passkeeper/internal/usecase/cli/blobsUC"
	"passkeeper/internal/usecase/cli/filesUC"
	"passkeeper/internal/usecase/cli/syncUC"
	"passkeeper/internal/usecase/cli/usersUC"
)

// ClientUsecase describe client functionality
type ClientUsecase interface {
	usersUC.UsersUsecase
	blobsUC.BlobsActionsUsecase
	blobsUC.GetBlobsUsecase
	filesUC.FilesUsecaseInf
	syncUC.SyncUsecase
}

// ClientUC implement ClientUsecase
type ClientUC struct {
	Addr string

	Authed bool
	User   *structs.User
	Token  string // JWT token

	Creds         []*structs.Credential
	Cards         []*structs.Card
	Notes         []*structs.Note
	Files         []*structs.File
	viewPageFocus bool
	SyncTime      time.Duration
	m             *sync.RWMutex

	userSvc  pb.UserSvcClient
	passSvc  pb.UserChangePassSvcClient
	credsSvc pb.BlobSvcClient

	log *zerolog.Logger

	*filesUC.FilesUC
}

// NewClientUC is a constructor for ClientUC
func NewClientUC(opts ...ClientUCoptions) (cli *ClientUC, err error) {

	cli = new(ClientUC)
	for _, opt := range opts {
		opt(cli)
	}

	if cli.Addr == "" {
		cli.Addr = config.DefaultAddr
	}

	if cli.SyncTime.Seconds() > config.SyncMax.Seconds() ||
		cli.SyncTime.Seconds() < config.SyncMin.Seconds() {
		cli.SyncTime = config.SyncDefault
	}

	conn, err := grpc.NewClient(cli.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// transport
	cli.userSvc = pb.NewUserSvcClient(conn)
	cli.passSvc = pb.NewUserChangePassSvcClient(conn)
	cli.credsSvc = pb.NewBlobSvcClient(conn)

	cli.m = &sync.RWMutex{}
	cli.FilesUC = new(filesUC.FilesUC)

	return cli, nil
}

// ClientUCoptions new type to use in constructor of ClientUC
type ClientUCoptions func(*ClientUC)

// WithAddr add address of server to ClientUC
func WithAddr(addr string) ClientUCoptions {
	return func(cli *ClientUC) {
		cli.Addr = addr
	}
}

// WithSyncTime setup sync time to ClientUC
func WithSyncTime(t time.Duration) ClientUCoptions {
	return func(cli *ClientUC) {
		cli.SyncTime = t
	}
}

// WithLogger setup logger to ClientUC
func WithLogger(l *zerolog.Logger) ClientUCoptions {
	return func(cli *ClientUC) {
		cli.log = l
	}
}

// GetToken return user token
func (c *ClientUC) GetToken() string {
	return c.Token
}

// CredsLen return credentiaonals len
func (c *ClientUC) CredsLen() int {
	c.m.RLock()
	defer c.m.RUnlock()
	return len(c.Creds)
}

func (c *ClientUC) CredsNotNil() bool {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.Creds != nil
}

// GetSyncTime return preset sync time
func (c *ClientUC) GetSyncTime() time.Duration {
	return c.SyncTime
}

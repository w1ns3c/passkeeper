package cli

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"passkeeper/internal/config"
	"passkeeper/internal/entities"

	pb "passkeeper/internal/transport/grpc/protofiles/proto"
)

var (
	ErrPassNotSame = fmt.Errorf("pass and pass repeat are not the same")
)

type ClientUsecase interface {
	Login(ctx context.Context, login, password string) (err error)
	Register(ctx context.Context, email, login, password, repeat string) error
	Logout()
	IsAuthed() bool

	ListCreds(ctx context.Context) (err error)
	EditCred(ctx context.Context, cred *entities.Credential, ind int) (err error)
	AddCred(ctx context.Context, cred *entities.Credential) (err error)
	DelCred(ctx context.Context, ind int) (err error)
	GetCredByIND(credID int) (cred *entities.Credential, err error)
	CredsLen() int
	CredsNotNil() bool

	SyncCreds(ctx context.Context)
	StopSync()
	ContinueSync()
	CheckSync() bool

	// moved from tuiApp
	GetToken() string
	GetHeader() string
	GetCreds() []*entities.Credential
	GetCards() []*entities.Card
	GetNotes() []*entities.Note
	GetSyncTime() time.Duration
}

type ClientUC struct {
	Addr string

	Authed bool
	User   *entities.User
	Token  string // JWT token
	//UserID      string // userID from JWT token
	//CredsSecret string // full secret for decrypt user's creds, contains sha1(clear_pass+secret_from_server)

	Creds         []*entities.Credential
	Cards         []*entities.Card
	Notes         []*entities.Note
	m             *sync.RWMutex
	viewPageFocus bool

	userSvc  pb.UserSvcClient
	passSvc  pb.UserChangePassSvcClient
	credsSvc pb.CredSvcClient

	TokenHeader string
	SyncTime    time.Duration

	log *zerolog.Logger
}

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

	cli.userSvc = pb.NewUserSvcClient(conn)
	cli.passSvc = pb.NewUserChangePassSvcClient(conn)
	cli.credsSvc = pb.NewCredSvcClient(conn)
	cli.m = &sync.RWMutex{}

	return cli, nil
}

type ClientUCoptions func(*ClientUC)

func WithAddr(addr string) ClientUCoptions {
	return func(cli *ClientUC) {
		cli.Addr = addr
	}
}

func WithSyncTime(t time.Duration) ClientUCoptions {
	return func(cli *ClientUC) {
		cli.SyncTime = t
	}
}

func WithLogger(l *zerolog.Logger) ClientUCoptions {
	return func(cli *ClientUC) {
		cli.log = l
	}
}

func (c *ClientUC) GetToken() string {
	return c.Token
}

func (c *ClientUC) GetHeader() string {
	return c.TokenHeader
}

func (c *ClientUC) CredsLen() int {
	c.m.RLock()
	defer c.m.RUnlock()
	return len(c.Creds)
}

func (c *ClientUC) GetSyncTime() time.Duration {
	return c.SyncTime
}

func (c *ClientUC) CredsNotNil() bool {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.Creds != nil
}

package cli

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"passkeeper/internal/entities"
	"passkeeper/internal/entities/config"
	pb "passkeeper/internal/transport/grpc/protofiles/proto"
	"passkeeper/internal/usecase/cli/filesUC"
)

var (
	ErrPassNotSame = fmt.Errorf("pass and pass repeat are not the same")
)

type ClientUsecase interface {
	Login(ctx context.Context, login, password string) (err error)
	Register(ctx context.Context, email, login, password, repeat string) error
	Logout()
	IsAuthed() bool

	GetBlobs(ctx context.Context) (err error)
	EditBlob(ctx context.Context, cred entities.CredInf, ind int) (err error)
	AddBlob(ctx context.Context, cred entities.CredInf) (err error)
	DelBlob(ctx context.Context, ind int, blobType entities.BlobType) (err error)

	GetCredByIND(credIND int) (cred *entities.Credential, err error)
	GetCardByIND(cardIND int) (cred *entities.Card, err error)
	GetNoteByIND(noteIND int) (cred *entities.Note, err error)
	GetFileByIND(ind int) (file *entities.File, err error)

	CredsLen() int
	CredsNotNil() bool

	SyncBlobs(ctx context.Context)
	StopSync()
	ContinueSync()
	CheckSync() bool

	GetCreds() []*entities.Credential
	GetCards() []*entities.Card
	GetNotes() []*entities.Note
	GetFiles() []*entities.File

	// moved from tuiApp
	GetToken() string
	GetSyncTime() time.Duration

	filesUC.FilesUsecaseInf
}

type ClientUC struct {
	Addr string

	Authed bool
	User   *entities.User
	Token  string // JWT token

	Creds         []*entities.Credential
	Cards         []*entities.Card
	Notes         []*entities.Note
	Files         []*entities.File
	viewPageFocus bool
	SyncTime      time.Duration
	m             *sync.RWMutex

	userSvc  pb.UserSvcClient
	passSvc  pb.UserChangePassSvcClient
	credsSvc pb.BlobSvcClient

	log *zerolog.Logger

	*filesUC.FilesUC
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

	// transport
	cli.userSvc = pb.NewUserSvcClient(conn)
	cli.passSvc = pb.NewUserChangePassSvcClient(conn)
	cli.credsSvc = pb.NewBlobSvcClient(conn)

	cli.m = &sync.RWMutex{}
	cli.FilesUC = new(filesUC.FilesUC)

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
